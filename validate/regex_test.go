package validate_test

import (
	"testing"

	"bitbucket.org/syb-devs/goth/validate"
)

func TestRegexMatchKO(t *testing.T) {
	type data struct {
		Field string `validate:"regexp:^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$,allowEmpty:1" `
	}

	v := validate.New()
	err := v.Validate(data{Field: "foo"})

	if err != nil {
		t.Errorf(err.Error())
	}

	if v.Errors() == nil || v.Errors().Len() != 1 {
		t.Errorf("Expecting exactly one validation error")
	}
}

func TestRegexCompileKO(t *testing.T) {
	type data struct {
		Field string `validate:"regexp:((,allowEmpty:1" `
	}

	v := validate.New()
	err := v.Validate(data{Field: "foo"})

	if err == nil {
		t.Errorf("Expecting compile regexp error")
	}

	if v.Errors() != nil {
		t.Errorf("Expecting compile regexp error, input errors must be void")
	}
}

func TestRegexMatchOK(t *testing.T) {
	type data struct {
		Field string `validate:"regexp:^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$,allowEmpty:1" `
	}

	v := validate.New()
	err := v.Validate(data{Field: "foo@mail.com"})

	if err != nil {
		t.Errorf(err.Error())
	}

	if v.Errors() != nil {
		t.Errorf("Unexpected validation error: %s", v.Errors())
	}
}

func TestRegexCompileOK(t *testing.T) {
	type data struct {
		Field string `validate:"regexp:^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$,allowEmpty:1" `
	}

	v := validate.New()
	err := v.Validate(data{Field: "foo@mail.com"})

	if err != nil {
		t.Errorf("This regexp was supposed to compile")
	}

	if v.Errors() != nil {
		t.Errorf("No input errors were expected because there is a regexp compile error")
	}
}

func TestRegexAllowEmptyOK(t *testing.T) {
	type data1 struct {
		Field string `validate:"regexp:^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$,allowEmpty:1" `
	}

	v := validate.New()
	err := v.Validate(data1{Field: ""})
	if err != nil {
		t.Errorf("No logic error was expected, allowEmpty:1")
	}

	if v.Errors() != nil {
		t.Errorf("No input errors were expected, allowEmpty:1")
	}

	type data2 struct {
		Field string `validate:"regexp:^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$,allowEmpty:0" `
	}

	v = validate.New()
	err = v.Validate(data2{Field: ""})

	if err != nil {
		t.Errorf("No logic error was expected, allowEmpty:0")
	}

	if v.Errors() == nil || v.Errors().Len() != 1 {
		t.Errorf("Expecting exactly one validation error, allowEmpty:0")
	}

}
