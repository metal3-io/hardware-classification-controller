package classifier

import (
	"fmt"
	"testing"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func TestCheckRangeCapacity(t *testing.T) {
	assert.True(t, checkRangeCapacity(bmh.Capacity(0), bmh.Capacity(0), bmh.Capacity(99)))
	assert.True(t, checkRangeCapacity(bmh.Capacity(0), bmh.Capacity(100), bmh.Capacity(99)))
	assert.True(t, checkRangeCapacity(bmh.Capacity(1), bmh.Capacity(0), bmh.Capacity(99)))
	assert.False(t, checkRangeCapacity(bmh.Capacity(100), bmh.Capacity(0), bmh.Capacity(99)))
	assert.False(t, checkRangeCapacity(bmh.Capacity(0), bmh.Capacity(9), bmh.Capacity(99)))
}

func TestCheckDiskCount(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Disk
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
			Rule: &hwcc.Disk{
				MinimumCount: 0,
				MaximumCount: 0,
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "within-max",
			Rule: &hwcc.Disk{
				MinimumCount: 0,
				MaximumCount: 4,
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "within-min",
			Rule: &hwcc.Disk{
				MinimumCount: 1,
				MaximumCount: 0,
			},
			Actual:   2,
			Expected: true,
		},
		{
			Scenario: "under-min",
			Rule: &hwcc.Disk{
				MinimumCount: 4,
				MaximumCount: 0,
			},
			Actual:   2,
			Expected: false,
		},
		{
			Scenario: "over-max",
			Rule: &hwcc.Disk{
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
						Disk: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{},
				},
			}
			disks := []bmh.Storage{}
			for i := 0; i < tc.Actual; i++ {
				disks = append(disks, bmh.Storage{Name: fmt.Sprintf("dev%d", i)})
			}
			host.Status.HardwareDetails.Storage = disks
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host),
				fmt.Sprintf("rule=%v actual=%v", tc.Rule, tc.Actual))
		})
	}
}

func TestCheckDiskSize(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Disk
		Actual   bmh.Capacity
		Expected bool
	}{
		{
			Scenario: "nil",
			Rule:     nil,
			Actual:   20 * bmh.GibiByte,
			Expected: true,
		},
		{
			Scenario: "no-min-max",
			Rule: &hwcc.Disk{
				MinimumIndividualSizeGB: 0,
				MaximumIndividualSizeGB: 0,
			},
			Actual:   20 * bmh.GibiByte,
			Expected: true,
		},
		{
			Scenario: "within-max",
			Rule: &hwcc.Disk{
				MinimumIndividualSizeGB: 0,
				MaximumIndividualSizeGB: 40,
			},
			Actual:   20 * bmh.GibiByte,
			Expected: true,
		},
		{
			Scenario: "within-min",
			Rule: &hwcc.Disk{
				MinimumIndividualSizeGB: 10,
				MaximumIndividualSizeGB: 0,
			},
			Actual:   20 * bmh.GibiByte,
			Expected: true,
		},
		{
			Scenario: "under-min",
			Rule: &hwcc.Disk{
				MinimumIndividualSizeGB: 40,
				MaximumIndividualSizeGB: 0,
			},
			Actual:   20 * bmh.GibiByte,
			Expected: false,
		},
		{
			Scenario: "over-max",
			Rule: &hwcc.Disk{
				MinimumIndividualSizeGB: 0,
				MaximumIndividualSizeGB: 10,
			},
			Actual:   20 * bmh.GibiByte,
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Disk: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{},
				},
			}
			disks := []bmh.Storage{
				bmh.Storage{
					Name:      "/dev/sda",
					SizeBytes: tc.Actual,
				},
			}
			host.Status.HardwareDetails.Storage = disks
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host),
				fmt.Sprintf("rule=%v actual=%v", tc.Rule, tc.Actual))
		})
	}
}

func TestCheckDiskPattern(t *testing.T) {
	testCases := []struct {
		Scenario string
		Rule     *hwcc.Disk
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
			Scenario: "mismatch pattern",
			Rule: &hwcc.Disk{
				MinimumCount: 0,
				MaximumCount: 2,
				DiskSelector: []hwcc.DiskSelector{
					{
						HCTL:       "N:0:0:0",
						Rotational: true,
					},
				},
			},
			Actual:   2,
			Expected: false,
		}, {
			Scenario: "match pattern",
			Rule: &hwcc.Disk{
				MinimumCount: 0,
				MaximumCount: 4,
				DiskSelector: []hwcc.DiskSelector{{
					HCTL:       "0:0:N:0",
					Rotational: true,
				},
				},
			},
			Actual:   1,
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			profile := hwcc.HardwareClassification{
				Spec: hwcc.HardwareClassificationSpec{
					HardwareCharacteristics: hwcc.HardwareCharacteristics{
						Disk: tc.Rule,
					},
				},
			}
			host := bmh.BareMetalHost{
				Status: bmh.BareMetalHostStatus{
					HardwareDetails: &bmh.HardwareDetails{},
				},
			}
			disks := []bmh.Storage{}
			for i := 0; i < tc.Actual; i++ {
				if tc.Scenario == "match pattern" {
					disks = append(disks,
						bmh.Storage{Name: fmt.Sprintf("dev%d", i),
							Rotational:   true,
							SizeBytes:    bmh.Capacity(53687091200),
							Vendor:       "QEMU",
							Model:        "QEMU HARDDISK",
							SerialNumber: "drive-scsi0-0-0-0",
							HCTL:         "0:0:1:0"},
						bmh.Storage{Name: fmt.Sprintf("dev%d", i),
							Rotational:   true,
							SizeBytes:    bmh.Capacity(53687091200),
							Vendor:       "QEMU",
							Model:        "QEMU HARDDISK",
							SerialNumber: "drive-scsi0-0-0-0",
							HCTL:         "0:0:0:0"})
				} else {
					disks = append(disks,
						bmh.Storage{Name: fmt.Sprintf("dev%d", i),
							Rotational:   true,
							SizeBytes:    bmh.Capacity(53687091200),
							Vendor:       "QEMU",
							Model:        "QEMU HARDDISK",
							SerialNumber: "drive-scsi0-0-0-0",
							HCTL:         "0:0:1:0"})
				}
			}
			host.Status.HardwareDetails.Storage = disks
			assert.Equal(t, tc.Expected, ProfileMatchesHost(&profile, &host),
				fmt.Sprintf("rule=%v actual=%v", tc.Rule, tc.Actual))
		})
	}

}
