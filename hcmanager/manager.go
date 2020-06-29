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

	"github.com/go-logr/logr"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HardwareClassificationManager only contains a client
type HardwareClassificationManager struct {
	client  client.Client
	Log     logr.Logger
	Profile *hwcc.HardwareClassification
}

// HardwareClassificationInterface important function used in reconciler
type HardwareClassificationInterface interface {
	FetchBmhHostList(namespace string) ([]bmh.BareMetalHost, bmh.BareMetalHostList, error)
	ExtractAndValidateHardwareDetails(hwcc.HardwareCharacteristics, []bmh.BareMetalHost) map[string]map[string]interface{}
	ValidateExtractedHardwareProfile(hwcc.HardwareCharacteristics) error
	MinMaxComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.HardwareCharacteristics) []string
}

//NewHardwareClassificationManager return new hardware classification manager
func NewHardwareClassificationManager(client client.Client, log logr.Logger) HardwareClassificationInterface {
	return HardwareClassificationManager{
		client: client,
		Log:    log,
	}

}
