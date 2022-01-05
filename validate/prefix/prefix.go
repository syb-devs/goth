package prefix

import (
	"errors"
	"fmt"
	"strings"

	"github.com/syb-devs/goth/validate"
	"github.com/syb-devs/goth/validate/internal"
)

func init() {
	validate.RegisterRule("hasPrefix", &hasPrefixRule{})
}

// ErrParamCount is returned when the number of rule parameters does not match the expected
var ErrParamCount = errors.New("This rule needs one mandatory parameter")

// hasPrefixRule struct holds Validate() method to satisfy the Validator interface.
type hasPrefixRule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *hasPrefixRule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	if len(params) != 1 {
		errorLogic = ErrParamCount
		return
	}
	hasPrefixParam := params[0]

	fieldVal := internal.GetInterfaceValue(data, field)
	hasPrefixStr := internal.MustStringify(fieldVal)
	if hasPrefixStr == "" {
		return
	}
	if strings.HasPrefix(hasPrefixStr, hasPrefixParam) {
		return
	}
	errorInput = fmt.Errorf("The field %s = %s should start with %s.", field, hasPrefixStr, hasPrefixParam)
	return
}
