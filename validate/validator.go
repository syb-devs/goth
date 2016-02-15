// Copyright 2014 Simplify Your Bussiness S.L. All rights reserved.

// Package validate implements validation of struct types using rules defined inside struct tags
package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	ErrRuleNotFound       = errors.New("Rule not found")
	ErrStructExpected     = errors.New("The underlying type of the validation data must be struct or *struct")
	ErrUnsupportedType    = errors.New("Unsupported type for rule")
	ErrInvalidParamFormat = errors.New("Invalid format for validation rule parameters")
)

var (
	// TagName is used to specify the struct tag name that the validator should use
	TagName = "validate"

	// RuleSeparator is the character(s) used to separate different validation rules in the same tag
	RuleSeparator = "|"

	defaultValidator = New()
)

// Rule represents a validation rule that will be applied to a struct field value.
type Rule interface {
	Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error)
}

// ErrList is used to store struct validation errors grouped by field name.
type ErrList map[string][]error

// String returns a literal representation of the error list.
func (e ErrList) String() string {
	str := ""
	for field, errors := range e {
		str = str + field + ": "
		for _, err := range errors {
			str = str + err.Error() + ", "
		}
		str = str + "\n"
	}
	return str
}

// Len returns the number of elements in the error list.
func (e ErrList) Len() int {
	return len(e)
}

// ruleMap stores validation rules that will be accessed by its name.
type ruleMap map[string]Rule

// fieldPrefix contains the field prefix for fields of nested structs.
type fieldPrefix []string

// push appends a prefix to the slice.
func (fp *fieldPrefix) push(prefix string) {
	*fp = append(*fp, prefix)
}

// pop returns the last prefix form the slice.
func (fp *fieldPrefix) pop() (p string) {
	*fp, p = (*fp)[:len(*fp)-1], (*fp)[len(*fp)-1]
	return
}

// String returns a literal representation of the field prefix.
func (fp *fieldPrefix) String() string {
	ret := ""
	for _, prefix := range *fp {
		ret = ret + prefix + "."
	}
	return ret
}

// Validator extracts and checks validation rules from struct tags
type Validator struct {
	registeredRules ruleMap
	data            interface{}
	errors          ErrList
	logicError      error
	mu              sync.RWMutex
	tagName         string
	fieldPrefix     fieldPrefix
}

// RegisterRule registers a validation rule in the default validator.
func RegisterRule(name string, rule Rule) {
	defaultValidator.RegisterRule(name, rule)
}

// Validate validates the given struct using the default validator and returns any logic error that might happen.
// To get the actual validation errors, use the method Errors().
func Validate(data interface{}) error {
	return defaultValidator.Validate(data)
}

// SetTagName sets the name of the struct tag to extract validation rules from.
func SetTagName(name string) {
	defaultValidator.SetTagName(name)
}

// Zeroed validator returns a validator with all its fields initialised.
func zeroedValidator() *Validator {
	return &Validator{
		registeredRules: make(ruleMap, 0),
		errors:          make(ErrList, 0),
		fieldPrefix:     make(fieldPrefix, 0),
	}
}

// New returns a new validator, set up with the default rules and options.
func New() *Validator {
	v := zeroedValidator()
	v.tagName = TagName
	v.RegisterRule("length", &lengthRule{})
	v.RegisterRule("regexp", &regexpRule{})

	return v
}

// SetTagName sets the name of the struct tag to extract validation rules from.
func (v *Validator) SetTagName(name string) {
	v.tagName = name
}

// RegisterRule registers a Rule in for this validator under the given name.
func (v *Validator) RegisterRule(name string, rule Rule) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.registeredRules[name] = rule
}

// Copy returns a copy new validator using the same configuration.
func (v *Validator) Copy() *Validator {
	vc := zeroedValidator()
	vc.tagName = v.tagName
	vc.registeredRules = v.registeredRules

	return vc
}

// getRule retrieves a rule from the rule map using a given name.
func (v *Validator) getRule(name string) (Rule, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	r := v.registeredRules[name]
	if r != nil {
		return r, nil
	}
	return nil, ErrRuleNotFound
}

