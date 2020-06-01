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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HardwareClassificationSpec defines the desired state of HardwareClassification
type HardwareClassificationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ExpectedHardwareConfiguration defines expected hardware configurations for CPU, RAM, Disk, NIC.
	ExpectedHardwareConfiguration ExpectedHardwareConfiguration `json:"expectedValidationConfiguration"`
}

// ExpectedHardwareConfiguration details to match with the host
type ExpectedHardwareConfiguration struct {
	// +optional
	CPU *CPU `json:"CPU,omitempty"`
	// +optional
	Disk *Disk `json:"Disk,omitempty"`
	// +optional
	NIC *NIC `json:"NIC,omitempty"`
	// +optional
	RAM *RAM `json:"RAM,omitempty"`
}

// CPU contains CPU details extracted from the hardware profile
type CPU struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumCount int `json:"minimumCount" description:"minimum cpu count, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumCount int `json:"maximumCount" description:"maximum cpu count, greater than 0"`
	// +optional
	// +kubebuilder:validation:Pattern=`^(0\.\d*[1-9]\d*|[1-9]\d*(\.\d+)?)$`
	MinimumSpeed string `json:"minimumSpeed" description:"minimum speed of cpu, greater than 0"`
	// +optional
	// +kubebuilder:validation:Pattern=`^(0\.\d*[1-9]\d*|[1-9]\d*(\.\d+)?)$`
	MaximumSpeed string `json:"maximumSpeed" description:"maximum speed of cpu, greater than 0"`
}

// Disk contains disk details extracted from the hardware profile
type Disk struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumCount int `json:"minimumCount" description:"minimum count of disk, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumIndividualSizeGB int64 `json:"minimumIndividualSizeGB" description:"minimum individual size of disk, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumCount int `json:"maximumCount" description:"maximum count of disk, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumIndividualSizeGB int64 `json:"maximumIndividualSizeGB" description:"maximum individual size of disk, greater than 0"`
}

// NIC contains nic details extracted from the hardware profile
type NIC struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumCount int `json:"minimumCount" description:"minimum count of nics, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumCount int `json:"maximumCount" description:"maximum count of nics, greater than 0"`
}

// RAM contains ram details extracted from the hardware profile
type RAM struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinimumSizeGB int `json:"minimumSizeGB" description:"minimun size of ram, greater than 0"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaximumSizeGB int `json:"maximumSizeGB" description:"maximum size of ram, greater than 0"`
}

// ProfileMatchStatus represents the state of the HardwareClassification
type ProfileMatchStatus string

const (
	// ProfileMatchStatusEmpty is the default status value
	ProfileMatchStatusEmpty ProfileMatchStatus = ""

	// ProfileMatchStatusMatched is the status value when the profile
	// matches to one of the BareMtalHost.
	ProfileMatchStatusMatched ProfileMatchStatus = "matched"

	// ProfileMatchStatusUnMatched is the status value when the profile
	// not matches to one of the BareMtalHost.
	ProfileMatchStatusUnMatched ProfileMatchStatus = "unmatched"
)

// ErrorType indicates the class of problem that has caused the HCC resource
// to enter an error state.
type ErrorType string

const (
	// LabelUpdateFailure is an error condition occurring when the
	// controller is unalble to update label of BareMetalHost.
	LabelUpdateFailure ErrorType = "label update error"

	// LabelDeleteFailure is an error condition occurring when the
	// controller is unalble to delete label of BareMetalHost.
	LabelDeleteFailure ErrorType = "label delete error"

	// FetchBMHListFailure is an error condition occurring when the
	// controller is unable to fetch BMH from BMO
	FetchBMHListFailure ErrorType = "fetch BMH from BMO error"

	// ProfileMisConfigured is an error condition occurring when the
	// extracted profile is empty.
	ProfileMisConfigured ErrorType = "Empty Profile Error"

	// NoBMHHost is an error condition occurring when the
	// baremetal host is empty.
	NoBMHHost ErrorType = "No baremetal host found"
)

// HardwareClassificationStatus defines the observed state of HardwareClassification
type HardwareClassificationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ErrorType indicates the type of failure encountered
	ErrorType ErrorType `json:"errorType,omitempty"`

	// ProfileMatchStatus identifies whether a applied profile is matches or not
	ProfileMatchStatus ProfileMatchStatus `json:"profileMatchStatus"`

	// The last error message reported by the hardwareclassification system
	ErrorMessage string `json:"errorMessage"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ProfileMatchStatus",type="string",JSONPath=".status.profileMatchStatus",description="Profile Match Status"
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.errorMessage",description="Most recent error"

// HardwareClassification is the Schema for the hardwareclassifications API
type HardwareClassification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HardwareClassificationSpec   `json:"spec,omitempty"`
	Status HardwareClassificationStatus `json:"status,omitempty"`
}

// SetProfileMatchStatus updates the ProfileMatchStatus field and returns
// true when a change is made or false when no change is made.
func (hcc *HardwareClassification) SetProfileMatchStatus(status ProfileMatchStatus) bool {
	if hcc.Status.ProfileMatchStatus != status {
		hcc.Status.ProfileMatchStatus = status
		return true
	}
	return false
}

// SetErrorMessage updates the ErrorMessage in the HardwareClassification Status struct
// when necessary and returns true when a change is made or false when
// no change is made.
func (hcc *HardwareClassification) SetErrorMessage(errType ErrorType, message string) (dirty bool) {
	if hcc.Status.ErrorType != errType {
		hcc.Status.ErrorType = errType
		dirty = true
	}
	if hcc.Status.ErrorMessage != message {
		hcc.Status.ErrorMessage = message
		dirty = true
	}
	return dirty
}

// ClearError removes any existing error message.
func (hcc *HardwareClassification) ClearError() (dirty bool) {
	var emptyErrType ErrorType = ""
	if hcc.Status.ErrorType != emptyErrType {
		hcc.Status.ErrorType = emptyErrType
		dirty = true
	}
	if hcc.Status.ErrorMessage != "" {
		hcc.Status.ErrorMessage = ""
		dirty = true
	}
	return dirty
}

// +kubebuilder:object:root=true

// HardwareClassificationList contains a list of HardwareClassification
type HardwareClassificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HardwareClassification `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HardwareClassification{}, &HardwareClassificationList{})
}
