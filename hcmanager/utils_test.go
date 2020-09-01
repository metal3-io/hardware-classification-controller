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

	DescribeTable("Test fetch BaremetalHost list",
		func(namespace string) {
			hostTest := getHosts()
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
			hcManager := NewHardwareClassificationManager(c, klogr.New())
			result, _, err := hcManager.FetchBmhHostList(namespace)
			if err != nil {
				Expect(len(result)).To(BeZero())
			}
			if len(result) > 0 {
				Expect(len(result)).Should(Equal(2))
			}
		},
		Entry("Should fetch BaremetalHosts in ready state and under metal3 namespace", getNamespace()),
		Entry("Should return error while fetching BaremetalHosts", "sample"),
	)

	DescribeTable("Test user profiles",
		func(hwcc hwcc.HardwareCharacteristics) {
			hostTest := getHosts()
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
			hcManager := NewHardwareClassificationManager(c, klogr.New())
			err := hcManager.ValidateExtractedHardwareProfile(hwcc)
			if err != nil {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).Should(BeNil())
			}

		},
		Entry("Test Valid Profile", getExtractedHardwareProfile()),
		Entry("Test Empty profile", getEmptyProfile()),
		Entry("Test Invalid CPU details", getInvalidCPUProfile()),
		Entry("Test Invalid DISK details", getInvalidDiskProfile()),
		Entry("Test Invalid RAM details", getInvalidRAMProfile()),
		Entry("Test Invalid NIC details", getInvalidNicProfile()),
		Entry("Test Missing NIC details", getMissingNicDetails()),
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
