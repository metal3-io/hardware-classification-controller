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

	"github.com/go-logr/logr"
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	ExtractAndValidateHardwareDetails(hwcc.HardwareCharacteristics, []bmh.BareMetalHost) []bmh.HardwareDetails
	ValidateExtractedHardwareProfile(hwcc.HardwareCharacteristics) error
	MinMaxFilter(ProfileName string, HostList []bmh.HardwareDetails, expectedHardwareprofile hwcc.HardwareCharacteristics) []string
	UpdateLabels(ctx context.Context, hcMetaData v1.ObjectMeta, comparedHost []string, BMHList bmh.BareMetalHostList) []string
}

//NewHardwareClassificationManager return new hardware classification manager
func NewHardwareClassificationManager(client client.Client, log logr.Logger) HardwareClassificationInterface {
	return HardwareClassificationManager{
		client: client,
		Log:    log,
	}
}
