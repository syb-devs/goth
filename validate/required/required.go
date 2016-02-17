package required

import (
	"fmt"
	"reflect"

	"bitbucket.org/syb-devs/goth/validate"
	"bitbucket.org/syb-devs/goth/validate/internal"
)

func init() {
	validate.RegisterRule("required", &rule{})
}

type rule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *rule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	fieldVal := reflect.ValueOf(internal.GetInterfaceValue(data, field))
	if isZero(fieldVal) {
		errorInput = fmt.Errorf("a value is required for %s", field)
	}
	return
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}
