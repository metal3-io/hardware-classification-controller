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
	"errors"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteLabel delete existing label of the baremetal host
func (mgr HardwareClassificationManager) deleteLabel(ctx context.Context, hcMetaData v1.ObjectMeta, host bmh.BareMetalHost) error {
	labelKey := LabelName + hcMetaData.Name
	// Delete existing labels for the same profile.
	if host.Status.Provisioning.State == "ready" {
		existingLabels := host.GetLabels()
		if existingLabels != nil {
			for key, value := range existingLabels {
				if key == labelKey {
					mgr.Log.Info("Delete Label", "BareMetalHost", host.Name, key, value)
					delete(existingLabels, key)
				}
			}

			// Updating labels
			host.SetLabels(existingLabels)
			if err := mgr.client.Update(ctx, &host); err != nil {
				return errors.New(host.Name + " " + err.Error())
			}
		}
	}
	return nil
}

// setLabel set label of baremetal host
func (mgr HardwareClassificationManager) setLabel(ctx context.Context, hcMetaData v1.ObjectMeta, host bmh.BareMetalHost) error {
	labelKey := LabelName + hcMetaData.Name
	labels := host.GetLabels()

	if labels == nil {
		labels = make(map[string]string)
	}

	// Update user provided labels else set default label
	if hcMetaData.Labels != nil {
		if val, ok := hcMetaData.Labels[hcMetaData.Name]; ok {
			labels[labelKey] = val
		} else {
			labels[labelKey] = DefaultLabel
		}
	} else {
		labels[labelKey] = DefaultLabel
	}

	mgr.Log.Info("Set Label", "BareMetalHost", host.Name, labelKey, labels[labelKey])
	// set updated labels to host
	host.SetLabels(labels)
	if err := mgr.client.Update(ctx, &host); err != nil {
		return errors.New(host.Name + " " + err.Error())
	}
	return nil
}

//UpdateLabels update labels of baremetal host
func (mgr HardwareClassificationManager) UpdateLabels(ctx context.Context, hcMetaData v1.ObjectMeta, validHosts []string, BMHList bmh.BareMetalHostList) []string {

	var updateLabelError []string

	for _, host := range BMHList.Items {
		if getHostName(host.Name, validHosts) {
			if err := mgr.setLabel(ctx, hcMetaData, host); err != nil {
				updateLabelError = append(updateLabelError, err.Error())
			}
		} else {
			if err := mgr.deleteLabel(ctx, hcMetaData, host); err != nil {
				updateLabelError = append(updateLabelError, err.Error())
			}
		}
	}
	return updateLabelError
}

// getHostName checks if host name is present in validHosts
func getHostName(hostName string, validHosts []string) bool {
	for _, host := range validHosts {
		if hostName == host {
			return true
		}
	}
	return false
}
