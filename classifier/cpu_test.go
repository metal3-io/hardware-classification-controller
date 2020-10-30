package classifier

import (
	"testing"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func TestCheckRangeClockSpeed(t *testing.T) {
	assert.True(t, checkRangeClockSpeed(bmh.ClockSpeed(0), bmh.ClockSpeed(0), bmh.ClockSpeed(99)))
	assert.True(t, checkRangeClockSpeed(bmh.ClockSpeed(0), bmh.ClockSpeed(100), bmh.ClockSpeed(99)))
	assert.True(t, checkRangeClockSpeed(bmh.ClockSpeed(1), bmh.ClockSpeed(0), bmh.ClockSpeed(99)))
	assert.False(t, checkRangeClockSpeed(bmh.ClockSpeed(100), bmh.ClockSpeed(0), bmh.ClockSpeed(99)))
	assert.False(t, checkRangeClockSpeed(bmh.ClockSpeed(0), bmh.ClockSpeed(9), bmh.ClockSpeed(99)))
}

func TestCheckCPUSpeed(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Cpu
		Actual   bmh.ClockSpeed
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   bmh.ClockSpeed(99),
			Expected: true,
		},
		{
			Scenario: "no-min-max",
			Rule: &hwcc.Cpu{
				MinimumSpeedMHz: 0,
				MaximumSpeedMHz: 0,
			},
			Actual:   bmh.ClockSpeed(99),
			Expected: true,
		},
		{
			Scenario: "within-max",
			Rule: &hwcc.Cpu{
				MinimumSpeedMHz: 0,
				MaximumSpeedMHz: 100,
			},
			Actual:   bmh.ClockSpeed(99),
			Expected: true,
		},
		{
			Scenario: "within-min",
			Rule: &hwcc.Cpu{
				MinimumSpeedMHz: 1,
				MaximumSpeedMHz: 0,
			},
			Actual:   bmh.ClockSpeed(99),
			Expected: true,
		},
		{
			Scenario: "under-min",
			Rule: &hwcc.Cpu{
				MinimumSpeedMHz: 100,
				MaximumSpeedMHz: 0,
			},
			Actual:   bmh.ClockSpeed(99),
			Expected: false,
		},
		{
			Scenario: "over-max",
			Rule: &hwcc.Cpu{
				MinimumSpeedMHz: 0,
				MaximumSpeedMHz: 9,
			},
			Actual:   bmh.ClockSpeed(99),
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Cpu: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						CPU: bmh.CPU{
							ClockMegahertz: tc.Actual,
						},
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host))
		})
	}
}

func TestCheckCPUCount(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Cpu
		Actual   int
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   8,
			Expected: true,
		},
		{
			Scenario: "no-min-max",
			Rule: &hwcc.Cpu{
				MinimumCount: 0,
				MaximumCount: 0,
			},
			Actual:   8,
			Expected: true,
		},
		{
			Scenario: "within-max",
			Rule: &hwcc.Cpu{
				MinimumCount: 0,
				MaximumCount: 100,
			},
			Actual:   8,
			Expected: true,
		},
		{
			Scenario: "within-min",
			Rule: &hwcc.Cpu{
				MinimumCount: 1,
				MaximumCount: 0,
			},
			Actual:   8,
			Expected: true,
		},
		{
			Scenario: "under-min",
			Rule: &hwcc.Cpu{
				MinimumCount: 100,
				MaximumCount: 0,
			},
			Actual:   8,
			Expected: false,
		},
		{
			Scenario: "over-max",
			Rule: &hwcc.Cpu{
				MinimumCount: 0,
				MaximumCount: 4,
			},
			Actual:   8,
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Cpu: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						CPU: bmh.CPU{
							Count: tc.Actual,
						},
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host))
		})
	}
}
