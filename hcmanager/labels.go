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

//DeleteHWCCLabel deletes the hwcc label from the baremetal host
func (mgr HardwareClassificationManager) DeleteHWCCLabel(ctx context.Context, hcMetaData v1.ObjectMeta, hosts bmh.BareMetalHostList) []string {
	var deleteLabelError []string
	for _, host := range hosts.Items {
		if err := mgr.DeleteLabels(ctx, hcMetaData, host); err != nil {
			deleteLabelError = append(deleteLabelError, host.Name+" "+err.Error())
		}
	}
	return deleteLabelError
}

// DeleteLabels delete existing label of the baremetal host
func (mgr HardwareClassificationManager) DeleteLabels(ctx context.Context, hcMetaData v1.ObjectMeta, host bmh.BareMetalHost) error {
	labelKey := LabelName + hcMetaData.Name
	// Delete existing labels for the same profile.
	existingLabels := host.GetLabels()
	if existingLabels != nil {
		for key := range existingLabels {
			if key == labelKey {
				delete(existingLabels, key)
			}
		}
		host.SetLabels(existingLabels)
		err := mgr.client.Update(ctx, &host)
		if err != nil {
			return errors.New("Label Delete Failed" + host.Name)
		}
	}
	return nil
}

//SetLabel update labels of baremetal host
func (mgr HardwareClassificationManager) SetLabel(ctx context.Context, hcMetaData v1.ObjectMeta, validHosts []string, BMHList bmh.BareMetalHostList) []string {
	labelKey := LabelName + hcMetaData.Name
	var setLabelError []string

	for _, host := range BMHList.Items {
		if getHostName(host.Name, validHosts) {
			labels := host.GetLabels()
			if labels == nil {
				labels = make(map[string]string)
			}
			if hcMetaData.Labels != nil {
				for _, value := range hcMetaData.Labels {
					if value == "" {
						labels[labelKey] = DefaultLabel
					} else {
						labels[labelKey] = value
					}
				}
			} else {
				labels[labelKey] = DefaultLabel
			}
			mgr.Log.Info("Set Label", host.Name, labels)
			// set updated labels to host
			host.SetLabels(labels)
			if err := mgr.client.Update(ctx, &host); err != nil {
				setLabelError = append(setLabelError, host.Name+" "+err.Error())
			}
		} else {
			if err := mgr.DeleteLabels(ctx, hcMetaData, host); err != nil {
				setLabelError = append(setLabelError, host.Name+" "+err.Error())
			}
		}
	}
	return setLabelError
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
