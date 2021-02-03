package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
	"strings"
)

func checkNICs(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	nicDetails := profile.Spec.HardwareCharacteristics.Nic
	var nicVendors []string
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

	if len(nicDetails.NicSelector.Vendor) == 0 {
		return ok
	}

	for i := 0; i < len(host.Status.HardwareDetails.NIC); i++ {
		vendorProduct := strings.Fields(host.Status.HardwareDetails.NIC[i].Model)
		nicVendors = append(nicVendors, vendorProduct[0])
	}
	ok = checkVendor(
		nicVendors,
		nicDetails.NicSelector.Vendor,
	)

	log.Info("NIC",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"Require Nics vendor", nicDetails.NicSelector.Vendor,
		"Actual Nics Vendor", nicVendors,
		"ok", ok,
	)

	return ok
}

func checkVendor(actualNicVendors, requireNicVendor []string) bool {
	if len(requireNicVendor) > len(actualNicVendors) {
		return false
	}
	nicCount := len(requireNicVendor)
	for _, requireNic := range requireNicVendor {
		for _, actualNic := range actualNicVendors {
			if requireNic == actualNic {
				nicCount -= 1
				break
			}
		}
	}

	if nicCount == 0 {
		return true
	} else {
		return false
	}
}
