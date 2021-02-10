package classifier

import (
	"strconv"
	"strings"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

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
	log.Info("Disk Pattern",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minCount", diskDetails.MinimumCount,
		"maxCount", diskDetails.MinimumCount,
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

func checkRangeCapacity(min, max, count bmh.Capacity) bool {
	if min > 0 && count < min {
		return false
	}
	if max > 0 && count > max {
		return false
	}
	return true
}

func checkDisk(pattern []hwcc.DiskSelector, disks []bmh.Storage) ([]bmh.Storage, bool) {
	var diskNew []bmh.Storage

	for _, pattern := range pattern {
		matched := false
		for _, disk := range disks {
			replacedString := replaceCharacters(disk.HCTL)
			if pattern.HCTL == replacedString && pattern.Rotational == disk.Rotational {
				matched = true
				diskNew = append(diskNew, disk)
			}
		}

		if !matched {
			log.Info("Disk Pattern",
				"pattern", pattern,
				"ok", false,
			)

			return diskNew, false
		}
	}

	return diskNew, true
}

func replaceCharacters(HCTL string) string {

	if HCTL == "0:0:0:0" {
		return "0:0:N:0"

	} else {
		res1 := strings.Split(HCTL, ":")
		for index, st := range res1 {
			value, _ := strconv.Atoi(st)
			if value > 0 {
				res1[index] = "N"
			}
		}
		return strings.Join(res1, ":")
	}
}
