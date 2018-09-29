package meta

// Object lets you work with object metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field (Name, UID, Namespace on lists) will be a no-op and return
// a default value.
type Object interface {
	GetNamespace() string
	SetNamespace(namespace string)
	GetName() string
	SetName(name string)
	GetUUID() string
	SetUUID(uuid string)
	GetResourceVersion() string
	SetResourceVersion(version string)
	GetCreationTimestamp() Time
	SetCreationTimestamp(timestamp Time)
	GetDeletionTimestamp() *Time
	SetDeletionTimestamp(timestamp *Time)
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
	GetClusterName() string
	SetClusterName(clusterName string)
}
