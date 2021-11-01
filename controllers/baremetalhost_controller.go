/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
	"github.com/metal3-io/hardware-classification-controller/classifier"
)

const (
	defaultLabelName  = "hardwareclassification.metal3.io/"
	defaultLabelValue = "matches"
)

// BareMetalHostReconciler reconciles a BareMetalHost object
type BareMetalHostReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *BareMetalHostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	logger := r.Log.WithValues("host", req.NamespacedName)

	logger.Info("reconciling")

	host := &bmh.BareMetalHost{}
	err := r.Get(context.TODO(), req.NamespacedName, host)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after
			// reconcile request.  Owned objects are automatically
			// garbage collected. For additional cleanup logic use
			// finalizers.  Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, errors.Wrap(err, "could not load host data")
	}

	if host.Status.HardwareDetails == nil {
		logger.Info("no hardware details")
		return ctrl.Result{}, nil
	}

	profileList := hwcc.HardwareClassificationList{}
	opts := &client.ListOptions{
		// We only want to apply profiles in the same namespace as the
		// host.
		Namespace: req.Namespace,
	}
	err = r.List(context.TODO(), &profileList, opts)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "could not fetch classification profiles")
	}

	changed := false
	for _, profile := range profileList.Items {
		labelKey, labelValue := getLabelDetails(&profile)

		switch {
		case !profile.DeletionTimestamp.IsZero():
			logger.Info("profile is being deleted", "profile", profile.Name)
			changed = deleteLabel(host, labelKey) || changed
			if changed {
				logger.Info("removed label", "name", labelKey, "value", labelValue)
			}
		case !classifier.ProfileMatchesHost(&profile, host):
			changed = deleteLabel(host, labelKey) || changed
			if changed {
				logger.Info("removed label", "name", labelKey, "value", labelValue)
			}
		default:
			changed = setLabel(host, labelKey, labelValue) || changed
			if changed {
				logger.Info("set label", "name", labelKey, "value", labelValue)
			}
		}
	}

	if changed {
		if err := r.Update(context.TODO(), host); err != nil {
			return ctrl.Result{}, errors.Wrap(err,
				fmt.Sprintf("failed to update host %s/%s", host.Namespace, host.Name))
		}
	}

	return ctrl.Result{}, nil
}

func getLabelDetails(profile *hwcc.HardwareClassification) (key, value string) {
	key = defaultLabelName + profile.Name
	labels := profile.GetLabels()
	if labels != nil {
		if val, ok := labels[profile.Name]; ok {
			value = val
		}
	}
	if value == "" {
		value = defaultLabelValue
	}
	return
}

func deleteLabel(host *bmh.BareMetalHost, labelKey string) bool {
	labels := host.GetLabels()

	if labels == nil {
		return false
	}

	if _, ok := labels[labelKey]; !ok {
		return false
	}

	delete(labels, labelKey)
	host.SetLabels(labels)
	return true
}

func setLabel(host *bmh.BareMetalHost, labelKey string, labelValue string) bool {
	labels := host.GetLabels()

	if labels == nil {
		labels = make(map[string]string)
	}

	// If we already have a label with the same value no change is
	// needed.
	if val, ok := labels[labelKey]; ok {
		if val == labelValue {
			return false
		}
	}

	labels[labelKey] = labelValue
	host.SetLabels(labels)
	return true
}

func (r *BareMetalHostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bmh.BareMetalHost{}).
		Named("baremetalhost").
		Watches(
			&source.Kind{Type: &hwcc.HardwareClassification{}},
			handler.EnqueueRequestsFromMapFunc(r.Map),
		).
		Complete(r)
}

func (m *BareMetalHostReconciler) Map(obj client.Object) []ctrl.Request {
	requests := []ctrl.Request{}

	if bmhs, ok := obj.(*bmh.BareMetalHost); ok {

		log := ctrl.Log.WithName("controllers").WithName("BareMetalHost").WithName("mapper").
			WithValues("HardwareClassification",
				fmt.Sprintf("%s/%s", bmhs.ObjectMeta.Namespace, bmhs.ObjectMeta.Name))

		bmhHostList := bmh.BareMetalHostList{}
		opts := &client.ListOptions{
			// We only want to apply profiles to hosts in the same
			// namespace.
			Namespace: bmhs.ObjectMeta.Namespace,
		}
		err := m.Client.List(context.TODO(), &bmhHostList, opts)
		if err != nil {
			log.Error(err, "could not fetch host list")
			return nil
		}

		for _, host := range bmhHostList.Items {
			log.Info("found host", "name", host.Name)
			requests = append(requests, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      host.Name,
					Namespace: host.Namespace,
				},
			})
		}
	}
	return requests
}
