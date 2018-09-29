package etcd

import (
	"fmt"
	"path"
	"strings"

	"github.com/grepory/storage/apis/meta"
)

const (
	rootPrefix = "/sensu.io"
)

var (
	NameMayNotBe      = []string{".", ".."}
	NameMayNotContain = []string{"/", "%"}
)

func IsValidPathSegmentName(name string) error {
	for _, illegalName := range NameMayNotBe {
		if name == illegalName {
			return fmt.Errorf("name may not be '%s'", illegalName)
		}
	}

	for _, illegalContent := range NameMayNotContain {
		if strings.Contains(name, illegalContent) {
			return fmt.Errorf("name may not contain '%s'", illegalContent)
		}
	}

	return nil
}

func NamespaceKeyFunc(prefix string, obj meta.Object) (string, error) {
	name := obj.GetName()
	if err := IsValidPathSegmentName(name); err != nil {
		return "", err
	}

	return path.Join(prefix, obj.GetNamespace(), name), nil
}

func NoNamespaceKeyFunc(prefix string, obj meta.Object) (string, error) {
	name := obj.GetName()
	if err := IsValidPathSegmentName(name); err != nil {
		return "", err
	}

	return path.Join(prefix, name), nil
}
