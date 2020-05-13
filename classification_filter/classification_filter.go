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

package filter

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"
	valTypes "hardware-classification-controller/validation/validationModel"
	"strconv"
)

// MinMaxComparison it will compare the minimum and maximum comparison based on the value provided by the user and check for the valid host
func MinMaxComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) []string {

	fmt.Printf("\n\n\n")
	var comparedHost []string

	for hostname, details := range validatedHost {
		isHostValid := true

		for _, value := range details {

			isValid := false

			cpu, CPUOK := value.(valTypes.CPU)
			if CPUOK {
				if checkCPUCount(cpu, expectedHardwareprofile.CPU) {
					isValid = true
				}
			}

			ram, RAMOK := value.(valTypes.RAM)
			if RAMOK {
				if checkRAM(ram, expectedHardwareprofile.RAM) {
					isValid = true
				}
			}

			nics, NICSOK := value.(valTypes.NIC)
			if NICSOK {
				if checkNICS(nics, expectedHardwareprofile.NIC) {
					isValid = true
				}
			}

			disk, DISKOK := value.(valTypes.Storage)
			if DISKOK {
				if checkDiskDetailsl(disk, expectedHardwareprofile.Disk) {
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
			fmt.Println(hostname, " Matched profile ", ProfileName)
			fmt.Printf("\n\n\n")

		} else {
			fmt.Println(hostname, " Did not match profile ", ProfileName)
			fmt.Printf("\n\n\n")

		}

	}

	return comparedHost

}

//checkNICS this function checks the nics details for both min and max parameters
func checkNICS(nics valTypes.NIC, expectedNIC *hwcc.NIC) bool {
	fmt.Printf("\n")
	if (expectedNIC.MaximumCount > 0) && (expectedNIC.MinimumCount > 0) {
		fmt.Println("Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics.Count)
		fmt.Println("Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics.Count)
		if (expectedNIC.MinimumCount > nics.Count) || (expectedNIC.MaximumCount < nics.Count) {
			fmt.Println("NICS MINMAX count did not match")
			return false
		}
	} else if expectedNIC.MaximumCount > 0 {
		fmt.Println("Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics.Count)
		if expectedNIC.MaximumCount < nics.Count {
			fmt.Println("NICS MAX count did not match")
			return false
		}

	} else if expectedNIC.MinimumCount > 0 {
		fmt.Println("Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics.Count)
		if expectedNIC.MinimumCount > nics.Count {
			fmt.Println("NICS MIN count did not match")
			return false
		}

	}
	return true
}

//checkRAM this function checks the ram details for both min and max parameters
func checkRAM(ram valTypes.RAM, expectedRAM *hwcc.RAM) bool {
	fmt.Printf("\n")
	if (expectedRAM.MaximumSizeGB > 0) && (expectedRAM.MinimumSizeGB > 0) {
		fmt.Println("Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram.RAMGb)
		fmt.Println("Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if (expectedRAM.MinimumSizeGB > ram.RAMGb) || (expectedRAM.MaximumSizeGB < ram.RAMGb) {
			fmt.Println("RAM MINMAX SIZE did not match")
			return false
		}
	} else if expectedRAM.MaximumSizeGB > 0 {
		fmt.Println("Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if expectedRAM.MaximumSizeGB < ram.RAMGb {
			fmt.Println("RAM MAX SIZE did not match")
			return false
		}

	} else if expectedRAM.MinimumSizeGB > 0 {
		fmt.Println("Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if expectedRAM.MinimumSizeGB > ram.RAMGb {
			fmt.Println("RAM MIN SIZE did not match")
			return false
		}

	}
	return true
}

//checkCPUCount this function checks the CPU details for both min and max parameters
func checkCPUCount(cpu valTypes.CPU, expectedCPU *hwcc.CPU) bool {

	fmt.Printf("\n")

	if (expectedCPU.MaximumCount > 0) && (expectedCPU.MinimumCount > 0) {
		fmt.Println("Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		fmt.Println("Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if (expectedCPU.MinimumCount > cpu.Count) || (expectedCPU.MaximumCount < cpu.Count) {
			fmt.Println("CPU MINMAX COUNT did not match")
			return false
		}

	} else if expectedCPU.MaximumCount > 0 {
		fmt.Println("Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MaximumCount < cpu.Count {
			fmt.Println("CPU MAX COUNT did not match")
			return false
		}

	} else if expectedCPU.MinimumCount > 0 {
		fmt.Println("Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MinimumCount > cpu.Count {
			fmt.Println("CPU MIN COUNT did not match")
			return false
		}

	}

	if (expectedCPU.MaximumSpeed != "") && (expectedCPU.MinimumSpeed != "") {
		MaxSpeed, errMax := strconv.ParseFloat(expectedCPU.MaximumSpeed, 64)
		MinSpeed, errMin := strconv.ParseFloat(expectedCPU.MinimumSpeed, 64)
		if errMax == nil && errMin == nil {
			fmt.Println("Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			fmt.Println("Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			if MinSpeed > 0 && MaxSpeed > 0 {
				if (MinSpeed > cpu.ClockSpeed) || (MaxSpeed < cpu.ClockSpeed) {
					fmt.Println("CPU MINMAX ClockSpeed did not match")
					return false
				}

			}
		}

	} else if expectedCPU.MaximumSpeed != "" {
		MaxSpeed, errMax := strconv.ParseFloat(expectedCPU.MaximumSpeed, 64)
		if errMax == nil {
			fmt.Println("Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			if MaxSpeed > 0 {
				if MaxSpeed < cpu.ClockSpeed {
					fmt.Println("CPU MAX ClockSpeed did not match")
					return false
				}

			}
		}
	} else if expectedCPU.MinimumSpeed != "" {
		MinSpeed, errMin := strconv.ParseFloat(expectedCPU.MinimumSpeed, 64)
		if errMin == nil {
			fmt.Println("Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
			if MinSpeed > 0 {
				if MinSpeed > cpu.ClockSpeed {
					fmt.Println("CPU MIN ClockSpeed did not match")
					return false
				}

			}
		}
	}

	return true

}

//checkDiskDetailsl this function checks the Disk details for both min and max parameters
func checkDiskDetailsl(storage valTypes.Storage, expectedDisk *hwcc.Disk) bool {
	fmt.Printf("\n")
	if (expectedDisk.MaximumCount > 0) && (expectedDisk.MinimumCount > 0) {
		fmt.Println("Provided Minimum count for Disk", expectedDisk.MinimumCount, " and fetched count ", storage.Count)
		fmt.Println("Provided Maximum count for Disk", expectedDisk.MaximumCount, " and fetched count ", storage.Count)
		if (expectedDisk.MaximumCount >= storage.Count) && (expectedDisk.MinimumCount <= storage.Count) {
			for _, disk := range storage.Disk {

				if expectedDisk.MaximumIndividualSizeGB > 0 && expectedDisk.MinimumIndividualSizeGB > 0 {
					fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if (expectedDisk.MaximumIndividualSizeGB < disk.SizeGb) || (expectedDisk.MinimumIndividualSizeGB > disk.SizeGb) {
						fmt.Println("Disk MINMAX SIZE did not match")
						return false
					}
				} else if expectedDisk.MaximumIndividualSizeGB > 0 {
					fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb {
						fmt.Println("Disk MAX SIZE did not match")
						return false
					}
				} else if expectedDisk.MinimumIndividualSizeGB > 0 {
					fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {
						fmt.Println("Disk MIN SIZE did not match")
						return false
					}
				}

			}

		} else {
			fmt.Println("Disk MINMAX Count did not match")
			return false
		}
	} else if expectedDisk.MaximumCount > 0 {
		fmt.Println("Provided Maximum count for Disk", expectedDisk.MaximumCount, " and fetched count ", storage.Count)
		if expectedDisk.MaximumCount >= storage.Count {
			for _, disk := range storage.Disk {

				if expectedDisk.MaximumIndividualSizeGB > 0 && expectedDisk.MinimumIndividualSizeGB > 0 {
					fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if (expectedDisk.MaximumIndividualSizeGB < disk.SizeGb) || (expectedDisk.MinimumIndividualSizeGB > disk.SizeGb) {
						fmt.Println("Disk MINMAX SIZE did not match")
						return false
					}
				} else if expectedDisk.MaximumIndividualSizeGB > 0 {
					fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb {
						fmt.Println("Disk MAX SIZE did not match")
						return false
					}
				} else if expectedDisk.MinimumIndividualSizeGB > 0 {
					fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {
						fmt.Println("Disk MIN SIZE did not match")
						return false
					}
				}

			}
		} else {
			fmt.Println("Disk MAX Count did not match")
			return false
		}
	} else if expectedDisk.MinimumCount > 0 {
		fmt.Println("Provided Minimum count for Disk", expectedDisk.MinimumCount, " and fetched count ", storage.Count)
		if expectedDisk.MinimumCount <= storage.Count {
			for _, disk := range storage.Disk {

				if expectedDisk.MaximumIndividualSizeGB > 0 && expectedDisk.MinimumIndividualSizeGB > 0 {
					fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if (expectedDisk.MaximumIndividualSizeGB < disk.SizeGb) || (expectedDisk.MinimumIndividualSizeGB > disk.SizeGb) {
						fmt.Println("Disk MINMAX SIZE did not match")
						return false
					}
				} else if expectedDisk.MaximumIndividualSizeGB > 0 {
					fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb {
						fmt.Println("Disk MAX SIZE did not match")
						return false
					}
				} else if expectedDisk.MinimumIndividualSizeGB > 0 {
					fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
					if expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {
						fmt.Println("Disk MIN SIZE did not match")
						return false
					}
				}
			}
		} else {
			fmt.Println("Disk MIN Count did not match")
			return false
		}
	}

	return true
}
