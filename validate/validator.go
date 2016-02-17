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
	tagName       = "validate"
	ruleSeparator = "|"
	rules         = NewRuleMap()
)

// RegisterRule registers a rule in the default validator
func RegisterRule(name string, rule Rule) {
	rules.RegisterRule(name, rule)
}

// Rule represents a validation rule that will be applied to a struct field value.
type Rule interface {
	Validate(data interface{}, field string, params []string, namedParams map[string]string) (errorLogic, errorInput error)
}

// RuleMap stores validation rules
type RuleMap struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

// NewRuleMap allocates and returns a RuleMap
func NewRuleMap() *RuleMap {
	return &RuleMap{
		rules: make(map[string]Rule, 0),
	}
}

// RegisterRule registers a rule in the map
func (rm *RuleMap) RegisterRule(name string, rule Rule) {
	rm.mu.Lock()
	rm.rules[name] = rule
	rm.mu.Unlock()
}

// GetRule fetches a validation rule from the map
func (rm *RuleMap) GetRule(name string) (Rule, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	rule, ok := rm.rules[name]
	return rule, ok
}

// Result represents the outcome of a validation
type Result struct {
	LogicError error
	FieldErrors
}

// OK returns true if the data validation was succesfull
func (r *Result) OK() bool {
	return r.LogicError == nil && len(r.FieldErrors) == 0
}

// Validator extracts and checks validation rules from struct tags
// TODO(zareone) create a Rule cache?? map[reflect.Type]ruleParams
type Validator struct {
	*RuleMap
	tagName       string
	ruleSeparator string
}

// New returns a new validator, set up with the default rules and options.
func New() *Validator {
	v := zeroedValidator()
	v.tagName = tagName
	v.ruleSeparator = ruleSeparator
	v.RuleMap = rules
	return v
}

func zeroedValidator() *Validator {
	return &Validator{
		RuleMap: NewRuleMap(),
	}
}

// Validate runs the actual validation of the struct, applying the rules registered in the validator,
// returning any logic error that might happen.
// To get the actual validation errors, use the method Errors().
func (v *Validator) Validate(data interface{}) Result {
	result := Result{
		FieldErrors: make(FieldErrors, 0),
	}
	sv := reflect.ValueOf(data)
	for sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		result.LogicError = ErrStructExpected
		return result
	}

	numFields := sv.NumField()
	for curField := 0; curField < numFields; curField++ {
		field := sv.Type().Field(curField)
		if !fieldIsExported(field) {
			continue
		}
		rules := parseRulesTag(field.Name, field.Tag.Get(v.tagName), v.ruleSeparator)

		// fieldValue := sv.Field(curField).Interface()
		fieldErrs, err := v.checkRules(sv.Interface(), rules)
		if err != nil {
			result.LogicError = err
			return result
		}
		result.AppendErrors(field.Name, fieldErrs...)
	}
	return result
}

func (v *Validator) checkRules(data interface{}, rules []ruleParams) ([]error, error) {
	var errs []error
	for _, ruleData := range rules {
		ruleValidator, ok := v.GetRule(ruleData.RuleName)
		if !ok {
			return errs, ErrRuleNotFound
		}

		LogicError, err := ruleValidator.Validate(
			data,
			ruleData.FieldName,
			ruleData.Params,
			ruleData.NamedParams)
		if LogicError != nil {
			return errs, err
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs, nil
}

// parseRulesTag parses a struct tag contents and extracts validation rules from it
func parseRulesTag(fieldName, tag, sep string) []ruleParams {
	var ret []ruleParams
	rules := strings.Split(tag, sep)
	for _, rule := range rules {
		ruleData := newRuleParams()
		ruleData.FieldName = fieldName
		var paramsText string
		unpack(strings.SplitN(rule, ":", 2), &ruleData.RuleName, &paramsText)

		params := strings.Split(paramsText, ",")
		for _, param := range params {
			if strings.Index(param, ":") != -1 {
				var paramName, paramValue string
				unpack(strings.SplitN(param, ":", 2), &paramName, &paramValue)
				ruleData.NamedParams[paramName] = paramValue
			} else {
				ruleData.Params = append(ruleData.Params, param)
			}
		}
		ret = append(ret, *ruleData)
	}
	return ret
}

// unpack stores the contents of a slice of strings in separate strings passed by reference
func unpack(s []string, vars ...*string) {
	for i, str := range s {
		*vars[i] = str
	}
}

// ruleParams holds parameters for a given validation rule
type ruleParams struct {
	RuleName    string
	FieldName   string
	Params      []string
	NamedParams map[string]string
}

func newRuleParams() *ruleParams {
	return &ruleParams{
		NamedParams: map[string]string{},
	}
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

// FieldErrors is used to store struct validation errors grouped by field name.
type FieldErrors map[string][]error

// AppendErrors adds an error associated with the given field
func (e FieldErrors) AppendErrors(field string, errs ...error) {
	if len(errs) == 0 {
		return
	}
	e[field] = append(e[field], errs...)
}

// FieldErrors returns the errors registered for a given field
func (e FieldErrors) FieldErrors(field string) []error {
	return e[field]
}

// String returns a literal representation of the error list.
func (e FieldErrors) String() string {
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
func (e FieldErrors) Len() int {
	return len(e)
}
