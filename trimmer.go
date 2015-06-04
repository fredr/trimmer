package trimmer

import (
	"errors"
	"reflect"
	"strings"
)

// ErrInvalidType is returned when an invalid type is passed to TrimStrings
var ErrInvalidType = errors.New("ptr must be a pointer to a struct")

// TrimStrings will trim whitespace from all string fields in the struct ptr points to. Can be overridden with struct tag `trim:"false"`
func TrimStrings(ptr interface{}) error {
	val := reflect.ValueOf(ptr)

	if val.Kind() != reflect.Ptr {
		return ErrInvalidType
	}

	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrInvalidType
	}

	typ := val.Type()

	// trim string fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.String {
			if fieldType.Tag.Get("trim") != "false" {
				field.SetString(strings.TrimSpace(field.String()))
			}
			continue
		}

		// trim string fields of nested structs
		if field.Kind() == reflect.Struct && field.CanAddr() {
			if err := TrimStrings(field.Addr().Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
