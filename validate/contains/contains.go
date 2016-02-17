package contains

import (
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/syb-devs/goth/validate"
	"bitbucket.org/syb-devs/goth/validate/internal"
)

func init() {
	validate.RegisterRule("contains", &rule{})
}

// ErrParamCount is returned when the number of rule parameters does not match the expected
var ErrParamCount = errors.New("This rule needs one mandatory parameter")

type rule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *rule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	if len(params) != 1 {
		errorLogic = ErrParamCount
		return
	}
	containsParam := params[0]

	fieldVal := internal.GetInterfaceValue(data, field)
	containsStr := internal.MustStringify(fieldVal)
	if strings.Contains(containsStr, containsParam) {
		return
	}
	errorInput = fmt.Errorf("The field %s = %s should contain %s.", field, containsStr, containsParam)
	return
}
