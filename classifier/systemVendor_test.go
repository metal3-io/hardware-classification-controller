package classifier

import (
	"testing"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func TestSystemVendorDetails(t *testing.T) {
	assert.True(t, checkString("Dell Inc.", "Dell Inc."))
	assert.True(t, checkString("PowerEdge", "PowerEdge"))
	assert.True(t, checkString("", "Dell Inc."))
	assert.True(t, checkString("", "PowerEdge"))
	assert.False(t, checkString("Power", "PowerEdge"))
	assert.False(t, checkString("Dell", "Dell Inc."))

}

func TestSystemVendorManufacturer(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.SystemVendor
		Actual   string
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   "Dell Inc.",
			Expected: true,
		},
		{
			Scenario: "matched",
			Rule: &hwcc.SystemVendor{
				Manufacturer: "Dell Inc.",
			},
			Actual:   "Dell Inc.",
			Expected: true,
		},
		{
			Scenario: "unmatched",
			Rule: &hwcc.SystemVendor{
				Manufacturer: "Dell Inc.",
			},
			Actual:   "Dell",
			Expected: false,
		},
		{
			Scenario: "empty",
			Rule: &hwcc.SystemVendor{
				Manufacturer: "",
			},
			Actual:   "Dell Inc.",
			Expected: true,
		},
		{
			Scenario: "empty-host-details",
			Rule: &hwcc.SystemVendor{
				Manufacturer: "Dell Inc.",
			},
			Actual:   "",
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						SystemVendor: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						SystemVendor: bmh.HardwareSystemVendor{
							Manufacturer: tc.Actual,
						},
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host))
		})
	}
}

func TestCheckSystemVendorProductName(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.SystemVendor
		Actual   string
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   "PowerEdge",
			Expected: true,
		},
		{
			Scenario: "matched",
			Rule: &hwcc.SystemVendor{
				ProductName: "PowerEdge",
			},
			Actual:   "PowerEdge",
			Expected: true,
		},
		{
			Scenario: "unmatched",
			Rule: &hwcc.SystemVendor{
				ProductName: "PowerEdge",
			},
			Actual:   "Power",
			Expected: false,
		},
		{
			Scenario: "empty",
			Rule: &hwcc.SystemVendor{
				ProductName: "",
			},
			Actual:   "PowerEdge",
			Expected: true,
		},
		{
			Scenario: "empty-host-details",
			Rule: &hwcc.SystemVendor{
				ProductName: "PowerEdge",
			},
			Actual:   "",
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						SystemVendor: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						SystemVendor: bmh.HardwareSystemVendor{
							ProductName: tc.Actual,
						},
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host))
		})
	}
}
