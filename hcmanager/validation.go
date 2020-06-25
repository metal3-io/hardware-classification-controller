/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hcmanager

import (
	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

// ExtractAndValidateHardwareDetails this function will return map containing introspection details for a host.
func (mgr HardwareClassificationManager) ExtractAndValidateHardwareDetails(extractedProfile hwcc.HardwareCharacteristics,
	bmhList []bmh.BareMetalHost) map[string]map[string]interface{} {

	validatedHostMap := make(map[string]map[string]interface{})

	if extractedProfile != (hwcc.HardwareCharacteristics{}) {
		for _, host := range bmhList {
			hardwareDetails := make(map[string]interface{})

			if extractedProfile.Cpu != nil {
				// Get the CPU details from the baremetal host and validate it into new structure
				validCPU := bmh.CPU{
					Count:          host.Status.HardwareDetails.CPU.Count,
					ClockMegahertz: bmh.ClockSpeed(host.Status.HardwareDetails.CPU.ClockMegahertz) / 1000,
				}
				hardwareDetails[CPULabel] = validCPU
			}

			if extractedProfile.Disk != nil {
				// Get the Storage details from the baremetal host and validate it into new structure
				var disks []bmh.Storage

				for _, disk := range host.Status.HardwareDetails.Storage {
					disks = append(disks, bmh.Storage{Name: disk.Name, SizeBytes: ConvertBytesToGb(disk.SizeBytes)})
				}
				hardwareDetails[DISKLabel] = disks
			}

			if extractedProfile.Nic != nil {
				// Get the NIC details from the baremetal host and validate it into new structure
				hardwareDetails[NICLabel] = len(host.Status.HardwareDetails.NIC)
			}

			if extractedProfile.Ram != nil {
				// Get the RAM details from the baremetal host and validate it into new structure
				hardwareDetails[RAMLabel] = int64(host.Status.HardwareDetails.RAMMebibytes / 1024)
			}

			if len(hardwareDetails) != 0 {
				validatedHostMap[host.ObjectMeta.Name] = hardwareDetails
			}
		}
	}
	return validatedHostMap
}
