//go:build !ignore_autogenerated

/*
Copyright 2024 - Jimmy Lipham & Zeroscaler Authors

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

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScaleTargetRef) DeepCopyInto(out *ScaleTargetRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScaleTargetRef.
func (in *ScaleTargetRef) DeepCopy() *ScaleTargetRef {
	if in == nil {
		return nil
	}
	out := new(ScaleTargetRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZeroScaler) DeepCopyInto(out *ZeroScaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZeroScaler.
func (in *ZeroScaler) DeepCopy() *ZeroScaler {
	if in == nil {
		return nil
	}
	out := new(ZeroScaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ZeroScaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZeroScalerList) DeepCopyInto(out *ZeroScalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ZeroScaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZeroScalerList.
func (in *ZeroScalerList) DeepCopy() *ZeroScalerList {
	if in == nil {
		return nil
	}
	out := new(ZeroScalerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ZeroScalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZeroScalerSpec) DeepCopyInto(out *ZeroScalerSpec) {
	*out = *in
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Target = in.Target
	out.ScaleDownAfter = in.ScaleDownAfter
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZeroScalerSpec.
func (in *ZeroScalerSpec) DeepCopy() *ZeroScalerSpec {
	if in == nil {
		return nil
	}
	out := new(ZeroScalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZeroScalerStatus) DeepCopyInto(out *ZeroScalerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZeroScalerStatus.
func (in *ZeroScalerStatus) DeepCopy() *ZeroScalerStatus {
	if in == nil {
		return nil
	}
	out := new(ZeroScalerStatus)
	in.DeepCopyInto(out)
	return out
}
