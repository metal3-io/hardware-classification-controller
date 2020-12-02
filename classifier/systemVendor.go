package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func checkSystemVendor(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	systemVendorDetails := profile.Spec.HardwareCharacteristics.SystemVendor
	if systemVendorDetails == nil {
		return true
	}

	ok := checkStringEmpty(
		systemVendorDetails.Manufacturer,
		systemVendorDetails.ProductName)
	log.Info("SystemVendor",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"Manufacturer", systemVendorDetails.Manufacturer,
		"ProductName", systemVendorDetails.ProductName,
		"actual Manufacturer", host.Status.HardwareDetails.SystemVendor.Manufacturer,
		"actual ProductName", host.Status.HardwareDetails.SystemVendor.ProductName,
		"ok", ok,
	)
	if !ok {
		return false
	}
	if systemVendorDetails.Manufacturer != host.Status.HardwareDetails.SystemVendor.Manufacturer && systemVendorDetails.ProductName != host.Status.HardwareDetails.SystemVendor.ProductName {
		ok = false
		log.Info("SystemVendor",
			"host", host.Name,
			"profile", profile.Name,
			"namespace", host.Namespace,
			"Manufacturer", systemVendorDetails.Manufacturer,
			"ProductName", systemVendorDetails.ProductName,
			"actual Manufacturer", host.Status.HardwareDetails.SystemVendor.Manufacturer,
			"actual ProductName", host.Status.HardwareDetails.SystemVendor.ProductName,
			"ok", ok,
		)
		return false
	}

	return true
}
