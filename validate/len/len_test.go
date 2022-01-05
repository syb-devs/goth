package len_test

import (
	"errors"
	"testing"

	"github.com/syb-devs/goth/validate"
	"github.com/syb-devs/goth/validate/internal"
	"github.com/syb-devs/goth/validate/len"
)

var findErrors = internal.FindErrors

func TestLen(t *testing.T) {
	var tests = []struct {
		input         interface{}
		valid         bool
		errorPatterns map[string][]string
		logicErr      error
	}{
		{
			input: struct {
				Name string `validate:"len:>=,3,zz"`
			}{
				Name: "Jon",
			},
			logicErr: len.ErrParamCount,
		},
		{
			input: struct {
				Name string `validate:"len:>=,x"`
			}{
				Name: "Jon",
			},
			logicErr: errors.New("strconv.Atoi: parsing \"x\": invalid syntax"),
		},
		{
			input: struct {
				Name string `validate:"len:3"`
			}{
				Name: "Jon",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"len:=,3"`
			}{
				Name: "Jon",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"len:<,7"`
			}{
				Name: "Basil",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"len:3"`
			}{
				Name: "Johnny",
			},
			errorPatterns: map[string][]string{"Name": {"The field Name should have a length equal to 3. Actual length: 6"}},
		},
		{
			input: struct {
				Name string `validate:"len:>=,3|len:<=,10"`
			}{
				Name: "Johnny",
			},
			valid: true,
		},
		{
			input: struct {
				Name string `validate:"len:x,3"`
			}{
				Name: "Johnny",
			},
			errorPatterns: map[string][]string{"Name": {"Invalid operator"}},
		},
		{
			input: struct {
				Name string `validate:"len:>,3"`
			}{
				Name: "Jon",
			},
			errorPatterns: map[string][]string{"Name": {"The field Name should have a length greater than 3. Actual length: 3"}},
		},
	}

	for _, test := range tests {
		v := validate.New()
		res := v.Validate(test.input)
		err := res.LogicError
		errs := res.FieldErrors

		if test.logicErr != nil {
			if err.Error() != test.logicErr.Error() {
				t.Errorf("expecting logic error: %v, got: %v", test.logicErr, err)
			}
		} else if err != nil {
			t.Errorf(err.Error())
		}
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
