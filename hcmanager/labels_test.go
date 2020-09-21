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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HCManager", func() {

	type testCaseLabel struct {
		obejctMeta               metav1.ObjectMeta
		validHosts               []string
		bmHosts                  bmh.BareMetalHostList
		expectedError            bool
		expectedUpdateLabelError []string
	}

	DescribeTable("Test Update Labels on BMHs",
		func(tc testCaseLabel) {
			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), getHosts()...)
			hcManager := NewHardwareClassificationManager(c, klogr.New())

			updateLabelError := hcManager.UpdateLabels(context.TODO(), tc.obejctMeta, tc.validHosts, tc.bmHosts)

			if tc.expectedError {
				Expect(len(updateLabelError)).To(Equal(len(tc.expectedUpdateLabelError)))
			} else {
				Expect(len(updateLabelError)).To(BeZero())
			}
		},
		Entry("Update labels on valid hosts",
			testCaseLabel{
				obejctMeta:               getObjectMeta(),
				validHosts:               getExpectedComparedHost(),
				bmHosts:                  getExpectedBMHList(),
				expectedError:            false,
				expectedUpdateLabelError: []string{},
			},
		),
		Entry("Delete label on unmatched host",
			testCaseLabel{
				obejctMeta:               getObjectMeta(),
				validHosts:               []string{"host-1"},
				bmHosts:                  getExpectedBMHList(),
				expectedError:            false,
				expectedUpdateLabelError: []string{},
			},
		),
	)
})
