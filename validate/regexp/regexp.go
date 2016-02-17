package regexp

import (
	"fmt"
	"regexp"

	"bitbucket.org/syb-devs/goth/validate/internal"
)

// Rule is a helper for other rules that are based on regular expressions
type Rule struct {
	Regexp *regexp.Regexp
}

// Validate checks that the field value matches the regexp passed in the val parameter
func (r *Rule) Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error) {
	fieldVal := internal.GetInterfaceValue(data, field)
	if fieldVal == "" {
		return
	}
	if !r.Regexp.MatchString(fieldVal.(string)) {
		errorInput = fmt.Errorf("The value of field %s does not match regexp %s", field, r.Regexp.String())
		return
	}
	return
}
