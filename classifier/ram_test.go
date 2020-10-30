package classifier

import (
	"fmt"
	"testing"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func TestCheckRAM(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Ram
		Actual   int
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   32 * 1024,
			Expected: true,
		},
		{
			Scenario: "no-min-max",
			Rule: &hwcc.Ram{
				MinimumSizeGB: 0,
				MaximumSizeGB: 0,
			},
			Actual:   32 * 1024,
			Expected: true,
		},
		{
			Scenario: "within-max",
			Rule: &hwcc.Ram{
				MinimumSizeGB: 0,
				MaximumSizeGB: 100,
			},
			Actual:   32 * 1024,
			Expected: true,
		},
		{
			Scenario: "within-min",
			Rule: &hwcc.Ram{
				MinimumSizeGB: 1,
				MaximumSizeGB: 0,
			},
			Actual:   32 * 1024,
			Expected: true,
		},
		{
			Scenario: "under-min",
			Rule: &hwcc.Ram{
				MinimumSizeGB: 100,
				MaximumSizeGB: 0,
			},
			Actual:   99 * 1024,
			Expected: false,
		},
		{
			Scenario: "over-max",
			Rule: &hwcc.Ram{
				MinimumSizeGB: 0,
				MaximumSizeGB: 1,
			},
			Actual:   32 * (1024 ^ 2),
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Ram: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						RAMMebibytes: tc.Actual,
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host),
				fmt.Sprintf("rule=%v actual=%v", tc.Rule, tc.Actual))
		})
	}
}
