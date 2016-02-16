package validate

import (
	"errors"
	"fmt"
	"strings"
)

var ErrContainsParamCount = errors.New("This rule needs one mandatory parameter")

// containsRule struct holds Validate() method to satisfy the Validator interface.
type containsRule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *containsRule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	if len(params) == 0 || len(params) > 1 {
		errorLogic = ErrContainsParamCount
	}
	containsParam := params[0]

	fieldVal := getInterfaceValue(data, field)
	containsStr := mustStringify(fieldVal)
	if strings.Contains(containsStr, containsParam) {
		return
	}
	errorInput = fmt.Errorf("The field %s = %s should contain %s.", field, containsStr, containsParam)
	return
}