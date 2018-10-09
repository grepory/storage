package storage

// A Store provides the methods necessary for interacting with
// objects in some form of storage.
type Store interface {
	Create(key string, obj interface{}) error
	Update(key string, obj interface{}) error
	Get(key string, obj interface{}) error
	List(key string, objs interface{}) error
}
