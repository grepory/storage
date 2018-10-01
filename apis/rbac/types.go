package rbac

import "github.com/grepory/storage/apis/meta"

// A Rule holds information which describes an action that can be taken.
type Rule struct {
	// Verbs is a list of verbs that apply to all of the listed
	// resources for this rule. These include "get", "list", "watch",
	// "create", "update", "delete".
	// TODO: add support for "patch" (this is expensive and should be
	// delayed until a further release).
	// TODO: add support for "watch" (via websockets)
	Verbs []string `json:"verbs" protobuf:"bytes,1,rep,name=verbs"`

	// Resources is a list of resources that this rule applies to.
	// "*" represents all resources.
	// TODO: enumerate "resources"
	Resources []string `json:"resources" protobuf:"bytes,2,rep,name=resources"`

	// Names is an optional list of resource names that the rule
	// applies to.
	Names []string `json:"names" protobuf:"bytes,3,rep,name=names"`
}

// A Role applies only to a single Namespace.
type Role struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// Rules hold all of the Rules for this Role.
	Rules []Rule `json:"rules" protobuf:"bytes,2,rep,name=rules"`
}

// ClusterRole is a role that applies to all Namespaces within
// a cluster.
type ClusterRole struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// Rules hold all of the Rules for this Role.
	Rules []Rule `json:"rules" protobuf:"bytes,2,rep,name=rules"`
}

// RoleRef is used to map groups to Roles or ClusterRoles.
type RoleRef struct {
	// Type is the type of role being referenced.
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`

	// Name is the name of the resource being referenced.
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

// A Group is a grouping of RoleRefs.
type Group struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// RoleRefs is the list of RoleRefs that make up the roles for this
	// Group.
	RoleRefs []RoleRef `json:"roleRefs" protobuf:"bytes,2,rep,name=roleRefs"`
}