// Validate runs the actual validation of the struct, applying the rules registered in the validator,
// returning any logic error that might happen.
// To get the actual validation errors, use the method Errors().
func (v *Validator) Validate(data interface{}) error {
	sv := reflect.ValueOf(data)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return v.Validate(sv.Elem().Interface())
	}
	if !IsStruct(data) {
		return ErrStructExpected
	}

	v.data = data
	numFields := reflect.ValueOf(v.data).NumField()

	for curField := 0; curField < numFields; curField++ {
		err := v.validateField(curField)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateField validates a single field of the struct and returns a logic error if something goes wrong.
func (v *Validator) validateField(i int) error {

	elem := reflect.TypeOf(v.data).Field(i)
	if !fieldIsExported(elem) {
		return nil
	}
	fieldName := elem.Name

	//TODO: check if field is a pointer
	fieldVal := reflect.ValueOf(v.data).Field(i).Interface()
	if IsStruct(fieldVal) {
		v.fieldPrefix.push(fieldName)
		defer v.fieldPrefix.pop()

		err := v.Validate(fieldVal)

		if err != nil {
			return err
		}
		return nil
	}

	tag := elem.Tag.Get(v.tagName)
	if tag == "" {
		return nil
	}

	for _, ruleStr := range strings.Split(tag, RuleSeparator) {
		var j = strings.Index(tag, ":")
		var ruleParamsStr = ruleStr[j+1:]
		var namedParams map[string]string
		var ruleParams []string

		var ruleName = ruleStr[0:j]

		namedParams = make(map[string]string, 0)

		for _, paramPart := range strings.Split(ruleParamsStr, ",") {
			isNamed := strings.Index(paramPart, ":") != -1
			if isNamed {
				var tmpParam = strings.Split(paramPart, ":")
				if len(tmpParam) != 2 {
					return ErrInvalidParamFormat
				}
				namedParams[tmpParam[0]] = tmpParam[1]
			} else {
				ruleParams = append(ruleParams, paramPart)
			}
		}

		var fieldCheck = func() {
			rule, err := v.getRule(ruleName)
			if err != nil {
				v.logicError = err
				return
			}

			logicErr, inputErr := rule.Validate(v.data, fieldName, ruleParams, namedParams)
			if logicErr != nil {
				v.logicError = logicErr
				return
			}
			if inputErr != nil {
				key := v.fieldPrefix.String() + fieldName
				v.errors[key] = append(v.errors[key], inputErr)
			}
		}

		v.safeExec(fieldCheck)
		if v.logicError != nil {
			return v.logicError
		}
	}
	return nil
}

// Errors returns a list of validation errors.
func (v *Validator) Errors() *ErrList {
	errors := v.errors
	if len(errors) == 0 {
		return nil
	}
	return &errors
}

// ErrorsByField returns a list of validation errors for a given field.
func (v *Validator) ErrorsByField(field string) *[]error {
	if field == "" {
		return nil
	}

	errors := v.errors[field]
	if errors == nil {
		return nil
	}
	return &errors
}

type safeFunc func()

// safeExec executes a given function and stores any recovered panic as a logic error inside de validator.
func (v *Validator) safeExec(f safeFunc) {
	defer func() {
		if recErr := recover(); recErr != nil {
			switch errv := recErr.(type) {
			case string:
				v.logicError = errors.New(errv)
			case error:
				v.logicError = errv
			default:
				v.logicError = fmt.Errorf("Panic recovered with type: %+v", recErr)
			}
		}
	}()
	f()
}

// IsStruct checks if the given value is a struct of a pointer to a struct.
func IsStruct(data interface{}) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		return IsStruct(v.Elem().Interface())
	}
	return v.Kind() == reflect.Struct
}

// fieldIsExported  returns true if the struct field is exported.
func fieldIsExported(f reflect.StructField) bool {
	return len(f.PkgPath) == 0
}

// getInterfaceValue returns the value of a given interface using reflection.
func getInterfaceValue(data interface{}, name string) interface{} {
	return reflect.ValueOf(data).FieldByName(name).Interface()
}

// toString returns a literal representation of a given value.
// The second parameter indicates whether a conversion was possible or not.
func toString(value interface{}) (string, bool) {
	switch v := value.(type) {
	case string, *string, int, *int, int32, *int32, int64, *int64:
		return fmt.Sprintf("%v", v), true
	default:
		return "", false
	}
}

// mustStringify tries to convert the given value to string type and panics if not possible.
func mustStringify(value interface{}) string {
	strVal, ok := toString(value)
	if ok == false {
		panic(ErrUnsupportedType)
	}
	return strVal
}
