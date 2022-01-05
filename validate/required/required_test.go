package required_test

import (
	"testing"

	"github.com/syb-devs/goth/validate"
	"github.com/syb-devs/goth/validate/internal"
)

var findErrors = internal.FindErrors

func TestRequired(t *testing.T) {
	var tests = []struct {
		input         interface{}
		valid         bool
		errorPatterns map[string][]string
	}{
		{
			input: struct {
				String      string   `validate:"required"`
				Int         int      `validate:"required"`
				Uint        uint     `validate:"required"`
				Float64     float64  `validate:"required"`
				StringSlice []string `validate:"required"`
				StringPtr   *string  `validate:"required"`
			}{
				String:      "s",
				Int:         -10,
				Uint:        10,
				Float64:     0.32,
				StringSlice: []string{"one"},
				StringPtr:   stringPtr("al"),
			},
			valid: true,
		},
		{
			input: struct {
				String      string   `validate:"required"`
				Int         int      `validate:"required"`
				Uint        uint     `validate:"required"`
				Float64     float64  `validate:"required"`
				StringSlice []string `validate:"required"`
				StringPtr   *string  `validate:"required"`
			}{},
			errorPatterns: map[string][]string{
				"String":      []string{"a value is required for String"},
				"Int":         []string{"a value is required for Int"},
				"Uint":        []string{"a value is required for Uint"},
				"Float64":     []string{"a value is required for Float64"},
				"StringSlice": []string{"a value is required for StringSlice"},
				"StringPtr":   []string{"a value is required for StringPtr"},
			},
		},
	}

	for _, test := range tests {
		v := validate.New()
		res := v.Validate(test.input)
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

func stringPtr(str string) *string {
	return &str
}
