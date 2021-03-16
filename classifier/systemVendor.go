package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"strings"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func checkSystemVendor(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	systemVendorDetails := profile.Spec.HardwareCharacteristics.SystemVendor
	if systemVendorDetails == nil {
		return true
	}

	ok := checkString(systemVendorDetails.Manufacturer, host.Status.HardwareDetails.SystemVendor.Manufacturer)
	log.Info("System Vendor",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"Manufacturer", systemVendorDetails.Manufacturer,
		"actual Manufacturer", host.Status.HardwareDetails.SystemVendor.Manufacturer,
		"ok", ok,
	)

	if !ok {
		return false
	}

	ok = checkSubString(systemVendorDetails.ProductName, host.Status.HardwareDetails.SystemVendor.ProductName)
	log.Info("System Vendor",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"ProductName", systemVendorDetails.ProductName,
		"actual ProductName", host.Status.HardwareDetails.SystemVendor.ProductName,
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

func checkSubString(expected, hostSpecific string) bool {
	if expected != "" {
		if !strings.Contains(hostSpecific, expected) {
			return false
		}

	}
	return true
}
