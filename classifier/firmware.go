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

	ok := checkString(firmwareDetails.BIOS.Vendor, host.Status.HardwareDetails.Firmware.BIOS.Vendor)

	log.Info("Firmware",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"vendor", firmwareDetails.BIOS.Vendor,
		"actualVendor", host.Status.HardwareDetails.Firmware.BIOS.Vendor,
		"ok", ok,
	)
	if !ok {
		return false
	}

	ok = checkString(firmwareDetails.BIOS.Version, host.Status.HardwareDetails.Firmware.BIOS.Version)

	log.Info("Firmware",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"version", firmwareDetails.BIOS.Version,
		"actualVersion", host.Status.HardwareDetails.Firmware.BIOS.Version,
		"ok", ok,
	)
	if !ok {
		return false
	}

	return true
}

func checkString(expected, hostSpecific string) bool {
	if expected != "" {
		if expected != hostSpecific {
			return false
		}

	}
	return true
}
