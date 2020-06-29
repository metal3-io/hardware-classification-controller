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

// MinMaxComparison it will compare the minimum and maximum comparison based on the value provided by the user and check for the valid host
func (mgr HardwareClassificationManager) MinMaxComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.HardwareCharacteristics) []string {

	var comparedHost []string

	for hostname, details := range validatedHost {
		isHostValid := true
		for _, value := range details {
			isValid := false
			cpu, CPUOK := value.(bmh.CPU)
			if CPUOK {
				if checkCPUCount(mgr, cpu, expectedHardwareprofile.Cpu) {
					isValid = true
				}
			}

			ram, RAMOK := value.(int64)
			if RAMOK {
				if checkRAM(mgr, ram, expectedHardwareprofile.Ram) {
					isValid = true
				}
			}

			nics, NICSOK := value.(int)
			if NICSOK {
				if checkNICS(mgr, nics, expectedHardwareprofile.Nic) {
					isValid = true
				}
			}

			disk, DISKOK := value.([]bmh.Storage)
			if DISKOK {
				if checkDiskDetails(mgr, disk, expectedHardwareprofile.Disk) {
					isValid = true
				}
			}

			if !isValid {
				isHostValid = false
				break
			}

		}

		if isHostValid {
			comparedHost = append(comparedHost, hostname)
			mgr.Log.Info(hostname, " Matches profile ", ProfileName)

		} else {
			mgr.Log.Info(hostname, " Did not matches profile ", ProfileName)
		}

	}

	return comparedHost

}

