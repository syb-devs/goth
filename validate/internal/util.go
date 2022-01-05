package internal

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/syb-devs/goth/validate"
)

var ErrUnsupportedType = errors.New("Unsupported type for rule")

// GetInterfaceValue returns the value of a given interface using reflection.
func GetInterfaceValue(data interface{}, name string) interface{} {
	return reflect.ValueOf(data).FieldByName(name).Interface()
}

// toString returns a literal representation of a given value.
// The second parameter indicates whether a conversion was possible or not.
func toString(value interface{}) (string, bool) {
	switch v := value.(type) {
	case string, *string, int, *int, int32, *int32, int64, *int64:
		return fmt.Sprintf("%v", v), true
	default:
		return "", false
	}
}

// MustStringify tries to convert the given value to string type and panics if not possible.
func MustStringify(value interface{}) string {
	strVal, ok := toString(value)
	if ok == false {
		panic(ErrUnsupportedType)
	}
	return strVal
}

func FindErrors(t *testing.T, errList validate.FieldErrors, patterns map[string][]string) {
	for field, patErrs := range patterns {
		valErrs := errList[field]
		for _, patErr := range patErrs {
			found := false
			for _, valErr := range valErrs {
				if patErr == valErr.Error() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error '%s' for field '%s', but was not found", patErr, field)
			}
		}
	}
}
