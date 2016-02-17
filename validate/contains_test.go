package validate_test

import (
	"testing"

	"bitbucket.org/syb-devs/goth/validate"
)

func TestContains(t *testing.T) {
	var tests = []struct {
		input         interface{}
		valid         bool
		errorPatterns map[string][]string
		logicErr      error
	}{
		{
			input: struct {
				Name string `validate:"contains:foo,bar"`
			}{
				Name: "Food",
			},
			logicErr: validate.ErrContainsParamCount,
		},
		{
			input: struct {
				Name string `validate:"contains:foo"`
			}{
				Name: "food",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"contains:foo"`
			}{
				Name: "Bar",
			},
			errorPatterns: map[string][]string{"Name": []string{"The field Name = Bar should contain foo."}},
		},
	}

	for _, test := range tests {
		v := validate.New()
		res := v.Validate(test.input)
		if test.logicErr != nil {
			if res.LogicError.Error() != test.logicErr.Error() {
				t.Errorf("expecting logic error: %v, got: %v", test.logicErr, res.LogicError)
			}
		} else if res.LogicError != nil {
			t.Errorf(res.LogicError.Error())
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

func findErrors(t *testing.T, errList validate.FieldErrors, patterns map[string][]string) {
	for field, patErrs := range patterns {
		valErrs := errList[field]
		for _, patErr := range patErrs {
			found := false
			for _, valErr := range valErrs {
				if patErr == valErr.Error() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error '%s' for field '%s', but was not found", patErr, field)
			}
		}
	}
}
