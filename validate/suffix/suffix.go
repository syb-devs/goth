package suffix

import (
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/syb-devs/goth/validate"
	"bitbucket.org/syb-devs/goth/validate/internal"
)

func init() {
	validate.RegisterRule("hasSuffix", &rule{})
}

// ErrParamCount is returned when the number of rule parameters does not match the expected
var ErrParamCount = errors.New("This rule needs one mandatory parameter")

type rule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *rule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	if len(params) == 0 || len(params) > 1 {
		errorLogic = ErrParamCount
		return
	}
	hasSuffixParam := params[0]

	fieldVal := internal.GetInterfaceValue(data, field)
	hasSuffixStr := internal.MustStringify(fieldVal)
	if strings.HasSuffix(hasSuffixStr, hasSuffixParam) {
		return
	}
	errorInput = fmt.Errorf("The field %s = %s should end with %s.", field, hasSuffixStr, hasSuffixParam)
	return
}
