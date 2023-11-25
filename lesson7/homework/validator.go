package homework

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
var ErrInvalidTag = errors.New("invalid tag provided for field")

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := ""
	for i, err := range v {
		res += err.Err.Error()

		if i != len(v)-1 {
			res += "\n"
		}
	}

	return res
}

func (v *ValidationErrors) Add(errors ...ValidationError) {
	*v = append(*v, errors...)
}

func getAllFields(val reflect.Value) ([]reflect.Value, []reflect.StructField) {

	n := val.Type().NumField()

	fieds := make([]reflect.Value, n)
	structFieds := make([]reflect.StructField, n)
	for i := 0; i < n; i++ {
		fieds[i] = val.Field(i)
		structFieds[i] = val.Type().Field(i)
	}

	return fieds, structFieds
}

// todo: cover with tests

type elem interface {
	int | string
}

func sliceToString[E elem](s []E) string {

	res := ""
	for i, v := range s {

		res += fmt.Sprint(v)
		if i != len(s)-1 {
			res += " "
		}
	}

	return res
}

func handleLenTag(val string, field reflect.Value, structField reflect.StructField) error {

	constraint, err := strconv.Atoi(val)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}

	// check if field has string type
	if field.Type() != reflect.TypeOf("") {
		return ErrInvalidTag
	}

	if len(field.String()) != constraint {
		return ValidationError{
			Err: errors.New("length of field " + structField.Name + " is not " + strconv.Itoa(constraint)),
		}
	}
	return nil
}

func handleMinTag(val string, field reflect.Value, structField reflect.StructField) error {

	constraint, err := strconv.Atoi(val)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}

	// check if field has type string or int
	if field.Type() == reflect.TypeOf("") {
		if len(field.String()) < constraint {
			return ValidationError{
				Err: errors.New("length of field " + structField.Name + " is less than " + strconv.Itoa(constraint)),
			}
		}
	} else if field.Type() == reflect.TypeOf(constraint) {
		if field.Int() < int64(constraint) {
			return ValidationError{
				Err: errors.New("field " + structField.Name + " is less than " + strconv.Itoa(constraint)),
			}
		}
	} else {
		return ErrInvalidTag
	}

	return nil
}

func handleMaxTag(val string, field reflect.Value, structField reflect.StructField) error {

	constraint, err := strconv.Atoi(val)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}

	// check if field has type string or int
	if field.Type() == reflect.TypeOf("") {
		if len(field.String()) > constraint {
			return ValidationError{
				Err: errors.New("length of field " + structField.Name + " is more than " + strconv.Itoa(constraint)),
			}
		}
	} else if field.Type() == reflect.TypeOf(constraint) {
		if field.Int() > int64(constraint) {
			return ValidationError{
				Err: errors.New("field " + structField.Name + " is more than " + strconv.Itoa(constraint)),
			}
		}
	} else {
		return ErrInvalidTag
	}

	return nil
}

func handleInTag(val string, field reflect.Value, structField reflect.StructField) error {

	spltd := strings.Split(val, ",")

	if len(spltd) == 1 && spltd[0] == "" {
		return ValidationError{
			Err: errors.New("in tag cannot contain empty constraint"),
		}
	}

	isInt := true

	constraints := make([]int, len(spltd))
	for i, v := range spltd {

		out, err := strconv.Atoi(v)
		if err != nil {
			isInt = false
		} else {
			constraints[i] = out
		}
	}

	if isInt {

		if field.Type() != reflect.TypeOf(1) {
			return ErrInvalidTag
		}

		if !slices.Contains(constraints, int(field.Int())) {
			return ValidationError{
				Err: errors.New("field " + structField.Name + " is not in " + sliceToString(constraints)),
			}
		}

	} else {

		if field.Type() != reflect.TypeOf("") {
			return ValidationError{
				Err: ErrInvalidTag,
			}
		}

		if !slices.Contains(spltd, field.String()) {
			return ValidationError{
				Err: errors.New("field " + structField.Name + " is not in " + sliceToString(spltd)),
			}
		}
	}

	return nil
}

func handleTag(tag string, field reflect.Value, structField reflect.StructField) error {

	spltd := strings.Split(tag, ":")

	key, value := spltd[0], spltd[1]

	switch key {

	case "len":
		return handleLenTag(value, field, structField)
	case "min":
		return handleMinTag(value, field, structField)
	case "max":
		return handleMaxTag(value, field, structField)
	case "in":
		return handleInTag(value, field, structField)
	}

	return nil
}

func Validate(s any) error {

	validErrors := make(ValidationErrors, 0)

	val := reflect.ValueOf(s)
	if val.Type().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	fields, structFields := getAllFields(val)
	for i, v := range structFields {

		tag, exists := v.Tag.Lookup("validate")
		if exists {

			if !v.IsExported() {
				validErrors.Add(ValidationError{Err: ErrValidateForUnexportedFields})
				continue
			}

			err := handleTag(tag, fields[i], v)
			if err != nil {

				if errors.Is(err, ValidationError{}) {
					validErrors.Add(err.(ValidationError))
				} else {
					validErrors.Add(ValidationError{Err: err})
				}
			}
		}

	}

	if len(validErrors) == 0 {
		return nil
	} else {
		return validErrors
	}
}
