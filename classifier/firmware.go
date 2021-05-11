package classifier

import (
	"fmt"
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"strings"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

// checkFirmware it filters the bmh host as per the hardware details provided by user
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

	ok = checkVersion(firmwareDetails.BIOS.MinorVersion,
		firmwareDetails.BIOS.MajorVersion,
		host.Status.HardwareDetails.Firmware.BIOS.Version)

	log.Info("Firmware",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minor version", firmwareDetails.BIOS.MinorVersion,
		"major version", firmwareDetails.BIOS.MajorVersion,
		"actualVersion", host.Status.HardwareDetails.Firmware.BIOS.Version,
		"ok", ok,
	)
	if !ok {
		return false
	}

	return true
}

// Check the version range
func checkVersion(minorVersion, majorVersion, hostVersion string) bool {

	if minorVersion != "" && majorVersion != "" {

		convertedMinorVersion := appendZeros(minorVersion, 4)
		convertedMajorVersion := appendZeros(majorVersion, 4)
		convertedHostVersion := appendZeros(hostVersion, 4)

		if convertedMinorVersion != "" && convertedMajorVersion != "" && convertedHostVersion != "" {
			if convertedMinorVersion > convertedHostVersion || convertedMajorVersion < convertedHostVersion {
				return false
			}
		}
	}
	return true
}

// Append zero to make it a number
func appendZeros(version string, zeros int) string {
	splitVersions := strings.Split(version, ".")
	appendVersionString := ""

	for _, value := range splitVersions {
		appendVersionString += fmt.Sprintf("%0*s", zeros, value)
	}

	return appendVersionString
}
