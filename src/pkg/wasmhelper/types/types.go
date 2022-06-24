package types

import (
	"fmt"
	"strings"

	k8stypes "k8s.io/apimachinery/pkg/types"
)

type Operation string

// Operation constants
const (
	Create  Operation = "CREATE"
	Update  Operation = "UPDATE"
	Delete  Operation = "DELETE"
	Connect Operation = "CONNECT"
)

type HookRequest struct {
	SubResource        string       `json:"subResource,omitempty" protobuf:"bytes,4,opt,name=subResource"`
	RequestSubResource string       `json:"requestSubResource,omitempty" protobuf:"bytes,15,opt,name=requestSubResource"`
	Name               string       `json:"name,omitempty" protobuf:"bytes,5,opt,name=name"`
	Namespace          string       `json:"namespace,omitempty" protobuf:"bytes,6,opt,name=namespace"`
	Operation          Operation    `json:"operation" protobuf:"bytes,7,opt,name=operation"`
	DryRun             *bool        `json:"dryRun,omitempty" protobuf:"varint,11,opt,name=dryRun"`
	UID                k8stypes.UID `json:"uid" protobuf:"bytes,1,opt,name=uid"`

	Kind            GroupVersionKind      `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	Resource        GroupVersionResource  `json:"resource" protobuf:"bytes,3,opt,name=resource"`
	RequestKind     *GroupVersionKind     `json:"requestKind,omitempty" protobuf:"bytes,13,opt,name=requestKind"`
	RequestResource *GroupVersionResource `json:"requestResource,omitempty" protobuf:"bytes,14,opt,name=requestResource"`
	UserInfo        UserInfo              `json:"userInfo" protobuf:"bytes,8,opt,name=userInfo"`
	Object          RawExtension          `json:"object,omitempty" protobuf:"bytes,9,opt,name=object"`
	OldObject       RawExtension          `json:"oldObject,omitempty" protobuf:"bytes,10,opt,name=oldObject"`
	Options         RawExtension          `json:"options,omitempty" protobuf:"bytes,12,opt,name=options"`

	Result bool
}

type GroupVersionKind struct {
	Group   string `json:"group" protobuf:"bytes,1,opt,name=group"`
	Version string `json:"version" protobuf:"bytes,2,opt,name=version"`
	Kind    string `json:"kind" protobuf:"bytes,3,opt,name=kind"`
}

func (gvk GroupVersionKind) String() string {
	return gvk.Group + "/" + gvk.Version + ", Kind=" + gvk.Kind
}

type GroupVersionResource struct {
	Group    string `json:"group" protobuf:"bytes,1,opt,name=group"`
	Version  string `json:"version" protobuf:"bytes,2,opt,name=version"`
	Resource string `json:"resource" protobuf:"bytes,3,opt,name=resource"`
}

func (gvr *GroupVersionResource) String() string {
	if gvr == nil {
		return "<nil>"
	}
	return strings.Join([]string{gvr.Group, "/", gvr.Version, ", Resource=", gvr.Resource}, "")
}

type RawExtension struct {
	Raw    []byte `json:"-" protobuf:"bytes,1,opt,name=raw"`
	Object Object `json:"-"`
}

type Object interface {
	GetObjectKind() ObjectKind
	DeepCopyObject() Object
}

type ObjectKind interface {
	SetGroupVersionKind(kind GroupVersionKind)
	GroupVersionKind() GroupVersionKind
}

type UserInfo struct {
	// The name that uniquely identifies this user among all active users.
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,1,opt,name=username"`
	// A unique value that identifies this user across time. If this user is
	// deleted and another user by the same name is added, they will have
	// different UIDs.
	// +optional
	UID string `json:"uid,omitempty" protobuf:"bytes,2,opt,name=uid"`
	// The names of groups this user is a part of.
	// +optional
	Groups []string `json:"groups,omitempty" protobuf:"bytes,3,rep,name=groups"`
	// Any additional information provided by the authenticator.
	// +optional
	Extra map[string]ExtraValue `json:"extra,omitempty" protobuf:"bytes,4,rep,name=extra"`
}

type ExtraValue []string

func (t ExtraValue) String() string {
	return fmt.Sprintf("%v", []string(t))
}
