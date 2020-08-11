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
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

// ExtractAndValidateHardwareDetails this function will return map containing introspection details for a host.
func (mgr HardwareClassificationManager) ExtractAndValidateHardwareDetails(extractedProfile hwcc.HardwareCharacteristics,
	bmhList []bmh.BareMetalHost) []bmh.HardwareDetails {

	var validHosts []bmh.HardwareDetails
	if extractedProfile != (hwcc.HardwareCharacteristics{}) {
		for _, host := range bmhList {
			var hardwareDetails bmh.HardwareDetails
			// Get the CPU details from the baremetal host and validate it
			if extractedProfile.Cpu != nil {
				hardwareDetails.CPU = bmh.CPU{
					Count:          host.Status.HardwareDetails.CPU.Count,
					ClockMegahertz: bmh.ClockSpeed(host.Status.HardwareDetails.CPU.ClockMegahertz),
				}
			}
			// Get the Storage details from the baremetal host and validate it
			if extractedProfile.Disk != nil {
				var disks []bmh.Storage
				for _, disk := range host.Status.HardwareDetails.Storage {
					disks = append(disks, bmh.Storage{Name: disk.Name, SizeBytes: ConvertBytesToGb(disk.SizeBytes)})
				}
				hardwareDetails.Storage = disks
			}
			// Get the NIC details from the baremetal host and validate it
			if extractedProfile.Nic != nil {
				hardwareDetails.NIC = host.Status.HardwareDetails.NIC
			}
			// Get the RAM details from the baremetal host and validate it
			if extractedProfile.Ram != nil {
				hardwareDetails.RAMMebibytes = host.Status.HardwareDetails.RAMMebibytes / 1024
			}

			//Check if hardware details are not empty to add it in validhost
			hardwareDetails.Hostname = host.ObjectMeta.Name
			validHosts = append(validHosts, hardwareDetails)
		}
	}
	return validHosts
}
