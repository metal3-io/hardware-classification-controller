package controllers

import (
	"context"
	hwcc "hardware-classification-controller/api/v1alpha1"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func (hcReconciler *HardwareClassificationReconciler) updateProfileMatchStatus(req ctrl.Request, hc *hwcc.HardwareClassification, status hwcc.ProfileMatchStatus) {
	hcReconciler.Log.Info("Current status is:", "ProfileMatchStatus", status)
	if hc.SetProfileMatchStatus(status) {
		hcReconciler.Log.Info("clearing previous error message")
		hc.ClearError()

		hcReconciler.Log.Info("Upating status as:", "ProfileMatchStatus", status)
		err := hcReconciler.saveHWCCStatus(hc)
		if err != nil {
			hcReconciler.Log.Error(err, "Error while saving ProfileMatchStatus", "", "")
		}
	}
}

func (hcReconciler *HardwareClassificationReconciler) handleErrorConditions(req ctrl.Request, hc *hwcc.HardwareClassification, errorType hwcc.ErrorType, message string, status hwcc.ProfileMatchStatus) {
	hcReconciler.setErrorCondition(req, hc, hwcc.FetchBMHListFailure, message)

	hcReconciler.Log.Info("Upating status as:", "ProfileMatchStatus", status)
	hc.SetProfileMatchStatus(status)
	err := hcReconciler.saveHWCCStatus(hc)
	if err != nil {
		hcReconciler.Log.Error(err, "Error while saving ProfileMatchStatus", "", "")
	}
}

func (hcReconciler *HardwareClassificationReconciler) setErrorCondition(request ctrl.Request, hardwareClassification *hwcc.HardwareClassification, errType hwcc.ErrorType, message string) (changed bool, err error) {
	var log = logf.Log.WithName("controller.HardwareClassification")

	reqLogger := log.WithValues("Request.Namespace",
		request.Namespace, "Request.Name", request.Name)

	changed = hardwareClassification.SetErrorMessage(errType, message)
	if changed {
		reqLogger.Info(
			"adding error message",
			"message", message,
		)
		err = hcReconciler.saveHWCCStatus(hardwareClassification)
		if err != nil {
			err = errors.Wrap(err, "failed to update error message")
		}
	} else {
		reqLogger.Info(
			"aleady added error message",
			"message", message,
		)
	}

	return
}

func (hcReconciler *HardwareClassificationReconciler) saveHWCCStatus(hcc *hwcc.HardwareClassification) error {

	//Refetch hwcc again
	obj := hcc.Status.DeepCopy()
	err := hcReconciler.Client.Get(context.TODO(),

		client.ObjectKey{
			Name:      hcc.Name,
			Namespace: hcc.Namespace,
		},
		hcc,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to update HardwareClassification Status")
	}

	hcc.Status = *obj
	err = hcReconciler.Client.Status().Update(context.TODO(), hcc)
	return err
}
