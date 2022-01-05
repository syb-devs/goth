package prefix_test

import (
	"testing"

	"github.com/syb-devs/goth/validate"
	"github.com/syb-devs/goth/validate/internal"
	"github.com/syb-devs/goth/validate/prefix"
)

var findErrors = internal.FindErrors

func TestHasPrefix(t *testing.T) {
	var tests = []struct {
		input         interface{}
		valid         bool
		errorPatterns map[string][]string
		logicErr      error
	}{
		{
			input: struct {
				Name string `validate:"hasPrefix:foo,bar"`
			}{
				Name: "Food",
			},
			logicErr: prefix.ErrParamCount,
		},
		{
			input: struct {
				Name string `validate:"hasPrefix:foo"`
			}{
				Name: "food",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"hasPrefix:foo"`
			}{
				Name: "",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"hasPrefix:foo"`
			}{
				Name: "Bar",
			},
			errorPatterns: map[string][]string{"Name": []string{"The field Name = Bar should start with foo."}},
		},
	}

	for _, test := range tests {
		v := validate.New()
		res := v.Validate(test.input)
		err := res.LogicError

		if test.logicErr != nil {
			if err != nil && err.Error() != test.logicErr.Error() {
				t.Errorf("expecting logic error: %v, got: %v", test.logicErr, err)
			}
		} else if err != nil {
			t.Errorf(err.Error())
		}
		errs := res.FieldErrors
		if test.valid && errs != nil && errs.Len() > 0 {
			t.Errorf("expecting zero errors, found %s", errs.String())
		}

		if test.errorPatterns != nil {
			if errs == nil {
				t.Errorf("validator did not return any errors, expected: %+v", test.errorPatterns)
			} else {
				findErrors(t, errs, test.errorPatterns)
			}
		}
	}
}
