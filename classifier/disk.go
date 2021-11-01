package classifier

import (
	"strconv"
	"strings"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

// checkDisks it filters the bmh host as per the hardware details provided by user
func checkDisks(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	diskDetails := profile.Spec.HardwareCharacteristics.Disk
	if diskDetails == nil {
		return true
	}

	newDisk := host.Status.HardwareDetails.Storage
	if diskDetails.DiskSelector != nil {

		filteredDisk, matched := checkDisk(diskDetails.DiskSelector, host.Status.HardwareDetails.Storage)

		if !matched {
			log.Info("Disk Pattern",
				"host", host.Name,
				"profile", profile.Name,
				"namespace", host.Namespace,
				"ok", false,
			)
			return false
		} else if len(filteredDisk) > 0 {
			newDisk = filteredDisk
		}
	}

	ok := checkRangeInt(
		diskDetails.MinimumCount,
		diskDetails.MaximumCount,
		len(newDisk),
	)
	log.Info("Disk Count",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minCount", diskDetails.MinimumCount,
		"maxCount", diskDetails.MaximumCount,
		"actualCount", len(newDisk),
		"ok", ok,
	)
	if !ok {
		return false
	}

	for i, disk := range newDisk {

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

// checkRangeCapacity check the range of disk count
func checkRangeCapacity(min, max, count bmh.Capacity) bool {
	if min > 0 && count < min {
		return false
	}
	if max > 0 && count > max {
		return false
	}
	return true
}

// checkDisk it filter outs the disks from bmh disk array as per hardware details
func checkDisk(pattern []hwcc.DiskSelector, disks []bmh.Storage) ([]bmh.Storage, bool) {
	var diskNew []bmh.Storage

	for _, pattern := range pattern {
		matched := false
		for _, disk := range disks {
			validatedPattern := validatePattern(disk.HCTL)
			validatedExpectedPattern := validateExpectedPattern(pattern.HCTL)

			if validatedExpectedPattern == validatedPattern && pattern.Rotational == disk.Rotational {
				matched = true
				diskNew = append(diskNew, disk)
				log.Info("Disk Pattern",
					"expectedPattern", pattern.HCTL,
					"expectedRotational", pattern.Rotational,
					"actualPattern", disk.HCTL,
					"actualRotational", disk.Rotational,
					"ok", true,
				)
			}
		}

		if !matched {
			log.Info("Disk Pattern",
				"expectedPattern", pattern.HCTL,
				"expectedRotational", pattern.Rotational,
				"ok", false,
			)

			return diskNew, false
		}
	}

	return diskNew, true
}

// validatePattern finds out the disk with the pattern provided in hardware profile
func validatePattern(HCTL string) string {

	if HCTL == "0:0:0:0" {
		return "0:0:N:0"

	} else {
		resp := strings.Split(HCTL, ":")
		for index, st := range resp {
			value, _ := strconv.Atoi(st)
			if value > 0 {
				resp[index] = "N"
			}
		}
		return strings.Join(resp, ":")
	}
}

// validateExpectedPattern returns default pattern
func validateExpectedPattern(HCTL string) string {
	if HCTL == "0:0:0:0" {
		return "0:0:N:0"
	}
	return HCTL
}
