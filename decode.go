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
	for fieldName, unparsedValue := range m {
		reflectField := reflect.ValueOf(s).Elem()
		reflectFieldValue := reflectField.FieldByNameFunc(func(n string) bool {
			return strings.EqualFold(n, fieldName)
		})

		// ignore keys that don't exist
		if !reflectFieldValue.IsValid() {
			continue
		}

		if !reflectFieldValue.CanSet() {
			return fmt.Errorf("value of field '%s' cannot be changed.", fieldName)
		}

		fieldType := reflectFieldValue.Type()
		fieldValue := reflect.ValueOf(unparsedValue)
		if fieldType == reflect.TypeOf(Internal{}) {
			internalSettings := Internal{
				GenerateEPUB: true,
			}

			if _, ok := unparsedValue.(map[string]any); !ok {
				return fmt.Errorf("internal settings format invalid. Ensure it is in the form of a map of key/value pairs.")
			}

			if err := mapToStruct(&internalSettings, unparsedValue.(map[string]any)); err != nil {
				return err
			}
			fieldValue = reflect.ValueOf(internalSettings)
		}

		if fieldType != fieldValue.Type() {
			return fmt.Errorf(
				"mismatch types: value '%v' (%v) must have the same type as field '%v' (%v).",
				fieldValue, fieldValue.Type(), fieldName, fieldType)
		}

		reflectFieldValue.Set(fieldValue)
	}
	return nil
}

// Decode a structured directory with a bookgen configuration file
// into a Collection.
func DecodeCollection(workingDir string) (Collection, error) {
	// ---
	// Read file
	// ---
	pathTOML := filepath.Join(workingDir, "bookgen.toml")
	dataTOML, err := os.ReadFile(pathTOML)
	if err != nil {
		return Collection{}, fmt.Errorf("collection: failed to read file '%v'. %w", pathTOML, err)
	}

	// ---
	// Decode TOML
	// ---
	c := Collection{
		Internal: Internal{
			GenerateEPUB: true,
		},
		LanguageCode: "en",
	}

	if _, err := toml.Decode(string(dataTOML), &c.Params); err != nil {
		return c, fmt.Errorf("collection: failed to decode TOML in '%v'. %w", pathTOML, err)
	}

	if err := mapToStruct(&c, c.Params); err != nil {
		return c, fmt.Errorf("collection: failed to decode TOML in '%v'. %w", pathTOML, err)
	}

	// ---
	// Check requirements
	// ---
	if err := c.CheckRequirementsForParsing(); err != nil {
		return c, fmt.Errorf("collection: failed to meet requirements. %w", err)
	}

	// ---
	// Decode books
	// ---
	c.Books = make([]Book, 0)
	booksDir := filepath.Join(workingDir, "books")
	items, err := os.ReadDir(booksDir)
	if err != nil {
		if os.IsNotExist(err) { // no error, do nothing
			return c, nil
		}

		// error
		return c, fmt.Errorf("collection: failed to read books directory %v. %w", booksDir, err)
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		bookWorkingDir := filepath.Join(booksDir, item.Name())
		book, err := DecodeBook(bookWorkingDir, &c)
		if err != nil {
			return c, err
		}

		c.Books = append(c.Books, book)
	}

	return c, nil
}

// Decode a structured directory with a bookgen-book configuration
// file into a Book.
func DecodeBook(workingDir string, parent *Collection) (Book, error) {
	// ---
	// Read file
	// ---
	pathTOML := filepath.Join(workingDir, "bookgen-book.toml")
	dataTOML, err := os.ReadFile(pathTOML)
	if err != nil {
		return Book{}, fmt.Errorf("book: failed to read file '%v'. %w", err)
	}

	// ---
	// Decode toml
	// ---
	b := Book{
		Parent:   parent,
		PageName: filepath.Base(workingDir),
	}

	if parent != nil {
		b.Internal.GenerateEPUB = parent.Internal.GenerateEPUB
		b.LanguageCode = parent.LanguageCode
	}

	if _, err := toml.Decode(string(dataTOML), &b.Params); err != nil {
		return b, fmt.Errorf("book '%v': failed to decode TOML in '%v'. %w", b.PageName, pathTOML, err)
	}

	if err := mapToStruct(&b, b.Params); err != nil {
		return b, fmt.Errorf("book '%v': failed to decode TOML in '%v'. %w", b.PageName, pathTOML, err)
	}

	// ---
	// Check requirements
	// ---
	if err := b.CheckRequirementsForParsing(workingDir); err != nil {
		return b, fmt.Errorf("book '%v': failed to meet requirements. %w", b.PageName, err)
	}

	return b, nil
}
