package validate_test

import (
	"testing"

	"bitbucket.org/syb-devs/goth/validate"
)

func TestValidate(t *testing.T) {
	type data struct {
		Field string `validate:"length:>,4" `
	}

	v := validate.New()
	err := v.Validate(data{})

	if err != nil {
		t.Errorf(err.Error())
	}

	// Validate passing a pointer
	err = v.Validate(&data{})

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestEmptyValidationTag(t *testing.T) {
	type data struct {
		Field string
	}

	v := validate.New()
	err := v.Validate(data{})

	if err != nil {
		t.Errorf("Error during validation")
	}
}

func TestNotStruct(t *testing.T) {
	v := validate.New()
	err := v.Validate("string")

	if err != validate.ErrStructExpected {
		t.Errorf("Expected: %s, got: %s", validate.ErrStructExpected, err)
	}
}

func TestEmbeddedStruct(t *testing.T) {
	type embed struct {
		InnerField string `validate:"length:>,4" `
	}

	type data struct {
		OuterField embed
	}

	v := validate.New()
	err := v.Validate(data{OuterField: embed{InnerField: "foo"}})
	if err != nil {
		t.Error(err.Error())
	}

	errors := v.ErrorsByField("OuterField.InnerField")
	if errors == nil {
		t.Fatalf("No errors retrieved for OuterField.InnerField")
	}
	numErrors := len(*errors)

	if numErrors != 1 {
		t.Errorf("Expected exactly 1 validation error, got %d", numErrors)
	}
}

type foo struct{}

var isStructTests = []struct {
	data     interface{} // input
	expected bool        // expected result
}{
	{"gopher", false},
	{1845, false},
	{foo{}, true},
	{&foo{}, true},
}

func TestIsStruct(t *testing.T) {
	for _, test := range isStructTests {
		actual := validate.IsStruct(test.data)
		if test.expected != actual {
			t.Errorf("IsStruct(%v): expected %v, actual %v", test.data, test.expected, actual)
		}
	}
}
