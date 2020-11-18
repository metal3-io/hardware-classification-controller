package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func checkDisks(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	diskDetails := profile.Spec.HardwareCharacteristics.Disk
	if diskDetails == nil {
		return true
	}

	ok := checkRangeInt(
		diskDetails.MinimumCount,
		diskDetails.MaximumCount,
		len(host.Status.HardwareDetails.Storage),
	)
	log.Info("DiskCount",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minCount", diskDetails.MinimumCount,
		"maxCount", diskDetails.MinimumCount,
		"actualCount", len(host.Status.HardwareDetails.Storage),
		"ok", ok,
	)
	if !ok {
		return false
	}

	for i, disk := range host.Status.HardwareDetails.Storage {

		// The disk size is reported on the host in bytes and the
		// classification rule is given in GB, so we have to convert
		// the values to the same units. Reducing bytes to GiB loses
		// detail, so we convert GB to bytes.
		minSize := bmh.Capacity(diskDetails.MinimumIndividualSizeGB) * bmh.GigaByte
		maxSize := bmh.Capacity(diskDetails.MaximumIndividualSizeGB) * bmh.GigaByte

		ok := checkRangeCapacity(
			minSize,
			maxSize,
			disk.SizeBytes,
		)
		log.Info("DiskSize",
			"host", host.Name,
			"profile", profile.Name,
			"namespace", host.Namespace,
			"minSize", minSize,
			"maxSize", maxSize,
			"actualSize", disk.SizeBytes,
			"diskNum", i,
			"diskName", disk.Name,
			"ok", ok,
		)
		if !ok {
			return false
		}
	}

	return true
}

func checkRangeCapacity(min, max, count bmh.Capacity) bool {
	if min > 0 && count < min {
		return false
	}
	if max > 0 && count > max {
		return false
	}
	return true
}
