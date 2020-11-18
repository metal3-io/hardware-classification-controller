package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func checkRAM(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	ramDetails := profile.Spec.HardwareCharacteristics.Ram
	if ramDetails == nil {
		return true
	}

	// The size reported on the host is in MiB and the classification
	// rule is given in GB, so we have to convert the values to the
	// same units. Reducing MiB to GB loses detail, so we convert GB
	// to MiB.
	actualSize := host.Status.HardwareDetails.RAMMebibytes
	minSize := ramDetails.MinimumSizeGB * 1024
	maxSize := ramDetails.MaximumSizeGB * 1024

	ok := checkRangeInt(minSize, maxSize, actualSize)
	log.Info("RAM",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minSize", minSize,
		"maxSize", maxSize,
		"actualSize", actualSize,
		"ok", ok,
	)

	return ok
}
