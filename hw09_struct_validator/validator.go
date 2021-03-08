package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTag = "validate"
	nestedTag   = "nested"
)

var (
	ErrTypeNotImplemented = errors.New("validation for field type isn't implemented")
	ErrStringValueLen     = errors.New("string value has incompatible length")
	ErrStringValueRegexp  = errors.New("string value does not belong to the require regexp")
	ErrStringValueSet     = errors.New("string value does not belong to the require set")
	ErrIntValueMin        = errors.New("int value has incompatible min")
	ErrIntValueMax        = errors.New("int value does not belong to the require max")
	ErrIntValueRange      = errors.New("int value does not belong to the require range")
	ErrTagParseTwoInts    = errors.New("failed to parse validation value with 2 int")
	ErrParseTag           = fmt.Errorf("failed to parse tag")

	errTagValueNotFound = errors.New("validation value not found")
)

type ValidationError struct {
	Field string
	Err   error
}

type (
	ValidationErrors []ValidationError
	arrayErrors      []error
)

type validationTag struct {
	values map[string]string
}

type intValidationTag struct {
	base validationTag
}

type stringValidationTag struct {
	base validationTag
}

func (v ValidationErrors) Error() string {
	s := make([]string, 0, len(v)*5+1)
	s = append(s, "{")

	for i, e := range v {
		errStr := ""
		if e.Err != nil {
			errStr = e.Err.Error()
		}
		s = append(s, "\"", e.Field, "\": ", errStr)

		if i != len(v)-1 {
			s = append(s, "; ")
		}
	}

	s = append(s, "}")

	return strings.Join(s, "")
}

func (a arrayErrors) Error() string {
	s := make([]string, 0, len(a)*2+1)
	s = append(s, "[")

	for i, e := range a {
		errStr := ""
		if e != nil {
			errStr = e.Error()
		}
		s = append(s, errStr)
		if i != len(a)-1 {
			s = append(s, "; ")
		}
	}

	s = append(s, "]")

	return strings.Join(s, "")
}

func (t validationTag) getStringValue(key string) (string, error) {
	if v, ok := t.values[key]; ok {
		return v, nil
	}

	return "", errTagValueNotFound
}

func (t validationTag) getIntValue(key string) (int64, error) {
	if v, ok := t.values[key]; ok {
		m, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return m, nil
	}

	return 0, errTagValueNotFound
}

