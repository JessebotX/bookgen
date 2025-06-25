package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/yuin/goldmark-meta"
)

var (
	MarkdownToHTML = goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(),
	)
	MarkdownToXHTML = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
	)
)

func convertMarkdownToHTML(content []byte, useXHTML bool) (template.HTML, map[string]any, error) {
	var buffer bytes.Buffer
	context := parser.NewContext()

	var md goldmark.Markdown
	if useXHTML {
		md = MarkdownToXHTML
	} else {
		md = MarkdownToHTML
	}

	if err := md.Convert(content, &buffer, parser.WithContext(context)); err != nil {
		return template.HTML(""), nil, err
	}

	metadata := meta.Get(context)

	return template.HTML(buffer.String()), metadata, nil
}

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
			return fmt.Errorf("value of field `%s` cannot be changed.", fieldName)
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
				"mismatch types: value `%v` (%v) must have the same type as field `%v` (%v).",
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
		return Collection{}, fmt.Errorf("collection: failed to read file `%v`. %w", pathTOML, err)
	}

	// ---
	// Decode TOML
	// ---
	var c Collection
	c.InitializeDefaults()

	if _, err := toml.Decode(string(dataTOML), &c.Params); err != nil {
		return c, fmt.Errorf("collection: failed to decode TOML in `%v`. %w", pathTOML, err)
	}

	if err := mapToStruct(&c, c.Params); err != nil {
		return c, fmt.Errorf("collection: failed to decode TOML in `%v`. %w", pathTOML, err)
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
		return Book{}, fmt.Errorf("book: failed to read file `%v`. %w", err)
	}

	// ---
	// Decode toml
	// ---
	var b Book
	b.InitializeDefaults(workingDir, parent)

	if _, err := toml.Decode(string(dataTOML), &b.Params); err != nil {
		return b, fmt.Errorf("book `%v`: failed to decode TOML in `%v`. %w", b.PageName, pathTOML, err)
	}

	if err := mapToStruct(&b, b.Params); err != nil {
		return b, fmt.Errorf("book `%v`: failed to decode TOML in `%v`. %w", b.PageName, pathTOML, err)
	}

	// ---
	// Check requirements
	// ---
	if err := b.CheckRequirementsForParsing(workingDir); err != nil {
		return b, fmt.Errorf("book `%v`: failed to meet requirements. %w", b.PageName, err)
	}

	// ---
	// Parse markdown
	// ---
	rawMarkdownPath := filepath.Join(workingDir, "index.md")
	rawMarkdown, err := os.ReadFile(rawMarkdownPath)
	if err != nil && os.IsExist(err) {
		return b, fmt.Errorf("book `%v`: failed to read book content file at `%v`, %w", rawMarkdownPath, err)
	}
	b.Content.Raw = string(rawMarkdown)

	contentHTML, _, err := convertMarkdownToHTML(rawMarkdown, false)
	if err != nil {
		return b, fmt.Errorf("book `%v`: failed to convert markdown to HTML. %w", err)
	}
	b.Content.HTML = contentHTML

	// ---
	// TODO: Check existence of files like cover image
	// ---

	// ---
	// TODO: Read chapters
	// ---

	return b, nil
}
