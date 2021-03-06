// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpaCert) DeepCopyInto(out *IpaCert) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpaCert.
func (in *IpaCert) DeepCopy() *IpaCert {
	if in == nil {
		return nil
	}
	out := new(IpaCert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IpaCert) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpaCertData) DeepCopyInto(out *IpaCertData) {
	*out = *in
	in.Issued.DeepCopyInto(&out.Issued)
	in.Expiry.DeepCopyInto(&out.Expiry)
	if in.DnsNames != nil {
		in, out := &in.DnsNames, &out.DnsNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpaCertData.
func (in *IpaCertData) DeepCopy() *IpaCertData {
	if in == nil {
		return nil
	}
	out := new(IpaCertData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpaCertList) DeepCopyInto(out *IpaCertList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IpaCert, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpaCertList.
func (in *IpaCertList) DeepCopy() *IpaCertList {
	if in == nil {
		return nil
	}
	out := new(IpaCertList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IpaCertList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpaCertSpec) DeepCopyInto(out *IpaCertSpec) {
	*out = *in
	if in.AdditionalNames != nil {
		in, out := &in.AdditionalNames, &out.AdditionalNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpaCertSpec.
func (in *IpaCertSpec) DeepCopy() *IpaCertSpec {
	if in == nil {
		return nil
	}
	out := new(IpaCertSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpaCertStatus) DeepCopyInto(out *IpaCertStatus) {
	*out = *in
	in.CertData.DeepCopyInto(&out.CertData)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpaCertStatus.
func (in *IpaCertStatus) DeepCopy() *IpaCertStatus {
	if in == nil {
		return nil
	}
	out := new(IpaCertStatus)
	in.DeepCopyInto(out)
	return out
}
