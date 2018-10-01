package etcd_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/grepory/storage/apis/simple"
	"github.com/grepory/storage/runtime/codec"
	"github.com/grepory/storage/storage/etcd"
	"go.etcd.io/etcd/clientv3"
)

func TestDeserialize(t *testing.T) {
	testObj := &simple.Simple{
		Field: "value",
	}

	serialized, err := proto.Marshal(testObj)
	if err != nil {
		log.Println("failed to serialize test object")
		t.FailNow()
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Println("failed creating etcd client")
		t.FailNow()
	}

	client.Put(context.TODO(), "key", string(serialized))

	store := etcd.NewStorage(client, codec.UniversalCodec())

	into := &simple.Simple{}

	if err := store.Get("key", into); err != nil {
		log.Println("failed getting from store")
		t.FailNow()
	}

	if into.Field != "value" {
		log.Printf("unfaithfully deserialized value, got: %s", into)
		t.FailNow()
	}

	log.Printf("got it: %s", into)
}
