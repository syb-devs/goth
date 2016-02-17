package validate

import (
	"errors"
	"fmt"
	"strings"
)

func init() {
	RegisterRule("hasSuffix", &hasSufixRule{})
}

var ErrHasSufixParamCount = errors.New("This rule needs one mandatory parameter")

// hasSufixRule struct holds Validate() method to satisfy the Validator interface.
type hasSufixRule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *hasSufixRule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	if len(params) == 0 || len(params) > 1 {
		errorLogic = ErrHasSufixParamCount
	}
	hasSufixParam := params[0]

	fieldVal := getInterfaceValue(data, field)
	hasSufixStr := mustStringify(fieldVal)
	if strings.HasPrefix(hasSufixStr, hasSufixParam) {
		return
	}
	errorInput = fmt.Errorf("The field %s = %s should end with %s.", field, hasSufixStr, hasSufixParam)
	return
}
