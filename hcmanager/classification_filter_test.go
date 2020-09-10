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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HCManager", func() {

	type testCaseVC struct {
		hwcProfile               hwcc.HardwareCharacteristics
		profileName              string
		bmHosts                  []bmh.BareMetalHost
		expectedHardwareDetails  []bmh.HardwareDetails
		expectedComparisonResult []string
	}

	DescribeTable("Test Classification and Validation Framework",
		func(tc testCaseVC) {
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), getHosts()...)
			hcManager := NewHardwareClassificationManager(c, klogr.New())

			validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(
				tc.hwcProfile, tc.bmHosts)

			if len(tc.expectedHardwareDetails) == 0 {
				Expect(len(validatedHardwareDetails)).To(BeZero())
			} else {
				Expect(validatedHardwareDetails).Should(Equal(tc.expectedHardwareDetails))

				comparedHost := hcManager.MinMaxFilter(
					tc.profileName,
					validatedHardwareDetails,
					tc.hwcProfile)

				if len(tc.expectedComparisonResult) == 0 {
					Expect(len(comparedHost)).To(BeZero())
				} else {
					Expect(comparedHost).Should(Equal(tc.expectedComparisonResult))
				}
			}
		},
		Entry("Validation and Comparison Should return Hardware Details List and Compared Host List",
			testCaseVC{
				hwcProfile:               getExtractedHardwareProfile(),
				profileName:              getTestProfileName(),
				bmHosts:                  getExpectedResult(),
				expectedHardwareDetails:  getExpectedHardwareDetails(),
				expectedComparisonResult: getExpectedComparedHost(),
			},
		),

		Entry("Validation Should Return Empty Hardware Details List",
			testCaseVC{
				hwcProfile:               getEmptyProfile(),
				profileName:              getTestProfileName(),
				bmHosts:                  getExpectedResult(),
				expectedHardwareDetails:  []bmh.HardwareDetails{},
				expectedComparisonResult: []string{},
			},
		),
		Entry("Comparison Should Return Empty Host List",
			testCaseVC{
				hwcProfile:               getMissingNicDetails(),
				profileName:              getTestProfileName(),
				bmHosts:                  getExpectedResult(),
				expectedHardwareDetails:  getExpectedMissingNICHardwareDetails(),
				expectedComparisonResult: []string{},
			},
		),
	)
})
