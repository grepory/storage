package etcd

import (
	"path"

	"github.com/grepory/storage/runtime/codec"
	"github.com/grepory/storage/storage"
	"go.etcd.io/etcd/clientv3"
)

const (
	etcdRoot = "/sensu.io"
)

// NewPrefixedStorage returns an etcd store whose objects exist under a path prefix.
func NewPrefixedStorage(prefix string, client *clientv3.Client, codec codec.Codec) storage.Store {
	return &storageDecorator{
		storage: &Storage{
			client: client,
			codec:  codec,
		},
		prefix: prefix,
	}
}

// StorageDecorator decorates an etcd store.
type storageDecorator struct {
	storage *Storage
	prefix  string
}

func (sd *storageDecorator) withPrefix(key string) string {
	return path.Join(sd.prefix, key)
}

func (sd *storageDecorator) Create(key string, obj interface{}) error {
	return sd.storage.Create(sd.withPrefix(key), obj)
}

func (sd *storageDecorator) Update(key string, obj interface{}) error {
	return sd.storage.Update(sd.withPrefix(key), obj)
}

func (sd *storageDecorator) Get(key string, obj interface{}) error {
	return sd.storage.Get(sd.withPrefix(key), obj)
}

func (sd *storageDecorator) List(key string, objs interface{}) error {
	return sd.storage.List(sd.withPrefix(key), objs)
}