func (t validationTag) getIntPairValue(key string) ([]int64, error) {
	if v, ok := t.values[key]; ok {
		vals := strings.Split(v, ",")
		if len(vals) != 2 {
			return nil, ErrTagParseTwoInts
		}
		res := make([]int64, 2)
		var err error
		res[0], err = strconv.ParseInt(vals[0], 10, 64)
		if err != nil {
			return nil, err
		}
		res[1], err = strconv.ParseInt(vals[1], 10, 64)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, errTagValueNotFound
}

func (t validationTag) getStringArrayValue(key string) ([]string, error) {
	if v, ok := t.values[key]; ok {
		vals := strings.Split(v, ",")
		return vals, nil
	}

	return nil, errTagValueNotFound
}

func (t intValidationTag) getMax() (int64, error) {
	return t.base.getIntValue("max")
}

func (t intValidationTag) getMin() (int64, error) {
	return t.base.getIntValue("min")
}

func (t intValidationTag) getRange() ([]int64, error) {
	return t.base.getIntPairValue("in")
}

func (t stringValidationTag) getLen() (int64, error) {
	return t.base.getIntValue("len")
}

func (t stringValidationTag) getSet() ([]string, error) {
	return t.base.getStringArrayValue("in")
}

func (t stringValidationTag) getRegexp() (string, error) {
	return t.base.getStringValue("regexp")
}

func getValidationTag(tag string) (validationTag, error) {
	var res validationTag
	vals := strings.Split(tag, "|")
	res.values = make(map[string]string, len(vals))
	for _, v := range vals {
		pair := strings.Split(v, ":")
		if len(pair) != 2 {
			return res, ErrParseTag
		}
		res.values[pair[0]] = pair[1]
	}

	return res, nil
}

func getStringValidationTag(tag string) (stringValidationTag, error) {
	base, err := getValidationTag(tag)
	return stringValidationTag{base}, err
}

func getIntValidationTag(tag string) (intValidationTag, error) {
	base, err := getValidationTag(tag)
	return intValidationTag{base}, err
}

func checkStringLen(v string, strValTag stringValidationTag) error {
	l, err := strValTag.getLen()
	if err != nil {
		if !errors.Is(err, errTagValueNotFound) {
			return err
		}
	} else {
		if len(v) != int(l) {
			return ErrStringValueLen
		}
	}

	return nil
}

func checkStringRegexp(v string, strValTag stringValidationTag) error {
	r, err := strValTag.getRegexp()
	if err != nil {
		if !errors.Is(err, errTagValueNotFound) {
			return err
		}
	} else {
		res, err := regexp.MatchString(r, v)
		if err != nil {
			return err
		}
		if res {
			return nil
		}
		return ErrStringValueRegexp
	}

	return nil
}

func checkStringSet(v string, strValTag stringValidationTag) error {
	set, err := strValTag.getSet()
	if err != nil {
		if !errors.Is(err, errTagValueNotFound) {
			return err
		}
	} else {
		for _, s := range set {
			if v == s {
				return nil
			}
		}
		return ErrStringValueSet
	}

	return nil
}

func checkIntMax(v int64, intValTag intValidationTag) error {
	m, err := intValTag.getMax()
	if err != nil {
		if !errors.Is(err, errTagValueNotFound) {
			return err
		}
	} else if v > m {
		return ErrIntValueMax
	}

	return nil
}

func checkIntMin(v int64, intValTag intValidationTag) error {
	m, err := intValTag.getMin()
	if err != nil {
		if !errors.Is(err, errTagValueNotFound) {
			return err
		}
	} else if v < m {
		return ErrIntValueMin
	}

	return nil
}

func checkIntRange(v int64, intValTag intValidationTag) error {
	m, err := intValTag.getRange()
	if err != nil {
		if !errors.Is(err, errTagValueNotFound) {
			return err
		}
	} else if v < m[0] || v > m[1] {
		return ErrIntValueRange
	}

	return nil
}

func validateString(v string, tag string) error {
	strValTag, err := getStringValidationTag(tag)
	if err != nil {
		return err
	}

	var res arrayErrors

	err = checkStringLen(v, strValTag)
	if err != nil {
		res = append(res, err)
	}

	err = checkStringRegexp(v, strValTag)
	if err != nil {
		res = append(res, err)
	}

	err = checkStringSet(v, strValTag)
	if err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

func validateInt(v int64, tag string) error {
	intValTag, err := getIntValidationTag(tag)
	if err != nil {
		return err
	}

	var res arrayErrors

	err = checkIntMax(v, intValTag)
	if err != nil {
		res = append(res, err)
	}

	err = checkIntMin(v, intValTag)
	if err != nil {
		res = append(res, err)
	}

	err = checkIntRange(v, intValTag)
	if err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

func checkNestedTag(tag string) bool {
	return tag == nestedTag
}

func validateField(v reflect.Value, tag string) error {
	switch v.Kind() {
	case reflect.Struct:
		if checkNestedTag(tag) {
			return Validate(v.Interface())
		}

		return nil
	case reflect.String:
		return validateString(v.String(), tag)
	case reflect.Int:
		return validateInt(v.Int(), tag)
	case reflect.Array:
	case reflect.Slice:
		var res ValidationErrors
		for i := 0; i < v.Len(); i++ {
			if err := validateField(v.Index(i), tag); err != nil {
				res = append(res, ValidationError{strconv.Itoa(i), err})
			}
		}

		if len(res) > 0 {
			return res
		}

		return nil
	}

	return ErrTypeNotImplemented
}

func validateImpl(v reflect.Value, tag reflect.StructTag) error {
	vTag, ok := tag.Lookup(validateTag)
	if !ok {
		return nil
	}

	return validateField(v, vTag)
}

func Validate(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but %T was received", v)
	}

	t := v.Type()
	var res ValidationErrors
	for i := 0; i < t.NumField(); i++ {
		fV := v.Field(i)
		if fV.CanInterface() {
			fT := t.Field(i)
			err := validateImpl(fV, fT.Tag)
			if err != nil {
				res = append(res, ValidationError{fT.Name, err})
			}
		}
	}

	if len(res) > 0 {
		return res
	}

	return nil
}
