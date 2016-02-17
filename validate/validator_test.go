package validate_test

import (
	"reflect"
	"testing"

	"bitbucket.org/syb-devs/goth/validate"
)

func TestValidate(t *testing.T) {
	data := struct {
		Name string `validate:"ruleA:a,b,c,foo:bar|ruleB:1,2,k:v"`
	}{
		Name: "John Doe",
	}

	ruleA := &dummyRule{}
	ruleB := &dummyRule{}
	validate.RegisterRule("ruleA", ruleA)
	validate.RegisterRule("ruleB", ruleB)

	val := validate.New()
	val.Validate(data)

	ruleAExpected := funcData{
		data:        data,
		field:       "Name",
		params:      []string{"a", "b", "c"},
		namedParams: map[string]string{"foo": "bar"},
	}
	ruleBExpected := funcData{
		data:        data,
		field:       "Name",
		params:      []string{"1", "2"},
		namedParams: map[string]string{"k": "v"},
	}

	if !reflect.DeepEqual(ruleAExpected, ruleA.fdata) {
		t.Errorf("expecting function data to be: %#v, but is: %#v", ruleAExpected, ruleA.fdata)
	}
	if !reflect.DeepEqual(ruleBExpected, ruleB.fdata) {
		t.Errorf("expecting function data to be: %#v, but is: %#v", ruleBExpected, ruleB.fdata)
	}
}

type funcData struct {
	data        interface{}
	field       string
	params      []string
	namedParams map[string]string
}

type dummyRule struct {
	fdata funcData
}

func (r *dummyRule) Validate(
	data interface{},
	field string,
	params []string,
	namedParams map[string]string) (errorLogic, errorInput error,
) {
	r.fdata = funcData{
		data:        data,
		field:       field,
		params:      params,
		namedParams: namedParams,
	}
	return nil, nil
}
