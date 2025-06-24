package main

import (
	"reflect"
	"strings"
	"fmt"

	"github.com/BurntSushi/toml"
)

// Decode a map m into struct s. Field/Key names are supposed to be
// case sensitive. Credit: <https://stackoverflow.com/a/26746461>
func mapToStruct(s any, m map[string]any) error {
	for fieldName, v := range m {
		reflectField := reflect.ValueOf(s).Elem()
		reflectFieldValue := reflectField.FieldByNameFunc(func (n string) bool {
			return strings.EqualFold(n, fieldName)
		})

		// ignore keys that don't exist
		if !reflectFieldValue.IsValid() {
			continue
		}

		if !reflectFieldValue.CanSet() {
			return fmt.Errorf("cannot set a value for field '%s'", fieldName)
		}

		fieldType := reflectFieldValue.Type()
		value := reflect.ValueOf(v)
		if fieldType != value.Type() {
			return fmt.Errorf("cannot set field '%v' (%v) to value '%v' (%v) because of mismatch types.\nValue must be of type %v.", fieldName, fieldType, value, value.Type(), fieldType)
		}

		reflectFieldValue.Set(value)
	}
	return nil
}

// Decode data and files into a collection of books
func DecodeCollection (data []byte) (Collection, error) {
	var collection Collection

	if _, err := toml.Decode(string(data), &collection.Params); err != nil {
		return collection, err
	}

	if err := mapToStruct(&collection, collection.Params); err != nil {
		return collection, err
	}

	return collection, nil
}

