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
	"strings"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
	"github.com/metal3-io/hardware-classification-controller/hcmanager"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(hardwareClassification, hcReconciler.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	defer func() {
		// Always attempt to Patch the hardwareClassification object and status after each reconciliation.
		if err := patchHelper.Patch(ctx, hardwareClassification); err != nil {
			hwcLog.Error(err, "Failed to Patch HardwareClassification")
		}
	}()

	// Get Expected Hardware Configuration from hardwareClassification
	extractedProfile := hardwareClassification.Spec.HardwareCharacteristics
	hwcLog.Info("Expected Hardware Configuration", "Profile", extractedProfile)

	// Get the new hardware classification manager
	hcManager := hcmanager.NewHardwareClassificationManager(hcReconciler.Client, hwcLog)

	ErrValidation := hcManager.ValidateExtractedHardwareProfile(extractedProfile)
	if ErrValidation != nil {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.ProfileMisConfigured, ErrValidation.Error())
		hwcLog.Error(ErrValidation, ErrValidation.Error())
		return ctrl.Result{}, nil
	}

	// stop processing the classification rules
	v := true
	if v {
		return ctrl.Result{}, nil
	}

	//Fetch baremetal host list for the given namespace
	hostList, bmhList, err := hcManager.FetchBmhHostList(hardwareClassification.ObjectMeta.Namespace)
	if err != nil {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.FetchBMHListFailure, err.Error())
		hwcLog.Error(err, err.Error())
		return ctrl.Result{}, nil
	}

	if len(hostList) == 0 {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.Empty, hwcc.NoBaremetalHost)
		hwcLog.Info(hwcc.NoBaremetalHost)
		return ctrl.Result{}, nil
	}

	//Extract the hardware details from the baremetal host list
	validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(extractedProfile, hostList)
	hwcLog.Info("Validated Hardware Details", "HardwareDetails", validatedHardwareDetails)

	//Compare the host list with extracted profile and fetch the valid host names
	validHosts := hcManager.MinMaxFilter(hardwareClassification.ObjectMeta.Name, validatedHardwareDetails, extractedProfile)
	hwcLog.Info("Filtered Bare metal hosts", "ValidHosts", validHosts)

	updateLabelError := hcManager.UpdateLabels(ctx, hardwareClassification.ObjectMeta, validHosts, bmhList)

	if len(updateLabelError) > 0 {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty,
			hwcc.LabelUpdateFailure, strings.Join(updateLabelError, ","))
		hwcLog.Error(nil, hwcc.UpdateLabelError, "Update Label Error", strings.Join(updateLabelError, ","))
	} else if len(validHosts) > 0 && len(updateLabelError) == 0 {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusMatched,
			hwcc.Empty, hwcc.NOError)
		hwcLog.Info("Updated profile status", "ProfileMatchStatus", hwcc.ProfileMatchStatusMatched)
	} else {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusUnMatched,
			hwcc.Empty, hwcc.NOError)
		hwcLog.Info("Updated profile status", "ProfileMatchStatus", hwcc.ProfileMatchStatusUnMatched)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager will add watches for this controller
func (hcReconciler *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Complete(hcReconciler)
}
