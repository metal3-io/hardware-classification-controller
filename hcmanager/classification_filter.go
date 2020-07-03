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

// MinMaxFilter it will perform the minimum and maximum comparison based on the value provided by the user and check for the valid host
func (mgr HardwareClassificationManager) MinMaxFilter(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.HardwareCharacteristics) []string {
	var comparedHost []string
	for hostname, details := range validatedHost {
		isHostValid := true
		for _, value := range details {
			isValid := false
			if cpu, CPUOK := value.(bmh.CPU); CPUOK {
				if checkCPUCount(mgr, cpu, expectedHardwareprofile.Cpu) {
					isValid = true
				}
			}
			if ram, RAMOK := value.(int64); RAMOK {
				if checkRAM(mgr, ram, expectedHardwareprofile.Ram) {
					isValid = true
				}
			}
			if nics, NICSOK := value.(int); NICSOK {
				if checkNICS(mgr, nics, expectedHardwareprofile.Nic) {
					isValid = true
				}
			}
			if disks, DISKOK := value.([]bmh.Storage); DISKOK {
				if checkDiskDetails(mgr, disks, expectedHardwareprofile.Disk) {
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
	expectedMinCPUCount := expectedCPU.MinimumCount
	expectedMaxCPUCount := expectedCPU.MaximumCount
	expectedMinSpeedHz := bmh.ClockSpeed(expectedCPU.MinimumSpeedMHz)
	expectedMaxSpeedHz := bmh.ClockSpeed(expectedCPU.MaximumSpeedMHz)

	mgr.Log.Info("", "Provided Minimum count for CPU", expectedMinCPUCount, " and fetched count ", cpu.Count)
	mgr.Log.Info("", "Provided Maximum count for CPU", expectedMaxCPUCount, " and fetched count ", cpu.Count)

	if expectedMaxCPUCount > 0 && expectedMinCPUCount > 0 {
		if (expectedMinCPUCount > cpu.Count) || (expectedMaxCPUCount < cpu.Count) {
			mgr.Log.Info("CPU MINMAX COUNT did not match")
			return false
		}
	} else if expectedMaxCPUCount > 0 {
		if expectedMaxCPUCount < cpu.Count {
			mgr.Log.Info("CPU MAX COUNT did not match")
			return false
		}
	} else if expectedMinCPUCount > 0 {
		if expectedMinCPUCount > cpu.Count {
			mgr.Log.Info("CPU MIN COUNT did not match")
			return false
		}
	}

	mgr.Log.Info("", "Provided Minimum ClockSpeed for CPU", expectedMinSpeedHz, " and fetched ClockSpeed ", cpu.ClockMegahertz)
	mgr.Log.Info("", "Provided Maximum ClockSpeed for CPU", expectedMaxSpeedHz, " and fetched ClockSpeed ", cpu.ClockMegahertz)

	if expectedMaxSpeedHz > 0 && expectedMinSpeedHz > 0 {
		if expectedMinSpeedHz > cpu.ClockMegahertz || expectedMaxSpeedHz < cpu.ClockMegahertz {
			mgr.Log.Info("CPU MINMAX ClockSpeed did not match")
			return false
		}
	} else if expectedMaxSpeedHz > 0 {
		if expectedMaxSpeedHz < cpu.ClockMegahertz {
			mgr.Log.Info("CPU MAX ClockSpeed did not match")
			return false
		}
	} else if expectedMinSpeedHz > 0 {
		if expectedMinSpeedHz > cpu.ClockMegahertz {
			mgr.Log.Info("CPU MIN ClockSpeed did not match")
			return false
		}
	}
	return true
}

//checkNICS this function checks the nics details for both min and max parameters
func checkNICS(mgr HardwareClassificationManager, nics int, expectedNIC *hwcc.Nic) bool {
	expectedMinNicCount := expectedNIC.MinimumCount
	expectedMaxNicCount := expectedNIC.MaximumCount

	mgr.Log.Info("", "Provided Minimum Count for NICS", expectedMinNicCount, " and fetched count ", nics)
	mgr.Log.Info("", "Provided Maximum count for NICS", expectedMaxNicCount, " and fetched count ", nics)

	if expectedMaxNicCount > 0 && expectedMinNicCount > 0 {
		if expectedMinNicCount > nics || expectedMaxNicCount < nics {
			mgr.Log.Info("NICS MINMAX count did not match")
			return false
		}
	} else if expectedMaxNicCount > 0 {
		if expectedMaxNicCount < nics {
			mgr.Log.Info("NICS MAX count did not match")
			return false
		}
	} else if expectedMinNicCount > 0 {
		if expectedMinNicCount > nics {
			mgr.Log.Info("NICS MIN count did not match")
			return false
		}
	}
	return true
}

//checkRAM this function checks the ram details for both min and max parameters
func checkRAM(mgr HardwareClassificationManager, ram int64, expectedRAM *hwcc.Ram) bool {
	expectedMinRAM := int64(expectedRAM.MinimumSizeGB)
	expectedMaxRAM := int64(expectedRAM.MaximumSizeGB)

	mgr.Log.Info("", "Provided Minimum Size for RAM", expectedMinRAM, " and fetched SIZE ", ram)
	mgr.Log.Info("", "Provided Maximum Size for RAM", expectedMaxRAM, " and fetched SIZE ", ram)

	if expectedMaxRAM > 0 && expectedMinRAM > 0 {
		if expectedMinRAM > ram || expectedMaxRAM < ram {
			mgr.Log.Info("RAM MINMAX SIZE did not match")
			return false
		}
	} else if expectedMaxRAM > 0 {
		if expectedMaxRAM < ram {
			mgr.Log.Info("RAM MAX SIZE did not match")
			return false
		}
	} else if expectedMinRAM > 0 {
		if expectedMinRAM > ram {
			mgr.Log.Info("RAM MIN SIZE did not match")
			return false
		}
	}
	return true
}

//checkDiskDetails this function checks the Disk details for both min and max parameters
func checkDiskDetails(mgr HardwareClassificationManager, disks []bmh.Storage, expectedDisk *hwcc.Disk) bool {
	expectedMaxDiskCount := expectedDisk.MaximumCount
	expectedMinDiskCount := expectedDisk.MinimumCount
	expectedMaxDiskSize := bmh.Capacity(expectedDisk.MaximumIndividualSizeGB)
	expectedMinDiskSize := bmh.Capacity(expectedDisk.MinimumIndividualSizeGB)

	mgr.Log.Info("", "Provided Minimum count for Disk", expectedMinDiskCount, " and fetched count ", len(disks))
	mgr.Log.Info("", "Provided Maximum count for Disk", expectedMaxDiskCount, " and fetched count ", len(disks))

	if expectedMaxDiskCount > 0 && expectedMinDiskCount > 0 {
		if expectedMinDiskCount > len(disks) || (expectedMaxDiskCount < len(disks)) {
			mgr.Log.Info("Disk MINMAX Count did not match")
			return false
		}
	} else if expectedMaxDiskCount > 0 {
		if expectedMaxDiskCount < len(disks) {
			mgr.Log.Info("Disk MAX Count did not match")
			return false
		}
	} else if expectedMinDiskCount > 0 {
		if expectedMinDiskCount > len(disks) {
			mgr.Log.Info("Disk MIN Count did not match")
			return false
		}
	}
	for _, disk := range disks {
		mgr.Log.Info("", "Provided Minimum Size for Disk", expectedMinDiskSize, " and fetched Size ", disk.SizeBytes)
		mgr.Log.Info("", "Provided Maximum Size for Disk", expectedMaxDiskSize, " and fetched Size ", disk.SizeBytes)
		if expectedMaxDiskSize > 0 && expectedMinDiskSize > 0 {
			if expectedMaxDiskSize < disk.SizeBytes || expectedMinDiskSize > disk.SizeBytes {
				mgr.Log.Info("Disk MINMAX SIZE did not match")
				return false
			}
		} else if expectedMaxDiskSize > 0 {
			if expectedMaxDiskSize < disk.SizeBytes {
				mgr.Log.Info("Disk MAX SIZE did not match")
				return false
			}
		} else if expectedMinDiskSize > 0 {
			if expectedMinDiskSize > disk.SizeBytes {
				mgr.Log.Info("Disk MIN SIZE did not match")
				return false
			}
		}
	}
	return true
}
