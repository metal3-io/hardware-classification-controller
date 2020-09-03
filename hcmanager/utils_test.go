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

	bmoapis "github.com/metal3-io/baremetal-operator/pkg/apis"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HCManager", func() {

	type testCaseFetchBMH struct {
		namespace     string
		hosts         []runtime.Object
		expectedError bool
	}

	DescribeTable("Test fetch BaremetalHost list",
		func(tc testCaseFetchBMH) {
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), tc.hosts...)
			hcManager := NewHardwareClassificationManager(c, klogr.New())
			result, _, err := hcManager.FetchBmhHostList(tc.namespace)
			if tc.expectedError {
				Expect(len(result)).To(BeZero())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(result)).NotTo(BeZero())
			}
		},
		Entry("Should fetch BaremetalHosts in ready state and under metal3 namespace",
			testCaseFetchBMH{
				namespace:     getNamespace(),
				hosts:         getHosts(),
				expectedError: false},
		),
		Entry("Should return error while fetching BaremetalHosts",
			testCaseFetchBMH{
				namespace:     "sample",
				hosts:         getHosts(),
				expectedError: true},
		),
	)

	type testCaseValidation struct {
		hosts         []runtime.Object
		hwcProfile    hwcc.HardwareCharacteristics
		expectedError bool
	}

	DescribeTable("Test user profiles",
		func(tc testCaseValidation) {
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), tc.hosts...)
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
				hosts:         getHosts(),
				hwcProfile:    getExtractedHardwareProfile(),
				expectedError: false},
		),
		Entry("Test Empty Profile",
			testCaseValidation{
				hosts:         getHosts(),
				hwcProfile:    getEmptyProfile(),
				expectedError: true},
		),
		Entry("Test Invalid CPU Details",
			testCaseValidation{
				hosts:         getHosts(),
				hwcProfile:    getInvalidCPUProfile(),
				expectedError: true},
		),
		Entry("Test Invalid DISK Details",
			testCaseValidation{
				hosts:         getHosts(),
				hwcProfile:    getInvalidDiskProfile(),
				expectedError: true},
		),
		Entry("Test Invalid RAM Details",
			testCaseValidation{
				hosts:         getHosts(),
				hwcProfile:    getInvalidRAMProfile(),
				expectedError: true},
		),
		Entry("Test Invalid NIC Details",
			testCaseValidation{
				hosts:         getHosts(),
				hwcProfile:    getInvalidNicProfile(),
				expectedError: true},
		),
		Entry("Test Missing NIC Details",
			testCaseValidation{
				hosts:         getHosts(),
				hwcProfile:    getMissingNicDetails(),
				expectedError: false},
		),
	)
})

//setupSchemeMm Add the bmoapi to our scheme
func setupSchemeMm() *runtime.Scheme {
	s := runtime.NewScheme()
	if err := bmoapis.AddToScheme(s); err != nil {
		panic(err)
	}
	if err := hwcc.AddToScheme(s); err != nil {
		panic(err)
	}
	return s
}
