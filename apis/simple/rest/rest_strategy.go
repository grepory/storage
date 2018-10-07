package rest

import (
	"github.com/grepory/storage/apis/meta"
	"github.com/grepory/storage/apis/simple"
)

type SimpleRestStrategy struct{}

func (strategy SimpleRestStrategy) PrepareForUpdate() meta.Object {
	return &simple.Simple{
		TypeMeta: meta.TypeMeta{
			Kind:       "Simple",
			APIVersion: "simple/v1alpha1",
		},
	}
}

func (strategy SimpleRestStrategy) Validate(obj meta.Object) error {
	return nil
}
