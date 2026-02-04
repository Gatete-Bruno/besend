package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func (in *EmailSenderConfig) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func (in *EmailSenderConfigList) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func (in *Email) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func (in *EmailList) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func (in *EmailSenderConfig) DeepCopy() *EmailSenderConfig {
	if in == nil {
		return nil
	}
	out := new(EmailSenderConfig)
	in.DeepCopyInto(out)
	return out
}

func (in *EmailSenderConfig) DeepCopyInto(out *EmailSenderConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *EmailSenderConfigList) DeepCopy() *EmailSenderConfigList {
	if in == nil {
		return nil
	}
	out := new(EmailSenderConfigList)
	in.DeepCopyInto(out)
	return out
}

func (in *EmailSenderConfigList) DeepCopyInto(out *EmailSenderConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EmailSenderConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *Email) DeepCopy() *Email {
	if in == nil {
		return nil
	}
	out := new(Email)
	in.DeepCopyInto(out)
	return out
}

func (in *Email) DeepCopyInto(out *Email) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *EmailList) DeepCopy() *EmailList {
	if in == nil {
		return nil
	}
	out := new(EmailList)
	in.DeepCopyInto(out)
	return out
}

func (in *EmailList) DeepCopyInto(out *EmailList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Email, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}
