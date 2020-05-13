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

package vallidation

import (
	valTypes "hardware-classification-controller/validation/validationModel"
	"net"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

//CheckValidIP uses net package to check if the IP is valid or not
func CheckValidIP(NICIp string) bool {
	return net.ParseIP(NICIp) != nil
}

//ConvertBytesToGb it converts the Byte into GB
func ConvertBytesToGb(inBytes int64) int64 {
	inGb := (inBytes / 1024 / 1024 / 1024)
	return inGb
}

//Validation this function will validate the host and create a new m ap with structered details
func Validation(hostDetails map[string]map[string]interface{}) map[string]map[string]interface{} {

	validatedHostMap := make(map[string]map[string]interface{})

	for hostName, details := range hostDetails {
		hardwareDetails := make(map[string]interface{})
		for key, value := range details {

			// Get the CPU details from the ironic host and validate it into new structure
			cpu, ok := value.(bmh.CPU)
			if ok {
				validCPU := valTypes.CPU{
					Count:      cpu.Count,
					ClockSpeed: float64(cpu.ClockMegahertz) / 1000,
				}
				hardwareDetails[key] = validCPU
			}

			// Get the RAM details from the ironic host and validate it into new structure
			ram, ok := value.(int)
			if ok {
				validRAM := valTypes.RAM{
					RAMGb: ram / 1024,
				}
				hardwareDetails[key] = validRAM
			}

			// Get the NICS details from the ironic host and validate it into new structure
			nics, ok := value.([]bmh.NIC)
			if ok {
				var validNICS valTypes.NIC
				for _, NIC := range nics {
					if NIC.PXE && CheckValidIP(NIC.IP) {
						validNICS.Name = NIC.Name
						validNICS.PXE = NIC.PXE
					}
				}

				validNICS.Count = len(nics)
				hardwareDetails[key] = validNICS
			}

			// Get the Storage details from the ironic host and validate it into new structure
			storage, ok := value.([]bmh.Storage)
			if ok {
				var disks []valTypes.Disk

				for _, disk := range storage {
					disks = append(disks, valTypes.Disk{Name: disk.Name, SizeGb: ConvertBytesToGb(int64(disk.SizeBytes))})
				}

				validStorage := valTypes.Storage{
					Count: len(disks),
					Disk:  disks,
				}

				hardwareDetails[key] = validStorage
			}

		}

		validatedHostMap[hostName] = hardwareDetails

	}

	return validatedHostMap

}
