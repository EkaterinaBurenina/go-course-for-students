package homework

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) String() []string {
	var res []string
	for _, err := range v {
		res = append(res, err.Err.Error())
	}
	return res
}

func (v ValidationErrors) Error() string {
	return strings.Join(v.String()[:], ",")
}

func Validate(v any) error {
	var validationErrors ValidationErrors

	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	valueOf := reflect.ValueOf(v)

	for i := 0; i < valueOf.NumField(); i++ {
		tag := typeOf.Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		if !typeOf.Field(i).IsExported() {
			validationErrors = append(validationErrors, ValidationError{Err: ErrValidateForUnexportedFields})
			continue
		}

		tagPrefix, tagValue := strings.Split(tag, ":")[0], strings.Split(tag, ":")[1]
		intTagValue, ok := TryToInt(tagValue)

		typeTagValue := "string"
		if ok {
			typeTagValue = "int"
		}

		val := valueOf.Field(i)

		switch tagPrefix {
		case "len":
			intTagVal, err := strconv.Atoi(tagValue)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Err: ErrInvalidValidatorSyntax})
			} else {
				if len(val.String()) != intTagVal {
					validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("length of %s should be %d", typeOf.Field(i).Name, intTagVal)})
				}
			}
		case "min":
			if typeTagValue != "int" {
				validationErrors = append(validationErrors, ValidationError{Err: ErrInvalidValidatorSyntax})
			} else if val.Kind() == reflect.Int {
				if int(val.Int()) < intTagValue {
					validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("length of %s should be greater than %d", typeOf.Field(i).Name, intTagValue)})
				}
			} else if val.Kind() == reflect.String {
				//negative length
				//if int_tag_value < 0 {
				//	validationErrors = append(validationErrors, ValidationError{Err: ErrInvalidValidatorSyntax})
				//}
				if len(val.String()) < intTagValue {
					validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("length of %s should be greater than %d", typeOf.Field(i).Name, intTagValue)})
				}
			} else {
				validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("unsupported type %s", val.Kind())})
			}
		case "max":
			if typeTagValue != "int" {
				validationErrors = append(validationErrors, ValidationError{Err: ErrInvalidValidatorSyntax})
			} else if val.Kind() == reflect.Int {
				if int(val.Int()) > intTagValue {
					validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("length of %s should be less than %d", typeOf.Field(i).Name, intTagValue)})
				}
			} else if val.Kind() == reflect.String {
				//negative length
				//if int_tag_value < 0 {
				//	validationErrors = append(validationErrors, ValidationError{Err: ErrInvalidValidatorSyntax})
				//}
				if len(val.String()) > intTagValue {
					validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("length of %s should be less than %d", typeOf.Field(i).Name, intTagValue)})
				}
			} else {
				validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("unsupported type %s", val.Kind())})
			}
		case "in":
			if tagValue == "" {
				validationErrors = append(validationErrors, ValidationError{Err: ErrInvalidValidatorSyntax})
			}
			if val.Kind() == reflect.String {
				allowedValues := strings.Split(tagValue, ",")
				if !stringInSlice(val.String(), allowedValues) {
					validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("value of %s should be in %s", typeOf.Field(i).Name, tagValue)})
				}
			} else if val.Kind() == reflect.Int {
				contains := false
				allowedValues := strings.Split(tagValue, ",")
				for _, v := range allowedValues {
					intV, ok := TryToInt(v)
					if ok && int(val.Int()) == intV {
						contains = true
						break
					}
				}
				if !contains {
					validationErrors = append(validationErrors, ValidationError{Err: fmt.Errorf("value of %s should be in %s", typeOf.Field(i).Name, tagValue)})
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func TryToInt(s string) (v int, ok bool) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	} else {
		return v, true
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
