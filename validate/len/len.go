package len

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/syb-devs/goth/validate"
	"github.com/syb-devs/goth/validate/internal"
)

func init() {
	validate.RegisterRule("len", &rule{})
}

// ErrParamCount is returned when the number of rule parameters does not match the expected
var ErrParamCount = errors.New("this rule needs two mandatory params, operator and value")

// rule struct holds Validate() method to satisfy the Validator interface.
type rule struct{}

// Validate checks that the given data conforms to the length constraints given as parameters.
func (r *rule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	var op, lengthParam string

	switch len(params) {
	case 1:
		op = "="
		lengthParam = params[0]
	case 2:
		op = params[0]
		lengthParam = params[1]
	default:
		errorLogic = ErrParamCount
		return
	}

	requiredLength, errorLogic := strconv.Atoi(lengthParam)
	if errorLogic != nil {
		return
	}

	fieldVal := internal.GetInterfaceValue(data, field)
	length := utf8.RuneCountInString(internal.MustStringify(fieldVal))

	var ok bool
	var opLiteral string
	switch op {
	case "=":
		ok = length == requiredLength
		opLiteral = "equal to"
	case ">":
		ok = length > requiredLength
		opLiteral = "greater than"
	case ">=":
		ok = length >= requiredLength
		opLiteral = "greater than, or equal to"
	case "<":
		ok = length < requiredLength
		opLiteral = "lower than"
	case "<=":
		ok = length < requiredLength
		opLiteral = "lower than, or equal to"
	default:
		return nil, errors.New("invalid operator")
	}

	if !ok {
		errorInput = fmt.Errorf("the field %s should have a length %s %d. Actual length: %d", field, opLiteral, requiredLength, length)
		return
	}
	return
}
