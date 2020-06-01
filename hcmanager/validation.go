package hcmanager

import (
	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

// ExtractAndValidateHardwareDetails this function will return map containing introspection details for a host.
func (mgr HardwareClassificationManager) ExtractAndValidateHardwareDetails(extractedProfile hwcc.ExpectedHardwareConfiguration,
	bmhList []bmh.BareMetalHost) map[string]map[string]interface{} {

	validatedHostMap := make(map[string]map[string]interface{})

	if extractedProfile != (hwcc.ExpectedHardwareConfiguration{}) {
		for _, host := range bmhList {
			hardwareDetails := make(map[string]interface{})

			if extractedProfile.CPU != nil {
				// Get the CPU details from the baremetal host and validate it into new structure
				validCPU := CPU{
					Count:      host.Status.HardwareDetails.CPU.Count,
					ClockSpeed: float64(host.Status.HardwareDetails.CPU.ClockMegahertz) / 1000,
				}
				hardwareDetails[CPULable] = validCPU
			}

			if extractedProfile.Disk != nil {
				// Get the Storage details from the baremetal host and validate it into new structure
				var disks []Disk

				for _, disk := range host.Status.HardwareDetails.Storage {
					disks = append(disks, Disk{Name: disk.Name, SizeGb: ConvertBytesToGb(int64(disk.SizeBytes))})
				}

				validStorage := Storage{
					Count: len(disks),
					Disk:  disks,
				}
				hardwareDetails[DISKLable] = validStorage
			}

			if extractedProfile.NIC != nil {
				// Get the NIC details from the baremetal host and validate it into new structure
				var validNICS NIC
				for _, NIC := range host.Status.HardwareDetails.NIC {
					if NIC.PXE && CheckValidIP(NIC.IP) {
						validNICS.Name = NIC.Name
						validNICS.PXE = NIC.PXE
					}
				}

				validNICS.Count = len(host.Status.HardwareDetails.NIC)
				hardwareDetails[NICLable] = validNICS
			}

			if extractedProfile.RAM != nil {
				// Get the RAM details from the baremetal host and validate it into new structure
				validRAM := RAM{
					RAMGb: host.Status.HardwareDetails.RAMMebibytes / 1024,
				}
				hardwareDetails[RAMLable] = validRAM
			}

			if len(hardwareDetails) != 0 {
				validatedHostMap[host.ObjectMeta.Name] = hardwareDetails
			}
		}
	}
	return validatedHostMap
}
