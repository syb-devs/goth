package validate

import (
	"errors"
	"fmt"
	"strings"
)

func init() {
	RegisterRule("hasPrefix", &hasPrefixRule{})
}

var ErrHasPrefixParamCount = errors.New("This rule needs one mandatory parameter")

// hasPrefixRule struct holds Validate() method to satisfy the Validator interface.
type hasPrefixRule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *hasPrefixRule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	if len(params) == 0 || len(params) > 1 {
		errorLogic = ErrHasPrefixParamCount
	}
	hasPrefixParam := params[0]

	fieldVal := getInterfaceValue(data, field)
	hasPrefixStr := mustStringify(fieldVal)
	if strings.HasPrefix(hasPrefixStr, hasPrefixParam) {
		return
	}
	errorInput = fmt.Errorf("The field %s = %s should start with %s.", field, hasPrefixStr, hasPrefixParam)
	return
}
