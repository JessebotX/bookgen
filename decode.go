package main

import (
	"fmt"
	"os"
	"path/filepath"
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
func DecodeCollection(data []byte, workingDir string) (Collection, error) {
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
		return c, fmt.Errorf("failed to decode collection toml data. %v", err)
	}

	if err := mapToStruct(&c, c.Params); err != nil {
		return c, fmt.Errorf("failed to decode collection toml data. %v", err)
	}

	// ---
	// Validate
	// ---
	if err := c.ValidateFields(); err != nil {
		return c, fmt.Errorf("failed to validate collection fields. %v", err)
	}

	// ---
	// Decode books
	// ---
	c.Books = make([]Book, 0)
	booksDir := filepath.Join(workingDir, "books")
	items, err := os.ReadDir(booksDir)
	if err != nil && os.IsNotExist(err) {
		return c, nil
	}

	if err != nil {
		return c, fmt.Errorf("failed to read books directory %v. %v", booksDir, err)
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		bookWorkingDir := filepath.Join(booksDir, item.Name())
		tomlBody, err := os.ReadFile(filepath.Join(bookWorkingDir, "bookgen-book.toml"))
		if err != nil {
			return c, err
		}

		book, err := DecodeBook(tomlBody, bookWorkingDir, &c)
		if err != nil {
			return c, err
		}

		c.Books = append(c.Books, book)
	}

	return c, nil
}

func DecodeBook(data []byte, workingDir string, parent *Collection) (Book, error) {
	b := Book{
		Parent: parent,
	}

	if parent != nil {
		b.Internal.GenerateEPUB = parent.Internal.GenerateEPUB
		b.LanguageCode = parent.LanguageCode
	}

	if _, err := toml.Decode(string(data), &b.Params); err != nil {
		return b, fmt.Errorf("failed to decode book toml data @ %v. %v", workingDir, err)
	}

	if err := mapToStruct(&b, b.Params); err != nil {
		return b, fmt.Errorf("failed to decode book toml data @ %v. %v", workingDir, err)
	}

	return b, nil
}
