package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func checkFirmware(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	firmwareDetails := profile.Spec.HardwareCharacteristics.Firmware
	if firmwareDetails == nil {
		return true
	}

	ok := checkStringEmpty(
		firmwareDetails.Bios.Vendor,
		firmwareDetails.Bios.Version)
	log.Info("Firmware",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"vendor", firmwareDetails.Bios.Vendor,
		"version", firmwareDetails.Bios.Version,
		"actual Vendor", host.Status.HardwareDetails.Firmware.BIOS.Vendor,
		"actual Version", host.Status.HardwareDetails.Firmware.BIOS.Version,
		"ok", ok,
	)
	if !ok {
		return false
	}
	if firmwareDetails.Bios.Vendor != host.Status.HardwareDetails.Firmware.BIOS.Vendor && firmwareDetails.Bios.Version != host.Status.HardwareDetails.Firmware.BIOS.Version {
		ok = false
		log.Info("Firmware",
			"host", host.Name,
			"profile", profile.Name,
			"namespace", host.Namespace,
			"vendor", firmwareDetails.Bios.Vendor,
			"version", firmwareDetails.Bios.Version,
			"actual Vendor", host.Status.HardwareDetails.Firmware.BIOS.Vendor,
			"actual Version", host.Status.HardwareDetails.Firmware.BIOS.Version,
			"ok", ok,
		)
		return false
	}

	return true
}
