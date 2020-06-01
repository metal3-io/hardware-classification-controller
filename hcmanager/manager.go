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
	ExtractAndValidateHardwareDetails(hwcc.ExpectedHardwareConfiguration, []bmh.BareMetalHost) map[string]map[string]interface{}
	ValidateExtractedHardwareProfile(hwcc.ExpectedHardwareConfiguration) error
}

//NewHardwareClassificationManager return new hardware classification manager
func NewHardwareClassificationManager(client client.Client, log logr.Logger) HardwareClassificationInterface {
	return HardwareClassificationManager{
		client: client,
		Log:    log,
	}

}
