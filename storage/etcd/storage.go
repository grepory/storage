package etcd

import (
	"context"
	"errors"

	"github.com/grepory/storage/runtime/codec"
	"go.etcd.io/etcd/clientv3"
)

var (
	InvalidInterfaceError = errors.New("object cannot be serialized")
)

// Storage is a light wrapper around an etcd client.
type Storage struct {
	Client *clientv3.Client
	Codec  codec.Codec
}

// Get a key from storage and deserialize it into objPtr.
func (s *Storage) Get(key string, objPtr interface{}) error {
	resp, err := s.Client.Get(context.TODO(), key)
	if err != nil {
		return err
	}

	v := resp.Kvs[0].Value

	return s.Codec.Decode(v, objPtr)
}

// Create an object in the store.
func (s *Storage) Create(key string, objPtr interface{}) error {
	serialized, err := s.Codec.Encode(objPtr)
	if err != nil {
		return err
	}

	txn := s.Client.Txn(context.TODO()).If(
		keyNotFound(key),
	).Then(
		put(key, string(serialized)),
	)

	resp, err := txn.Commit()
	if err != nil {
		return err
	}

	if !resp.Succeeded {
		return errors.New("could not create existing object")
	}

	return nil
}

// Update a key given with the serialized object.
func (s *Storage) Update(key string, objPtr interface{}) error {
	serialized, err := s.Codec.Encode(objPtr)
	if err != nil {
		return err
	}

	txn := s.Client.Txn(context.TODO()).If(
		keyFound(key),
	).Then(
		put(key, string(serialized)),
	)

	resp, err := txn.Commit()
	if err != nil {
		return err
	}

	if !resp.Succeeded {
		return errors.New("could not update non-existent object")
	}

	return nil
}

func keyFound(key string) clientv3.Cmp {
	return clientv3.Compare(
		clientv3.CreateRevision(key),
		">",
		0,
	)
}

func keyNotFound(key string) clientv3.Cmp {
	return clientv3.Compare(
		clientv3.CreateRevision(key),
		"=",
		0,
	)
}

func put(key, value string) clientv3.Op {
	return clientv3.OpPut(key, value)
}
