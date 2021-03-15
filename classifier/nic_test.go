package classifier

import (
	"fmt"
	"testing"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func TestCheckNICCount(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Nic
		Actual   int
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "no-min-max",
			Rule: &hwcc.Nic{
				MinimumCount: 0,
				MaximumCount: 0,
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "within-max",
			Rule: &hwcc.Nic{
				MinimumCount: 0,
				MaximumCount: 4,
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "within-min",
			Rule: &hwcc.Nic{
				MinimumCount: 1,
				MaximumCount: 0,
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "under-min",
			Rule: &hwcc.Nic{
				MinimumCount: 4,
				MaximumCount: 0,
			},
			Actual:   2,
			Expected: false,
		},
		{
			Scenario: "over-max",
			Rule: &hwcc.Nic{
				MinimumCount: 0,
				MaximumCount: 1,
			},
			Actual:   2,
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Nic: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{},
				},
			}
			nics := []bmh.NIC{}
			for i := 0; i < tc.Actual; i++ {
				nics = append(nics, bmh.NIC{Name: fmt.Sprintf("eth%d", i)})
			}
			host.Status.HardwareDetails.NIC = nics
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host),
				fmt.Sprintf("rule=%v actual=%v", tc.Rule, tc.Actual))
		})
	}
}

func TestCheckNICVendor(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Nic
		Actual   int
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "under-min",
			Rule: &hwcc.Nic{
				MinimumCount: 4,
				MaximumCount: 0,
			},
			Actual:   2,
			Expected: false,
		},
		{
			Scenario: "only min & vendor",
			Rule: &hwcc.Nic{
				MinimumCount: 1,
				NicSelector:  hwcc.NicSelector{Vendor: []string{"0x1af4"}},
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "only max & vendor",
			Rule: &hwcc.Nic{
				MaximumCount: 1,
				NicSelector:  hwcc.NicSelector{Vendor: []string{"0x1af4"}},
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "min, max & vendor",
			Rule: &hwcc.Nic{
				MinimumCount: 1,
				MaximumCount: 2,
				NicSelector:  hwcc.NicSelector{Vendor: []string{"0x1af4"}},
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "mismatch vendors",
			Rule: &hwcc.Nic{
				MinimumCount: 1,
				MaximumCount: 2,
				NicSelector:  hwcc.NicSelector{Vendor: []string{"0x8086", "0x1af4"}},
			},
			Actual:   2,
			Expected: false,
		},
		{
			Scenario: "match vendors",
			Rule: &hwcc.Nic{
				MinimumCount: 1,
				MaximumCount: 2,
				NicSelector:  hwcc.NicSelector{Vendor: []string{"0x1af4", "0x1af4"}},
			},
			Actual:   2,
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Nic: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{},
				},
			}
			nics := []bmh.NIC{}
			for i := 0; i < tc.Actual; i++ {
				nics = append(nics, bmh.NIC{Name: fmt.Sprintf("eth%d", i), Model: "0x1af4 0x0001"})
			}
			host.Status.HardwareDetails.NIC = nics
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host),
				fmt.Sprintf("rule=%v actual=%v", tc.Rule, tc.Actual))
		})
	}

}
