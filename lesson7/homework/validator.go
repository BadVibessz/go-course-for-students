package homework

import (
	"github.com/pkg/errors"
	"reflect"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	// TODO: implement this
	return ""
}

func getAllFields(val reflect.Value) ([]reflect.Value, []reflect.StructField) {

	fieds := make([]reflect.Value, val.NumField())
	structFieds := make([]reflect.StructField, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fieds = append(fieds, val.Field(i))
		structFieds = append(structFieds, val.Type().Field(i))
	}

	return fieds, structFieds
}

// todo: cover with tests
func hasUnexportedFields(fields []reflect.StructField) bool {
	for _, v := range fields {
		if !v.IsExported() {
			return true
		}
	}
	return false
}

func Validate(s any) error {
	// TODO: implement this

	val := reflect.ValueOf(&s).Elem()

	fields, structFields := getAllFields(val)

	if hasUnexportedFields(structFields) {
		return ErrValidateForUnexportedFields
	}

	for _, v := range structFields {

		tagValue, exists := v.Tag.Lookup("len")
		if exists {
			// todo: validate syntax
		}

	}

	print(structFields[0].Tag)

	return nil
}
