package namer

import (
	"fmt"
	"strings"
)

var (
	illegalCharacters = []string{".", "/"}
)

func containsIllegalString(name string) error {
	for _, c := range illegalCharacters {
		if strings.Contains(name, c) {
			return fmt.Errorf("contains illegal character: %s", c)
		}
	}

	return nil
}

// ValidateName validates the name of an object or namespace.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("cannot be empty string")
	}

	return containsIllegalString(name)
}
