package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

var log = ctrl.Log.WithName("classifier")

func ProfileMatchesHost(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	if !checkCPU(profile, host) {
		return false
	}
	return true
}

func checkRangeInt(min, max, count int) bool {
	if min > 0 && count < min {
		return false
	}
	if max > 0 && count > max {
		return false
	}
	return true
}
