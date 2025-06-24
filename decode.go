package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

// Decode a map m into struct s. Field/Key names are supposed to be
// case sensitive. Credit: <https://stackoverflow.com/a/26746461>
func mapToStruct(s any, m map[string]any) error {
	for fieldName, v := range m {
		reflectField := reflect.ValueOf(s).Elem()
		reflectFieldValue := reflectField.FieldByNameFunc(func(n string) bool {
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
		if fieldType == reflect.TypeOf(Internal{}) {
			internalSettings := Internal{
				GenerateEPUB: true,
			}

			if _, ok := v.(map[string]any); !ok {
				return fmt.Errorf("internal settings must be defined as a toml map")
			}

			if err := mapToStruct(&internalSettings, v.(map[string]any)); err != nil {
				return err
			}
			value = reflect.ValueOf(internalSettings)
		}

		if fieldType != value.Type() {
			return fmt.Errorf(
				"cannot set field '%v' (%v) to value '%v' (%v) because of mismatch types. Value must be of type %v",
				fieldName, fieldType, value, value.Type(), fieldType)
		}

		reflectFieldValue.Set(value)
	}
	return nil
}

// Decode data and files into a collection of books
func DecodeCollection(data []byte) (Collection, error) {
	c := Collection{
		Internal: Internal{
			GenerateEPUB: true,
		},
		LanguageCode: "en",
	}

	// ---
	// Decode TOML
	// ---
	if _, err := toml.Decode(string(data), &c.Params); err != nil {
		return c, fmt.Errorf("decode: failed to decode toml data. %v", err)
	}

	if err := mapToStruct(&c, c.Params); err != nil {
		return c, fmt.Errorf("decode: failed to decode toml data. %v", err)
	}

	// ---
	// Validate
	// ---
	if err := c.ValidateFields(); err != nil {
		return c, fmt.Errorf("decode: failed to validate fields. %v", err)
	}

	return c, nil
}
