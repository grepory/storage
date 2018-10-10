package storage

import "path"

type prefixer struct {
	store  Store
	prefix string
}

// WithPrefix writes keys to a Store with a prefix.
func WithPrefix(prefix string, store Store) Store {
	return prefixer{
		store:  store,
		prefix: prefix,
	}
}

func (p prefixer) withPrefix(key string) string {
	return path.Join(p.prefix, key)
}

func (p prefixer) Create(key string, obj interface{}) error {
	return p.store.Create(p.withPrefix(key), obj)
}

func (p prefixer) Update(key string, obj interface{}) error {
	return p.store.Update(p.withPrefix(key), obj)
}

func (p prefixer) Get(key string, obj interface{}) error {
	return p.store.Get(p.withPrefix(key), obj)
}

func (p prefixer) List(key string, obj interface{}) error {
	return p.store.List(p.withPrefix(key), obj)
}
