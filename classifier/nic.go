package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func checkNICs(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	nicDetails := profile.Spec.HardwareCharacteristics.Nic
	if nicDetails == nil {
		return true
	}

	ok := checkRangeInt(
		nicDetails.MinimumCount,
		nicDetails.MaximumCount,
		len(host.Status.HardwareDetails.NIC),
	)
	log.Info("NIC",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minCount", nicDetails.MinimumCount,
		"maxCount", nicDetails.MaximumCount,
		"actualCount", len(host.Status.HardwareDetails.NIC),
		"ok", ok,
	)

	return ok
}
