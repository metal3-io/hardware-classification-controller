// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cpu) DeepCopyInto(out *Cpu) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cpu.
func (in *Cpu) DeepCopy() *Cpu {
	if in == nil {
		return nil
	}
	out := new(Cpu)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Disk) DeepCopyInto(out *Disk) {
	*out = *in
	if in.DiskSelector != nil {
		in, out := &in.DiskSelector, &out.DiskSelector
		*out = make([]DiskSelector, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Disk.
func (in *Disk) DeepCopy() *Disk {
	if in == nil {
		return nil
	}
	out := new(Disk)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiskSelector) DeepCopyInto(out *DiskSelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiskSelector.
func (in *DiskSelector) DeepCopy() *DiskSelector {
	if in == nil {
		return nil
	}
	out := new(DiskSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HardwareCharacteristics) DeepCopyInto(out *HardwareCharacteristics) {
	*out = *in
	if in.Cpu != nil {
		in, out := &in.Cpu, &out.Cpu
		*out = new(Cpu)
		**out = **in
	}
	if in.Disk != nil {
		in, out := &in.Disk, &out.Disk
		*out = new(Disk)
		(*in).DeepCopyInto(*out)
	}
	if in.Nic != nil {
		in, out := &in.Nic, &out.Nic
		*out = new(Nic)
		**out = **in
	}
	if in.Ram != nil {
		in, out := &in.Ram, &out.Ram
		*out = new(Ram)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HardwareCharacteristics.
func (in *HardwareCharacteristics) DeepCopy() *HardwareCharacteristics {
	if in == nil {
		return nil
	}
	out := new(HardwareCharacteristics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HardwareClassification) DeepCopyInto(out *HardwareClassification) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HardwareClassification.
func (in *HardwareClassification) DeepCopy() *HardwareClassification {
	if in == nil {
		return nil
	}
	out := new(HardwareClassification)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HardwareClassification) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HardwareClassificationList) DeepCopyInto(out *HardwareClassificationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HardwareClassification, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HardwareClassificationList.
func (in *HardwareClassificationList) DeepCopy() *HardwareClassificationList {
	if in == nil {
		return nil
	}
	out := new(HardwareClassificationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HardwareClassificationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HardwareClassificationSpec) DeepCopyInto(out *HardwareClassificationSpec) {
	*out = *in
	in.HardwareCharacteristics.DeepCopyInto(&out.HardwareCharacteristics)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HardwareClassificationSpec.
func (in *HardwareClassificationSpec) DeepCopy() *HardwareClassificationSpec {
	if in == nil {
		return nil
	}
	out := new(HardwareClassificationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HardwareClassificationStatus) DeepCopyInto(out *HardwareClassificationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HardwareClassificationStatus.
func (in *HardwareClassificationStatus) DeepCopy() *HardwareClassificationStatus {
	if in == nil {
		return nil
	}
	out := new(HardwareClassificationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nic) DeepCopyInto(out *Nic) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nic.
func (in *Nic) DeepCopy() *Nic {
	if in == nil {
		return nil
	}
	out := new(Nic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ram) DeepCopyInto(out *Ram) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ram.
func (in *Ram) DeepCopy() *Ram {
	if in == nil {
		return nil
	}
	out := new(Ram)
	in.DeepCopyInto(out)
	return out
}
