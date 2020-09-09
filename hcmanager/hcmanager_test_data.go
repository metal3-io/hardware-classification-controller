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
	"k8s.io/apimachinery/pkg/runtime"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getHost() bmh.BareMetalHost {
	host := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host-0",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateReady,
			}, HardwareDetails: &bmh.HardwareDetails{
				CPU:      bmh.CPU{Arch: "x86_64", Model: "Intel(R) Xeon(R) Gold 6226 CPU @ 2.70GHz", Count: 40, ClockMegahertz: 3400},
				Firmware: bmh.Firmware{BIOS: bmh.BIOS{Date: "", Vendor: "", Version: ""}},
				Hostname: "localhost.localdomain",
				NIC: []bmh.NIC{{IP: "", MAC: "xx:xx:xx:xx:xx:xx", Model: "0xXXXX 0xXXXX", Name: "eth11", PXE: false, SpeedGbps: 0, VLANID: 0},
					{IP: "192.xxx.xxx.xx", MAC: "xx:xx:xx:xx:xx:xx", Model: "0xXXXX 0xXXXX", Name: "eth6", PXE: true, SpeedGbps: 0, VLANID: 0}},
				RAMMebibytes: 196608,
				Storage: []bmh.Storage{{Name: "/dev/sda", SizeBytes: 599550590976},
					{Name: "/dev/sdb", SizeBytes: 599550590976},
					{Name: "/dev/sdc", SizeBytes: 599550590976},
					{Name: "/dev/sdd", SizeBytes: 599550590976},
					{Name: "/dev/sde", SizeBytes: 599550590976},
					{Name: "/dev/sdf", SizeBytes: 599550590976},
					{Name: "/dev/sdg", SizeBytes: 599550590976},
					{Name: "/dev/sdh", SizeBytes: 599550590976},
					{Name: "/dev/sdi", SizeBytes: 599550590976}},
				SystemVendor: bmh.HardwareSystemVendor{Manufacturer: "Dell Inc.", ProductName: "PowerEdge R740xd (SKU=NotProvided;ModelName=PowerEdge R740xd)"},
			},
		},
	}
	return host
}

func getNamespace() string {
	return "metal3"
}

func getExtractedHardwareProfile() hwcc.HardwareCharacteristics {

	return hwcc.HardwareCharacteristics{
		Cpu: &hwcc.Cpu{
			MaximumCount:    32,
			MinimumCount:    32,
			MaximumSpeedMHz: 5000,
			MinimumSpeedMHz: 3000,
		},
		Disk: &hwcc.Disk{
			MaximumCount:            9,
			MaximumIndividualSizeGB: 558,
			MinimumCount:            9,
			MinimumIndividualSizeGB: 558,
		},
		Nic: &hwcc.Nic{
			MaximumCount: 4,
			MinimumCount: 4,
		},
		Ram: &hwcc.Ram{
			MaximumSizeGB: 192,
			MinimumSizeGB: 192,
		},
	}

}

func getMissingNicDetails() hwcc.HardwareCharacteristics {

	return hwcc.HardwareCharacteristics{
		Cpu: &hwcc.Cpu{
			MaximumCount: 32,
			MinimumCount: 32,
		},
		Disk: &hwcc.Disk{
			MaximumCount:            9,
			MaximumIndividualSizeGB: 558,
			MinimumCount:            9,
			MinimumIndividualSizeGB: 558,
		},
		Ram: &hwcc.Ram{
			MaximumSizeGB: 192,
			MinimumSizeGB: 192,
		},
	}

}

func getEmptyProfile() hwcc.HardwareCharacteristics {
	return hwcc.HardwareCharacteristics{}
}

func getInvalidCPUProfile() hwcc.HardwareCharacteristics {

	return hwcc.HardwareCharacteristics{
		Cpu: &hwcc.Cpu{
			MaximumCount:    0,
			MinimumCount:    0,
			MaximumSpeedMHz: 0,
			MinimumSpeedMHz: 0,
		},
	}
}

func getInvalidDiskProfile() hwcc.HardwareCharacteristics {

	return hwcc.HardwareCharacteristics{
		Disk: &hwcc.Disk{
			MaximumCount:            0,
			MaximumIndividualSizeGB: 0,
			MinimumCount:            0,
			MinimumIndividualSizeGB: 0,
		},
	}
}

func getInvalidNicProfile() hwcc.HardwareCharacteristics {

	return hwcc.HardwareCharacteristics{
		Nic: &hwcc.Nic{
			MaximumCount: 0,
			MinimumCount: 0,
		},
	}
}

func getInvalidRAMProfile() hwcc.HardwareCharacteristics {

	return hwcc.HardwareCharacteristics{
		Ram: &hwcc.Ram{
			MaximumSizeGB: 0,
			MinimumSizeGB: 0,
		},
	}
}

func getObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{Name: "hardwareclassification-test"}
}

func getExtractedHardwareProfileRuntime() []runtime.Object {
	expectedHardwareClassification := hwcc.HardwareCharacteristics{
		Cpu: &hwcc.Cpu{
			MaximumCount:    1,
			MinimumCount:    1,
			MaximumSpeedMHz: 5000,
			MinimumSpeedMHz: 3000,
		},
		Disk: &hwcc.Disk{
			MaximumCount:            1,
			MaximumIndividualSizeGB: 500,
			MinimumCount:            1,
			MinimumIndividualSizeGB: 500,
		},
		Nic: &hwcc.Nic{
			MaximumCount: 1,
			MinimumCount: 2,
		},
		Ram: &hwcc.Ram{
			MaximumSizeGB: 8,
			MinimumSizeGB: 8,
		},
	}

	expectedHardwareConfiguration := hwcc.HardwareClassification{
		Spec: hwcc.HardwareClassificationSpec{
			HardwareCharacteristics: expectedHardwareClassification,
		},
	}

	return []runtime.Object{&expectedHardwareConfiguration}
}

func getExpectedResult() []bmh.BareMetalHost {
	host := getHost()
	host3 := host
	host3.ObjectMeta.Name = "host-3"
	return []bmh.BareMetalHost{host, host3}
}

func getHosts() []runtime.Object {
	host := getHost()
	host1 := host
	host1.ObjectMeta.Name = "host-1"
	host1.Status.Provisioning = bmh.ProvisionStatus{
		State: bmh.StateInspecting,
	}

	host2 := host
	host2.ObjectMeta.Name = "host-2"
	host2.ObjectMeta.Namespace = "test"
	host1.Status.Provisioning = bmh.ProvisionStatus{
		State: bmh.StateInspecting,
	}

	host3 := host
	host3.ObjectMeta.Name = "host-3"

	return []runtime.Object{&host, &host1, &host2, &host3}
}
