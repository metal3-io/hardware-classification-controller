package classifier

import (
	"testing"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func TestCheckString(t *testing.T) {
	assert.True(t, checkString("Dell Inc.", "Dell Inc."))
	assert.True(t, checkString("1.5.6", "1.5.6"))
	assert.True(t, checkString("", "Dell Inc."))
	assert.True(t, checkString("", "1.5.6"))
	assert.False(t, checkString("1.5", "1.5.6"))
	assert.False(t, checkString("Dell", "Dell Inc."))

}

func TestFirmwareVendor(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Firmware
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
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor: "Dell Inc.",
				},
			},
			Actual:   "Dell Inc.",
			Expected: true,
		},
		{
			Scenario: "unmatched",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor: "Dell Inc.",
				},
			},
			Actual:   "Dell",
			Expected: false,
		},
		{
			Scenario: "empty",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor:  "",
					Version: "1.5.6",
				},
			},
			Actual:   "Dell Inc.",
			Expected: true,
		},
		{
			Scenario: "empty-host-details",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor: "Dell Inc.",
				},
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
						Firmware: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						Firmware: bmh.Firmware{
							BIOS: bmh.BIOS{
								Vendor:  tc.Actual,
								Version: "1.5.6",
							},
						},
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host))
		})
	}
}

func TestCheckFirmwareVersion(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Firmware
		Actual   string
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   "1.5.6",
			Expected: true,
		},
		{
			Scenario: "matched",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Version: "1.5.6",
				},
			},
			Actual:   "1.5.6",
			Expected: true,
		},
		{
			Scenario: "unmatched",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Version: "1.5.6",
				},
			},
			Actual:   "1.5",
			Expected: false,
		},
		{
			Scenario: "empty",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Version: "",
				},
			},
			Actual:   "1.5.6",
			Expected: true,
		},
		{
			Scenario: "empty-host-details",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Version: "1.5.6",
				},
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
						Firmware: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						Firmware: bmh.Firmware{
							BIOS: bmh.BIOS{
								Version: tc.Actual,
							},
						},
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host))
		})
	}
}

func TestFirmware(t *testing.T) {
	testCases := []struct {
		Scenario      string
		Rule          *hwcc.Firmware
		ActualVendor  string
		ActualVersion string
		Expected      bool
	}{
		{
			Scenario:      "nil",
			Rule:          nil,
			ActualVendor:  "Dell Inc.",
			ActualVersion: "1.5.6",
			Expected:      true,
		},
		{
			Scenario: "both_matched",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor:  "Dell Inc.",
					Version: "1.5.6",
				},
			},
			ActualVendor:  "Dell Inc.",
			ActualVersion: "1.5.6",
			Expected:      true,
		},
		{
			Scenario: "both_unmatched",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor:  "Dell",
					Version: "1.5",
				},
			},
			ActualVendor:  "Dell Inc.",
			ActualVersion: "1.5.6",
			Expected:      false,
		},

		{
			Scenario: "vendor_unmatched",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor:  "Dell",
					Version: "1.5.6",
				},
			},
			ActualVendor:  "Dell Inc.",
			ActualVersion: "1.5.6",
			Expected:      false,
		},

		{
			Scenario: "version_unmatched",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor:  "Dell Inc.",
					Version: "1.5",
				},
			},
			ActualVendor:  "Dell Inc.",
			ActualVersion: "1.5.6",
			Expected:      false,
		},
		{
			Scenario: "empty_version",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor:  "Dell Inc.",
					Version: "",
				},
			},
			ActualVendor:  "Dell Inc.",
			ActualVersion: "1.5.6",
			Expected:      true,
		},

		{
			Scenario: "empty_vendor",
			Rule: &hwcc.Firmware{
				BIOS: hwcc.BIOS{
					Vendor:  "",
					Version: "1.5.6",
				},
			},
			ActualVendor:  "Dell Inc.",
			ActualVersion: "1.5.6",
			Expected:      true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Firmware: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{
						Firmware: bmh.Firmware{
							BIOS: bmh.BIOS{
								Vendor:  tc.ActualVendor,
								Version: tc.ActualVersion,
							},
						},
					},
				},
			}
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host))
		})
	}
}
