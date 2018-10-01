package simple_test

import (
	"testing"

	"github.com/grepory/storage/apis/meta"
	"github.com/grepory/storage/apis/simple"
	"github.com/grepory/storage/runtime/codec"
)

var serializedExample = `{
	"kind": "Simple",
	"apiVersion": "storage.apis.simple/v0",
	"metadata": {
		"name": "name",
		"namespace": "namespace"
	},
	"field": "field"
}`

func TestTypeMeta(t *testing.T) {
	var obj interface{}
	obj = &simple.Simple{}
	if _, ok := obj.(meta.GroupVersionKind); !ok {
		t.Error("unable to cast *Simple to GroupVersionKind")
	}
}

func TestDeserialization(t *testing.T) {
	var obj interface{}
	obj = &simple.Simple{}

	if err := codec.UniversalCodec().Decode([]byte(serializedExample), obj); err != nil {
		t.Errorf(err.Error())
	}

	typed := obj.(*simple.Simple)
	if typed.Kind != "Simple" {
		t.Errorf("expected object kind: %s, got %s", "Simple", typed.Kind)
	}

	if typed.Name != "name" {
		t.Errorf("expected metadata name: %s, got %s", "name", typed.Name)
	}

	gvk, ok := obj.(meta.GroupVersionKind)
	if !ok {
		t.Error("unable to cast Simple type to GroupVersionKind")
		t.FailNow()
	}

	if gvk.GetGroup() != "storage.apis.simple" {
		t.Errorf("expected group: %s, got %s", "storage.apis.simple", gvk.GetGroup())
	}

	if gvk.GetVersion() != "v0" {
		t.Errorf("expected version: %s, got %s", "v0", gvk.GetVersion())
	}
}
