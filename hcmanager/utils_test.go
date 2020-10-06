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
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HCManager", func() {

	type testCaseFetchBMH struct {
		namespace      string
		expectedError  bool
		expectedResult []bmh.BareMetalHost
	}

	DescribeTable("Test fetch BaremetalHost list",
		func(tc testCaseFetchBMH) {
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), getHosts()...)
			hcManager := NewHardwareClassificationManager(c, klogr.New())
			result, _, err := hcManager.FetchBmhHostList(tc.namespace)
			if tc.expectedError {
				Expect(err).To(HaveOccurred())
			} else if len(tc.expectedResult) == 0 {
				Expect(len(result)).To(BeZero())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(result).Should(Equal(tc.expectedResult))
			}
		},
		Entry("Should fetch BaremetalHosts in ready state and under metal3 namespace",
			testCaseFetchBMH{
				namespace:      getNamespace(),
				expectedError:  false,
				expectedResult: getExpectedResult()},
		),
		Entry("Should return empty result while fetching BaremetalHosts",
			testCaseFetchBMH{
				namespace:      "sample",
				expectedError:  false,
				expectedResult: []bmh.BareMetalHost{}},
		),
	)

	type testCaseValidation struct {
		hwcProfile    hwcc.HardwareCharacteristics
		expectedError bool
	}

	DescribeTable("Test user profiles",
		func(tc testCaseValidation) {
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), getHosts()...)
			hcManager := NewHardwareClassificationManager(c, klogr.New())
			err := hcManager.ValidateExtractedHardwareProfile(tc.hwcProfile)
			if tc.expectedError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("Test Valid Profile",
			testCaseValidation{
				hwcProfile:    getExtractedHardwareProfile(),
				expectedError: false},
		),
		Entry("Test Empty Profile",
			testCaseValidation{
				hwcProfile:    getEmptyProfile(),
				expectedError: true},
		),
		Entry("Test Invalid CPU Details",
			testCaseValidation{
				hwcProfile:    getInvalidCPUProfile(),
				expectedError: true},
		),
		Entry("Test Invalid DISK Details",
			testCaseValidation{
				hwcProfile:    getInvalidDiskProfile(),
				expectedError: true},
		),
		Entry("Test Invalid RAM Details",
			testCaseValidation{
				hwcProfile:    getInvalidRAMProfile(),
				expectedError: true},
		),
		Entry("Test Invalid NIC Details",
			testCaseValidation{
				hwcProfile:    getInvalidNicProfile(),
				expectedError: true},
		),
		Entry("Test Missing NIC Details",
			testCaseValidation{
				hwcProfile:    getMissingNicDetails(),
				expectedError: false},
		),
	)
})

//setupSchemeMm Add the bmoapi to our scheme
func setupSchemeMm() *runtime.Scheme {
	s := runtime.NewScheme()
	if err := bmh.AddToScheme(s); err != nil {
		panic(err)
	}
	if err := hwcc.AddToScheme(s); err != nil {
		panic(err)
	}
	return s
}
