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

func getClient(t *testing.T) *clientv3.Client {
	t.Helper()

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed creating etcd client: %s", err.Error())
	}

	return client
}

func TestDeserialize(t *testing.T) {
	testObj := &simple.Simple{
		Field: "value",
	}

	serialized, err := proto.Marshal(testObj)
	if err != nil {
		t.Fatalf("failed to serialize test object: %s", err.Error())
	}

	client := getClient(t)

	_, err = client.Put(context.TODO(), "key", string(serialized))
	if err != nil {
		t.Fatalf("failed to create a key: %s", err.Error())
	}

	store := etcd.NewStorage(client, codec.UniversalCodec())

	into := &simple.Simple{}

	if err := store.Get("key", into); err != nil {
		t.Fatalf("failed getting from store: %s", err.Error())
	}

	if into.Field != "value" {
		t.Fatalf("unfaithfully deserialized value, got: %s", into)
	}

	log.Printf("got it: %s", into)
}

func TestList(t *testing.T) {
	foo := &simple.Simple{
		Field: "foo",
	}
	bar := &simple.Simple{
		Field: "bar",
	}

	serializedFoo, err := proto.Marshal(foo)
	if err != nil {
		t.Fatalf("failed to serialize foo object: %s", err.Error())
	}
	serializedBar, err := proto.Marshal(bar)
	if err != nil {
		t.Fatalf("failed to serialize bar object: %s", err.Error())
	}

	client := getClient(t)
	_, err = client.Put(context.TODO(), "simple/foo", string(serializedFoo))
	if err != nil {
		t.Fatalf("failed to create a key: %s", err.Error())
	}
	_, err = client.Put(context.TODO(), "simple/bar", string(serializedBar))
	if err != nil {
		t.Fatalf("failed to create a key: %s", err.Error())
	}

	store := etcd.NewStorage(client, codec.UniversalCodec())
	into := []simple.Simple{}

	if err := store.List("simple", &into); err != nil {
		t.Fatalf("failed getting from store: %s", err.Error())
	}

	if len(into) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(into))
	}

	if into[0].Field != "bar" {
		t.Fatalf("unfaithfully deserialized value, got: %+v", into[0])
	}
	if into[1].Field != "foo" {
		t.Fatalf("unfaithfully deserialized value, got: %+v", into[1])
	}
}
