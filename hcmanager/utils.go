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
	"context"
	"errors"
	"net"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	//LabelName initial name to the label key as hardware classification group
	LabelName = "hardwareclassification.metal3.io/"
	//Status extract the baremetal host for status ready
	Status = "ready"
	//DefaultLabel if label is missing from the Extracted Hardware Profile
	DefaultLabel = "matches"
	//CPULabel label for extraction of hardware details
	CPULabel = "CPU"
	//NICLabel label for extraction of hardware details
	NICLabel = "NIC"
	//DISKLabel label for extraction of hardware details
	DISKLabel = "DISK"
	//RAMLabel label for extraction of hardware details
	RAMLabel = "RAM"
)

//FetchBmhHostList this function will fetch and return baremetal hosts in ready state
func (mgr HardwareClassificationManager) FetchBmhHostList(Namespace string) ([]bmh.BareMetalHost, bmh.BareMetalHostList, error) {
	ctx := context.Background()
	bmhHostList := bmh.BareMetalHostList{}
	var validHostList []bmh.BareMetalHost
	opts := &client.ListOptions{
		Namespace: Namespace,
	}
	// Get list of BareMetalHost from BMO
	err := mgr.client.List(ctx, &bmhHostList, opts)
	if err != nil {
		return validHostList, bmhHostList, errors.New(err.Error())
	}
	// Get hosts in ready status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.Provisioning.State == "ready" {
			validHostList = append(validHostList, host)
		}
	}
	return validHostList, bmhHostList, nil
}

//CheckValidIP uses net package to check if the IP is valid or not
func CheckValidIP(NICIp string) bool {
	return net.ParseIP(NICIp) != nil
}

//ConvertBytesToGb it converts the Byte into GB
func ConvertBytesToGb(inBytes bmh.Capacity) bmh.Capacity {
	inGb := (inBytes / 1024 / 1024 / 1024)
	return inGb
}

//SetStatus set error status for the hardware classification profile
func SetStatus(hwc *hwcc.HardwareClassification,
	status hwcc.ProfileMatchStatus,
	errorType hwcc.ErrorType, errorMessage string) {
	hwc.Status.ProfileMatchStatus = status
	hwc.Status.ErrorType = errorType
	hwc.Status.ErrorMessage = errorMessage
}

//ValidateExtractedHardwareProfile it will validate the extracted hardware profile and log the warnings
func (mgr HardwareClassificationManager) ValidateExtractedHardwareProfile(extractedProfile hwcc.HardwareCharacteristics) error {
	if (extractedProfile.Cpu == nil) &&
		(extractedProfile.Ram == nil) &&
		(extractedProfile.Disk == nil) &&
		(extractedProfile.Nic == nil) {
		return errors.New("Hardware profile details can not be empty")
	}

	if extractedProfile.Cpu == nil {
		mgr.Log.Info("CPU details are not present")
	} else {
		if extractedProfile.Cpu.MaximumCount == 0 &&
			extractedProfile.Cpu.MinimumCount == 0 &&
			extractedProfile.Cpu.MaximumSpeedMHz == 0 &&
			extractedProfile.Cpu.MinimumSpeedMHz == 0 {
			return errors.New("Invalid CPU Details")
		}
	}

	if extractedProfile.Ram == nil {
		mgr.Log.Info("RAM details are not present")
	} else {
		if extractedProfile.Ram.MaximumSizeGB == 0 &&
			extractedProfile.Ram.MinimumSizeGB == 0 {
			return errors.New("Invalid RAM Details")
		}
	}

	if extractedProfile.Disk == nil {
		mgr.Log.Info("DISK details are not present")
	} else {
		if extractedProfile.Disk.MaximumCount == 0 &&
			extractedProfile.Disk.MinimumCount == 0 &&
			extractedProfile.Disk.MaximumIndividualSizeGB == 0 &&
			extractedProfile.Disk.MinimumIndividualSizeGB == 0 {
			return errors.New("Invalid DISK Details")
		}
	}

	if extractedProfile.Nic == nil {
		mgr.Log.Info("NIC details is not present")
	} else {
		if extractedProfile.Nic.MaximumCount == 0 &&
			extractedProfile.Nic.MinimumCount == 0 {
			return errors.New("Invalid NIC details")
		}
	}
	return nil
}
