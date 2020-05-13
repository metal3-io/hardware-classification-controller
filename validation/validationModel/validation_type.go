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

package validationmodel

//RAM contains ram details fetched from the introspection data
type RAM struct {
	RAMGb int `json:"ramMebibytes"`
}

//HardwareSystemVendor contains hardware manufacturer details fetched from the introspection data
type HardwareSystemVendor struct {
	Manufacturer string `json:"manufacturer"`
}

//NIC contains the nic details fetched from the introspection data
type NIC struct {
	Name  string `json:"name"`
	PXE   bool   `json:"pxe"`
	Count int    `json:"count"`
}

//Storage contains disk details fetched from the introspection data
type Storage struct {
	Count int    `json:"count"`
	Disk  []Disk `json:"disk"`
}

//Disk contains disk size fetched from the introspection data
type Disk struct {
	Name   string `json:"name"`
	SizeGb int64  `json:"sizeBytes"`
}

//CPU contains the clockspeed and count details fetched from the introspection data
type CPU struct {
	Count      int     `json:"count"`
	ClockSpeed float64 `json:"clockspeed"`
}
