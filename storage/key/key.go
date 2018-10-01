// Package paths provides utility functions to produce the path to
// a resource in storage.

package key

import (
	"fmt"
	"path"

	"github.com/grepory/storage/apis/meta"
)

const (
	globalNamespace = ".global"
)

// A KeyerFunc returns the key of an object in storage.
type KeyerFunc func(prefix string, obj meta.Object) (string, error)

// NamespaceKeyerFunc returns a namespaced key for a meta.Object.
func NamespaceKeyerFunc(prefix string, obj meta.Object) (string, error) {
	gvk := meta.GetGroupVersionKind(obj)
	// We should never be passed an object that isn't an API object, but
	// guard against it here.
	if gvk == nil {
		return "", fmt.Errorf("cannot access type metadata for object: %v", obj)
	}

	return path.Join(prefix, gvk.GetKind(), obj.GetNamespace(), obj.GetName()), nil
}

// NoNamespaceKeyerFunc returns a non-namespaced key for a meta.Object
// that exists in the global namespace.
func NoNamespaceKeyerFunc(prefix string, obj meta.Object) (string, error) {
	gvk := meta.GetGroupVersionKind(obj)
	if gvk == nil {
		return "", fmt.Errorf("unable to access type metadata for object: %v", obj)
	}

	return path.Join(prefix, gvk.GetKind(), globalNamespace, obj.GetName()), nil
}
