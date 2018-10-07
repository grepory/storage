package server

import "github.com/grepory/storage/apis/meta"

// A RestStrategy is an interface to handle REST requests for all meta objects.
type RestStrategy interface {
	PrepareForUpdate() meta.Object
	Validate(obj meta.Object) error
}
