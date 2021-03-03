package helpers

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// WrapError adds message context to error.
func WrapError(err error, message string) error {
	return fmt.Errorf("%v: %w", message, err)
}

// GetFunctionName returns name of function
func GetFunctionName(i interface{}) (string, error) {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(fullName, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("wrong function format")
	}
	funcName := parts[len(parts)-1]
	parts = strings.Split(funcName, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("wrong function format")
	}
	return parts[1], nil
}
