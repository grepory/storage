package storage

import "github.com/grepory/storage/apis/meta"

type Store interface {
	Create(key string, obj interface{}) error
	Get(key string, obj interface{}) error
}

type KeyFunc func(m meta.TypeMeta, obj meta.Object) string
