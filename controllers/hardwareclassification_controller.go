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

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
	"github.com/metal3-io/hardware-classification-controller/utils"
	"github.com/pkg/errors"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	//HWControllerName Name to show in the logs
	HWControllerName = "HardwareClassification-Controller"
)

// HardwareClassificationReconciler reconciles a HardwareClassification object
type HardwareClassificationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reconcile function
func (hcReconciler *HardwareClassificationReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()

	// Initialize the logger with namespace
	hwcLog := hcReconciler.Log.WithValues("hardwareclassification", req.NamespacedName)

	// Get HardwareClassificationController to get values for Namespace and ExpectedHardwareConfiguration
	hardwareClassification := &hwcc.HardwareClassification{}

	if err := hcReconciler.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Add a finalizer to newly created objects.
	if hardwareClassification.DeletionTimestamp.IsZero() && !hasFinalizer(hardwareClassification) {
		hwcLog.Info(
			"adding finalizer",
			"existingFinalizers", hardwareClassification.Finalizers,
			"newValue", hwcc.Finalizer,
		)
		hardwareClassification.Finalizers = append(hardwareClassification.Finalizers, hwcc.Finalizer)
		err := hcReconciler.Update(context.TODO(), hardwareClassification)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to add finalizer")
		}
		return ctrl.Result{}, nil
	}

	bmhHostList := bmh.BareMetalHostList{}
	opts := &client.ListOptions{
		// We only want to apply profiles to hosts in the same
		// namespace.
		Namespace: hardwareClassification.Namespace,
	}
	err := hcReconciler.List(context.TODO(), &bmhHostList, opts)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "could not fetch host list")
	}

	// Count hosts with our label. We use this value to decide whether
	// we have matched and whether it is OK to delete this profile.
	labelKey, _ := getLabelDetails(hardwareClassification)
	matchCount := 0
	for _, host := range bmhHostList.Items {
		labels := host.GetLabels()
		if labels == nil {
			continue
		}
		if _, ok := labels[labelKey]; !ok {
			continue
		}
		hwcLog.Info("found host with label",
			"host", host.Name,
			"label", labelKey,
		)
		matchCount++
	}

	// Wait to delete the hardwareClassification resource until no
	// hosts are labeled as matching its rules.
	if !hardwareClassification.DeletionTimestamp.IsZero() {

		if matchCount > 0 {
			hwcLog.Info("waiting to delete")
			// We do not need to explicitly ask to be requeued here
			// because we will be invoked again when the host(s) found
			// with our label are modified.
			return ctrl.Result{}, nil
		}

		// No hosts are using our label, so we can allow the delete to
		// proceed.
		hardwareClassification.Finalizers = utils.FilterStringFromList(
			hardwareClassification.Finalizers, hwcc.Finalizer)
		err = hcReconciler.Update(context.TODO(), hardwareClassification)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to remove finalizer")
		}

		hwcLog.Info("deleting")
		return ctrl.Result{}, nil
	}

	// Update our status to report whether we have matched a host or not.
	status := hwcc.ProfileMatchStatusMatched
	if matchCount == 0 {
		status = hwcc.ProfileMatchStatusUnMatched
	}
	if len(bmhHostList.Items) == 0 {
		status = hwcc.NoBareMetalHosts
	}

	if hardwareClassification.Status.ProfileMatchStatus != status {
		hwcLog.Info("updating match status", "newValue", status)
		hardwareClassification.Status.ProfileMatchStatus = status
		err = hcReconciler.Status().Update(context.TODO(), hardwareClassification)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to update status")
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func hasFinalizer(profile *hwcc.HardwareClassification) bool {
	return utils.StringInList(profile.Finalizers, hwcc.Finalizer)
}

// SetupWithManager will add watches for this controller
func (hcReconciler *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {

	mapper := classificationMapper{
		client: mgr.GetClient(),
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Named("hardware-classification").
		Watches(&source.Kind{Type: &bmh.BareMetalHost{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: &mapper}).
		Complete(hcReconciler)
}

type classificationMapper struct {
	client client.Client
}

func (m *classificationMapper) Map(obj handler.MapObject) []ctrl.Request {
	log := ctrl.Log.WithName("controllers").WithName("HardwareClassification").WithName("mapper").
		WithValues("BareMetalHost",
			fmt.Sprintf("%s/%s", obj.Meta.GetNamespace(), obj.Meta.GetName()))

	hwcList := hwcc.HardwareClassificationList{}
	opts := &client.ListOptions{
		// We only want to apply profiles to classification rules in
		// the same namespace.
		Namespace: obj.Meta.GetNamespace(),
	}
	err := m.client.List(context.TODO(), &hwcList, opts)
	if err != nil {
		log.Error(err, "could not fetch hardware classification list")
		return nil
	}

	requests := []ctrl.Request{}
	for _, profile := range hwcList.Items {
		log.Info("found hardwareclassification", "name", profile.Name)
		requests = append(requests, ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      profile.Name,
				Namespace: profile.Namespace,
			},
		})
	}
	return requests
}
