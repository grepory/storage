package storage

type Store interface {
	Create(key string, obj interface{}) error
	Update(key string, obj interface{}) error
	Get(key string, obj interface{}) error
}
