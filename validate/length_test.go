package validate_test

import (
	"errors"
	"testing"

	"bitbucket.org/syb-devs/goth/validate"
)

func TestLength(t *testing.T) {
	var tests = []struct {
		input         interface{}
		valid         bool
		errorPatterns map[string][]string
		logicErr      error
	}{
		{
			input: struct {
				Name string `validate:"length:>=,3,zz"`
			}{
				Name: "Jon",
			},
			logicErr: validate.ErrLengthParamCount,
		},
		{
			input: struct {
				Name string `validate:"length:>=,x"`
			}{
				Name: "Jon",
			},
			logicErr: errors.New("strconv.ParseInt: parsing \"x\": invalid syntax"),
		},
		{
			input: struct {
				Name string `validate:"length:3"`
			}{
				Name: "Jon",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"length:=,3"`
			}{
				Name: "Jon",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"length:<,7"`
			}{
				Name: "Basil",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"length:=,3"`
			}{
				Name: "Johnny",
			},
			errorPatterns: map[string][]string{"Name": []string{"The field Name should have a length equal to 3. Actual length: 6"}},
		},
		{
			input: struct {
				Name string `validate:"length:>=,3|length:<=,10"`
			}{
				Name: "Johnny",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"length:x,3"`
			}{
				Name: "Johnny",
			},
			errorPatterns: map[string][]string{"Name": []string{"Invalid operator"}},
		},
	}

	for _, test := range tests {
		v := validate.New()
		err := v.Validate(test.input)
		if test.logicErr != nil {
			if err.Error() != test.logicErr.Error() {
				t.Errorf("expecting logic error: %v, got: %v", test.logicErr, err)
			}
		} else if err != nil {
			t.Errorf(err.Error())
		}
		errs := v.Errors()
		if test.valid && errs != nil && errs.Len() > 0 {
			t.Errorf("expecting zero errors, found %s", errs.String())
		}

		if test.errorPatterns != nil {
			if errs == nil {
				t.Errorf("validator did not return any errors, expected: %+v", test.errorPatterns)
			} else {
				findErrors(t, *errs, test.errorPatterns)
			}
		}
	}
}

func findErrors(t *testing.T, errList validate.ErrList, patterns map[string][]string) {
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
