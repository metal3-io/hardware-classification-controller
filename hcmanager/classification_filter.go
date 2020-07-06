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
func (mgr HardwareClassificationManager) MinMaxFilter(ProfileName string, HostList []bmh.HardwareDetails, expectedHardwareprofile hwcc.HardwareCharacteristics) []string {
	var validHost []string
	for _, hardwareDetail := range HostList {
		if !checkCPUCount(mgr, hardwareDetail.CPU, expectedHardwareprofile.Cpu) ||
			!checkRAM(mgr, hardwareDetail.RAMMebibytes, expectedHardwareprofile.Ram) ||
			!checkNICS(mgr, len(hardwareDetail.NIC), expectedHardwareprofile.Nic) ||
			!checkDiskDetails(mgr, hardwareDetail.Storage, expectedHardwareprofile.Disk) {
			continue
		}
		validHost = append(validHost, hardwareDetail.Hostname)
	}
	return validHost
}

//checkCPUCount this function checks the CPU details for both min and max parameters
func checkCPUCount(mgr HardwareClassificationManager, cpu bmh.CPU, expectedCPU *hwcc.Cpu) bool {
	if expectedCPU == nil {
		return true
	}
	if expectedCPU.MaximumCount > 0 {
		expectedMaxCPUCount := expectedCPU.MaximumCount
		mgr.Log.Info("", "Provided Maximum count for CPU", expectedMaxCPUCount, " and fetched count ", cpu.Count)
		if expectedMaxCPUCount < cpu.Count {
			mgr.Log.Info("CPU MAX COUNT did not match")
			return false
		}
	}
	if expectedCPU.MinimumCount > 0 {
		expectedMinCPUCount := expectedCPU.MinimumCount
		mgr.Log.Info("", "Provided Minimum count for CPU", expectedMinCPUCount, " and fetched count ", cpu.Count)
		if expectedMinCPUCount > cpu.Count {
			mgr.Log.Info("CPU MIN COUNT did not match")
			return false
		}
	}
	if expectedCPU.MaximumSpeedMHz > 0 {
		expectedMaxSpeedHz := bmh.ClockSpeed(expectedCPU.MaximumSpeedMHz)
		mgr.Log.Info("", "Provided Maximum ClockSpeed for CPU", expectedMaxSpeedHz, " and fetched ClockSpeed ", cpu.ClockMegahertz)
		if expectedMaxSpeedHz < cpu.ClockMegahertz {
			mgr.Log.Info("CPU MAX ClockSpeed did not match")
			return false
		}
	}
	if expectedCPU.MinimumSpeedMHz > 0 {
		expectedMinSpeedHz := bmh.ClockSpeed(expectedCPU.MinimumSpeedMHz)
		mgr.Log.Info("", "Provided Minimum ClockSpeed for CPU", expectedMinSpeedHz, " and fetched ClockSpeed ", cpu.ClockMegahertz)
		if expectedMinSpeedHz > cpu.ClockMegahertz {
			mgr.Log.Info("CPU MIN ClockSpeed did not match")
			return false
		}
	}
	return true
}

//checkNICS this function checks the nics details for both min and max parameters
func checkNICS(mgr HardwareClassificationManager, nics int, expectedNIC *hwcc.Nic) bool {
	if expectedNIC == nil {
		return true
	}
	if expectedNIC.MaximumCount > 0 {
		expectedMaxNicCount := expectedNIC.MaximumCount
		mgr.Log.Info("", "Provided Maximum count for NICS", expectedMaxNicCount, " and fetched count ", nics)
		if expectedMaxNicCount < nics {
			mgr.Log.Info("NICS MAX count did not match")
			return false
		}
	}
	if expectedNIC.MinimumCount > 0 {
		expectedMinNicCount := expectedNIC.MinimumCount
		mgr.Log.Info("", "Provided Minimum Count for NICS", expectedMinNicCount, " and fetched count ", nics)
		if expectedMinNicCount > nics {
			mgr.Log.Info("NICS MIN count did not match")
			return false
		}
	}
	return true
}

//checkRAM this function checks the ram details for both min and max parameters
func checkRAM(mgr HardwareClassificationManager, ram int, expectedRAM *hwcc.Ram) bool {
	if expectedRAM == nil {
		return true
	}
	if expectedRAM.MaximumSizeGB > 0 {
		expectedMaxRAM := expectedRAM.MaximumSizeGB
		mgr.Log.Info("", "Provided Maximum Size for RAM", expectedMaxRAM, " and fetched SIZE ", ram)
		if expectedMaxRAM < ram {
			mgr.Log.Info("RAM MAX SIZE did not match")
			return false
		}
	}
	if expectedRAM.MinimumSizeGB > 0 {
		expectedMinRAM := expectedRAM.MinimumSizeGB
		mgr.Log.Info("", "Provided Minimum Size for RAM", expectedMinRAM, " and fetched SIZE ", ram)
		if expectedMinRAM > ram {
			mgr.Log.Info("RAM MIN SIZE did not match")
			return false
		}
	}
	return true
}

//checkDiskDetails this function checks the Disk details for both min and max parameters
func checkDiskDetails(mgr HardwareClassificationManager, disks []bmh.Storage, expectedDisk *hwcc.Disk) bool {
	if expectedDisk == nil {
		return true
	}
	if expectedDisk.MaximumCount > 0 {
		expectedMaxDiskCount := expectedDisk.MaximumCount
		mgr.Log.Info("", "Provided Maximum count for Disk", expectedMaxDiskCount, " and fetched count ", len(disks))
		if expectedMaxDiskCount < len(disks) {
			mgr.Log.Info("Disk MAX Count did not match")
			return false
		}
	}
	if expectedDisk.MinimumCount > 0 {
		expectedMinDiskCount := expectedDisk.MinimumCount
		mgr.Log.Info("", "Provided Minimum count for Disk", expectedMinDiskCount, " and fetched count ", len(disks))
		if expectedMinDiskCount > len(disks) {
			mgr.Log.Info("Disk MIN Count did not match")
			return false
		}
	}
	for _, disk := range disks {
		if expectedDisk.MaximumIndividualSizeGB > 0 {
			expectedMaxDiskSize := bmh.Capacity(expectedDisk.MaximumIndividualSizeGB)
			mgr.Log.Info("", "Provided Maximum Size for Disk", expectedMaxDiskSize, " and fetched Size ", disk.SizeBytes)
			if expectedMaxDiskSize < disk.SizeBytes {
				mgr.Log.Info("Disk MAX SIZE did not match")
				return false
			}
		}
		if expectedDisk.MinimumIndividualSizeGB > 0 {
			expectedMinDiskSize := bmh.Capacity(expectedDisk.MinimumIndividualSizeGB)
			mgr.Log.Info("", "Provided Minimum Size for Disk", expectedMinDiskSize, " and fetched Size ", disk.SizeBytes)
			if expectedMinDiskSize > disk.SizeBytes {
				mgr.Log.Info("Disk MIN SIZE did not match")
				return false
			}
		}
	}
	return true
}
