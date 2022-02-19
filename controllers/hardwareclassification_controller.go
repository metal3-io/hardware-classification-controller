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
	"strings"

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
	failedLabelName  = "hardwareclassification-error"
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

	failedHostList := fetchFailedBmhHostList(bmhHostList)
	if len(failedHostList) > 0 {
		if changed := labelFailedHost(hcReconciler, failedHostList, ctx); changed != false {
			hwcLog.Info("set label ", "failed host list", failedHostList)
		}
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

	setHostCount(hardwareClassification, hwcc.MatchedCount(matchCount), hwcc.UnmatchedCount(len(bmhHostList.Items)-(len(failedHostList)+matchCount)))
	setErrHostCount(hardwareClassification, failedHostList)
	err = hcReconciler.Status().Update(context.TODO(), hardwareClassification)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to update status")
	}

	return ctrl.Result{}, nil
}

func setHostCount(hwc *hwcc.HardwareClassification, MatchedHost hwcc.MatchedCount, UnmatchedHost hwcc.UnmatchedCount) {
	hwc.Status.MatchedCount = MatchedHost
	fmt.Println("Updating Matched host count")
	hwc.Status.UnmatchedCount = UnmatchedHost
	fmt.Println("Updating Unmatched host count")
}

func setErrHostCount(hwc *hwcc.HardwareClassification, failedHosts []bmh.BareMetalHost) {
	registrationErrorCount := 0
	introspectionErrorCount := 0
	provisioningErrorCount := 0
	powerMgmtErrorCount := 0
	proviRegisErrorCount := 0
	preprationErrorCount := 0
	detachErrorCount := 0
	for _, host := range failedHosts {
		switch host.Status.ErrorType {
		case "registration error":
			registrationErrorCount += 1
		case "inspection error":
			introspectionErrorCount += 1
		case "provisioning error":
			provisioningErrorCount += 1
		case "power management error":
			powerMgmtErrorCount += 1
		case "provisioned registration error":
			proviRegisErrorCount += 1
		case "preparation error":
			preprationErrorCount += 1
		case "detach error":
			detachErrorCount += 1
		}
	}
	fmt.Println("Updating ErrorHost count")
	hwc.Status.ErrorHosts = hwcc.ErrorHosts(len(failedHosts))
	hwc.Status.RegistrationErrorHosts = hwcc.RegistrationErrorHosts(registrationErrorCount)
	hwc.Status.IntrospectionErrorHosts = hwcc.IntrospectionErrorHosts(introspectionErrorCount)
	hwc.Status.ProvisioningErrorHosts = hwcc.ProvisioningErrorHosts(provisioningErrorCount)
	hwc.Status.PowerMgmtErrorHosts = hwcc.PowerMgmtErrorHosts(powerMgmtErrorCount)
	hwc.Status.ProvisionedRegistrationErrorHosts = hwcc.ProvisionedRegistrationErrorHosts(proviRegisErrorCount)
	hwc.Status.PreparationErrorHosts = hwcc.PreparationErrorHosts(preprationErrorCount)
	hwc.Status.DetachErrorHosts = hwcc.DetachErrorHosts(detachErrorCount)
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

func labelFailedHost(hcReconciler *HardwareClassificationReconciler, failedHostList []bmh.BareMetalHost, ctx context.Context) bool {
	for _, host := range failedHostList {
		labels := host.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labelValue := strings.ReplaceAll(string(host.Status.ErrorType), " ", "-")
		// If we already have a label with the same value no change is
		// needed.
		if val, ok := labels[failedLabelName]; ok {
			if val == labelValue {
				return false
			}
		}
		labels[failedLabelName] = labelValue
		host.SetLabels(labels)
		if err := hcReconciler.Client.Update(ctx, &host); err != nil {
			return false
		}
	}
	return true
}

func fetchFailedBmhHostList(bmhHostList bmh.BareMetalHostList) (failedHostList []bmh.BareMetalHost) {
	// Get hosts in error status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.HardwareDetails == nil && host.Status.OperationalStatus == "error" {
			failedHostList = append(failedHostList, host)
		}
	}
	return failedHostList
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
