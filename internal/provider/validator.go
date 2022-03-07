package provider

import (
	"fmt"
	"strconv"
)

// validateModePermission checks that the given input string is a valid file permission string,
// expressed in numeric notation.
// See: https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation
func validateModePermission(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 3 || len(v) > 4 {
		es = append(es, fmt.Errorf("bad mode for file - string length should be 3 or 4 digits: %s", v))
	}

	fileMode, err := strconv.ParseInt(v, 8, 64)
	if err != nil || fileMode > 0777 || fileMode < 0 {
		es = append(es, fmt.Errorf("bad mode for file - must be expressed in octal numeric notation, of 3 or 4 digits: %s", v))
	}

	return
}