//checkCPUCount this function checks the CPU details for both min and max parameters
func checkCPUCount(mgr HardwareClassificationManager, cpu bmh.CPU, expectedCPU *hwcc.Cpu) bool {

	if (expectedCPU.MaximumCount > 0) && (expectedCPU.MinimumCount > 0) {

		mgr.Log.Info("", "Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		mgr.Log.Info("", "Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if (expectedCPU.MinimumCount > cpu.Count) || (expectedCPU.MaximumCount < cpu.Count) {
			mgr.Log.Info("CPU MINMAX COUNT did not match")
			return false
		}

	} else if expectedCPU.MaximumCount > 0 {

		mgr.Log.Info("", "Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MaximumCount < cpu.Count {

			mgr.Log.Info("CPU MAX COUNT did not match")
			return false
		}

	} else if expectedCPU.MinimumCount > 0 {

		mgr.Log.Info("", "Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MinimumCount > cpu.Count {
			mgr.Log.Info("CPU MIN COUNT did not match")
			return false
		}

	}

	if (expectedCPU.MaximumSpeedMHz > 0) && (expectedCPU.MinimumSpeedMHz > 0) {

		MinSpeed := bmh.ClockSpeed(expectedCPU.MinimumSpeedMHz)
		MaxSpeed := bmh.ClockSpeed(expectedCPU.MaximumSpeedMHz)

		mgr.Log.Info("", "Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockMegahertz)
		mgr.Log.Info("", "Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockMegahertz)
		if (MinSpeed > cpu.ClockMegahertz) || (MaxSpeed < cpu.ClockMegahertz) {
			mgr.Log.Info("CPU MINMAX ClockSpeed did not match")
			return false
		}
	} else if expectedCPU.MaximumSpeedMHz > 0 {
		MaxSpeed := bmh.ClockSpeed(float64(expectedCPU.MaximumSpeedMHz / 1000))
		mgr.Log.Info("", "Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockMegahertz)
		if MaxSpeed < cpu.ClockMegahertz {
			mgr.Log.Info("CPU MAX ClockSpeed did not match")
			return false
		}
	} else if expectedCPU.MinimumSpeedMHz > 0 {
		MinSpeed := bmh.ClockSpeed(float64(expectedCPU.MinimumSpeedMHz / 1000))
		mgr.Log.Info("", "Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockMegahertz)
		if MinSpeed > cpu.ClockMegahertz {
			mgr.Log.Info("CPU MIN ClockSpeed did not match")
			return false
		}

	}

	return true

}

//checkNICS this function checks the nics details for both min and max parameters
func checkNICS(mgr HardwareClassificationManager, nics int, expectedNIC *hwcc.Nic) bool {

	if (expectedNIC.MaximumCount > 0) && (expectedNIC.MinimumCount > 0) {

		mgr.Log.Info("", "Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics)
		mgr.Log.Info("", "Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics)
		if (expectedNIC.MinimumCount > nics) || (expectedNIC.MaximumCount < nics) {

			mgr.Log.Info("NICS MINMAX count did not match")
			return false
		}
	} else if expectedNIC.MaximumCount > 0 {

		mgr.Log.Info("", "Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics)
		if expectedNIC.MaximumCount < nics {

			mgr.Log.Info("NICS MAX count did not match")
			return false
		}

	} else if expectedNIC.MinimumCount > 0 {

		mgr.Log.Info("", "Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics)
		if expectedNIC.MinimumCount > nics {

			mgr.Log.Info("NICS MIN count did not match")
			return false
		}

	}
	return true
}

//checkRAM this function checks the ram details for both min and max parameters
func checkRAM(mgr HardwareClassificationManager, ram int64, expectedRAM *hwcc.Ram) bool {
	if (expectedRAM.MaximumSizeGB > 0) && (expectedRAM.MinimumSizeGB > 0) {

		mgr.Log.Info("", "Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram)
		mgr.Log.Info("", "Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram)
		if (expectedRAM.MinimumSizeGB > ram) || (expectedRAM.MaximumSizeGB < ram) {
			mgr.Log.Info("RAM MINMAX SIZE did not match")
			return false
		}
	} else if expectedRAM.MaximumSizeGB > 0 {
		mgr.Log.Info("", "Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram)
		if expectedRAM.MaximumSizeGB < ram {
			mgr.Log.Info("RAM MAX SIZE did not match")
			return false
		}

	} else if expectedRAM.MinimumSizeGB > 0 {

		mgr.Log.Info("", "Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram)
		if expectedRAM.MinimumSizeGB > ram {

			mgr.Log.Info("RAM MIN SIZE did not match")
			return false
		}

	}
	return true
}

//checkDiskDetails this function checks the Disk details for both min and max parameters
func checkDiskDetails(mgr HardwareClassificationManager, disk []bmh.Storage, expectedDisk *hwcc.Disk) bool {

	if (expectedDisk.MaximumCount > 0) && (expectedDisk.MinimumCount > 0) {
		mgr.Log.Info("", "Provided Minimum count for Disk", expectedDisk.MinimumCount, " and fetched count ", len(disk))
		mgr.Log.Info("", "Provided Maximum count for Disk", expectedDisk.MaximumCount, " and fetched count ", len(disk))

		if (expectedDisk.MinimumCount > len(disk)) || (expectedDisk.MaximumCount < len(disk)) {
			mgr.Log.Info("Disk MINMAX Count did not match")
			return false
		}

	} else if expectedDisk.MaximumCount > 0 {
		if expectedDisk.MaximumCount < len(disk) {
			mgr.Log.Info("Disk MAX Count did not match")
			return false
		}
	} else if expectedDisk.MinimumCount > 0 {
		if expectedDisk.MinimumCount > len(disk) {
			mgr.Log.Info("Disk MIN Count did not match")
			return false
		}

	}

	for _, disk := range disk {
		if expectedDisk.MaximumIndividualSizeGB > 0 && expectedDisk.MinimumIndividualSizeGB > 0 {

			mgr.Log.Info("", "Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeBytes)
			mgr.Log.Info("", "Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeBytes)
			if (bmh.Capacity(expectedDisk.MaximumIndividualSizeGB) < disk.SizeBytes) || (bmh.Capacity(expectedDisk.MinimumIndividualSizeGB) > disk.SizeBytes) {

				mgr.Log.Info("Disk MINMAX SIZE did not match")
				return false
			}
		} else if expectedDisk.MaximumIndividualSizeGB > 0 {

			mgr.Log.Info("", "Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeBytes)
			if bmh.Capacity(expectedDisk.MaximumIndividualSizeGB) < disk.SizeBytes {

				mgr.Log.Info("Disk MAX SIZE did not match")
				return false
			}
		} else if expectedDisk.MinimumIndividualSizeGB > 0 {

			mgr.Log.Info("", "Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeBytes)
			if bmh.Capacity(expectedDisk.MinimumIndividualSizeGB) > disk.SizeBytes {

				mgr.Log.Info("Disk MIN SIZE did not match")
				return false
			}
		}
	}

	return true
}
