package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func TestGetLabelDetails(t *testing.T) {
	testCases := []struct {
		Scenario string
		Profile  hwcc.HardwareClassification
		Label    string
		Value    string
	}{
		{
			Scenario: "default-value",
			Profile: hwcc.HardwareClassification{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "profile-name",
					Namespace: "profile-namespace",
				},
			},
			Label: "hardwareclassification.metal3.io/profile-name",
			Value: "matches",
		},
		{
			Scenario: "explicit-value",
			Profile: hwcc.HardwareClassification{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "profile-name",
					Namespace: "profile-namespace",
					Labels: map[string]string{
						"hardwareclassification.metal3.io/profile-name": "alternate-value",
					},
				},
			},
			Label: "hardwareclassification.metal3.io/profile-name",
			Value: "alternate-value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			label, value := getLabelDetails(&tc.Profile)
			assert.Equal(t, tc.Label, label)
			assert.Equal(t, tc.Value, value)
		})
	}
}

func TestSetLabel(t *testing.T) {
	testCases := []struct {
		Scenario string
		Host     bmh.BareMetalHost
		Labels   map[string]string
		Expected bool
	}{
		{
			Scenario: "no-existing-labels",
			Host: bmh.BareMetalHost{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
				},
			},
			Labels: map[string]string{
				"name": "value",
			},
			Expected: true,
		},
		{
			Scenario: "has-same-label",
			Host: bmh.BareMetalHost{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
					Labels: map[string]string{
						"name": "value",
					},
				},
			},
			Labels: map[string]string{
				"name": "value",
			},
			Expected: false,
		},
		{
			Scenario: "change-value",
			Host: bmh.BareMetalHost{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
					Labels: map[string]string{
						"name": "old-value",
					},
				},
			},
			Labels: map[string]string{
				"name": "value",
			},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			changed := setLabel(&tc.Host, "name", "value")
			assert.Equal(t, tc.Labels, tc.Host.GetLabels())
			assert.Equal(t, tc.Expected, changed)
		})
	}
}

func TestDeleteLabel(t *testing.T) {
	testCases := []struct {
		Scenario string
		Host     bmh.BareMetalHost
		Labels   map[string]string
		Expected bool
	}{
		{
			Scenario: "no-existing-labels",
			Host: bmh.BareMetalHost{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
				},
			},
			Labels:   nil,
			Expected: false,
		},
		{
			Scenario: "has-label",
			Host: bmh.BareMetalHost{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
					Labels: map[string]string{
						"name": "value",
					},
				},
			},
			Labels:   map[string]string{},
			Expected: true,
		},
		{
			Scenario: "has-other-label",
			Host: bmh.BareMetalHost{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
					Labels: map[string]string{
						"different-name": "value",
					},
				},
			},
			Labels: map[string]string{
				"different-name": "value",
			},
			Expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			changed := deleteLabel(&tc.Host, "name")
			assert.Equal(t, tc.Labels, tc.Host.GetLabels())
			assert.Equal(t, tc.Expected, changed)
		})
	}
}
