package main

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	mapstructure "github.com/go-viper/mapstructure/v2"

	yaml "github.com/goccy/go-yaml"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	MarkdownToHTML = goldmark.New(
		goldmark.WithExtensions(
			Meta,
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

// Decode a structured directory with a bookgen configuration file
// into a Collection.
func DecodeCollection(workingDir string) (Collection, error) {
	// ---
	// Read file
	// ---
	pathConfig := filepath.Join(workingDir, "bookgen.yml")
	dataConfig, err := os.ReadFile(pathConfig)
	if err != nil {
		return Collection{}, fmt.Errorf("collection: failed to read file `%v`. %w", pathConfig, err)
	}

	// ---
	// Decode config
	// ---
	var c Collection
	c.InitializeDefaults()

	if err := yaml.Unmarshal(dataConfig, &c.Params); err != nil {
		return c, fmt.Errorf("collection: failed to decode YAML in `%v`. %w", pathConfig, err)
	}

	if err := mapstructure.Decode(c.Params, &c); err != nil {
		return c, fmt.Errorf("collection: failed to decode YAML in `%v`. %w", pathConfig, err)
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
	pathConfig := filepath.Join(workingDir, "bookgen-book.yml")
	dataConfig, err := os.ReadFile(pathConfig)
	if err != nil {
		return Book{}, fmt.Errorf("book: failed to read file `%v`. %w", err)
	}

	// ---
	// Decode config
	// ---
	var b Book
	b.InitializeDefaults(workingDir, parent)

	if err := yaml.Unmarshal(dataConfig, &b.Params); err != nil {
		return b, fmt.Errorf("book `%v`: failed to decode YAML in `%v`. %w", b.PageName, pathConfig, err)
	}

	if err := mapstructure.Decode(b.Params, &b); err != nil {
		return b, fmt.Errorf("book `%v`: failed to decode YAML in `%v`. %w", b.PageName, pathConfig, err)
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
		return b, fmt.Errorf("book `%v`: failed to read book content file at `%v`, %w", b.PageName, rawMarkdownPath, err)
	}
	b.Content.Raw = string(rawMarkdown)

	contentHTML, _, err := convertMarkdownToHTML(rawMarkdown, false)
	if err != nil {
		return b, fmt.Errorf("book `%v`: failed to convert markdown to HTML. %w", b.PageName, err)
	}
	b.Content.HTML = contentHTML

	datePubParam, ok := b.Params["published"]
	if ok && b.DatePublished.IsZero() {
		b.DatePublished, err = getTimeFromParam(datePubParam)
		if err != nil {
			return b, fmt.Errorf("book `%v`: failed to parse date published: %w", b.PageName, err)
		}
	}

	dateModParam, ok := b.Params["modified"]
	if ok && b.DateModified.IsZero() {
		b.DateModified, err = getTimeFromParam(dateModParam)
		if err != nil {
			return b, fmt.Errorf("book `%v`: failed to parse date modified: %w", b.PageName, err)
		}
	}

	// ---
	// TODO: Check existence of files like cover image
	// ---

	// ---
	// Read chapters
	// ---
	chaptersDir := filepath.Join(workingDir, "chapters")
	items, err := os.ReadDir(chaptersDir)
	if err != nil {
		if os.IsExist(err) {
			return b, fmt.Errorf("book `%v`: failed to read chapters directory at `%v`. %w", b.PageName, chaptersDir, err)
		}
	}

	b.Chapters = make([]Chapter, 0)
	for _, item := range items {
		if item.IsDir() || !strings.HasSuffix(item.Name(), ".md") {
			continue
		}

		chapterSourcePath := filepath.Join(chaptersDir, item.Name())
		c, err := DecodeChapter(chapterSourcePath, &b)
		if err != nil {
			return b, fmt.Errorf("book %v: %w", b.PageName, err)
		}

		b.Chapters = append(b.Chapters, c)
	}

	// Sort chapters and fill Next and Previous pointers
	slices.SortFunc(b.Chapters, func(x, y Chapter) int {
		// Sort order: Order, Title.
		// TODO: compare other fields such as DatePublished.
		if n := cmp.Compare(x.Order, y.Order); n != 0 {
			return n
		}

		return strings.Compare(x.Title, y.Title)
	})

	for i := range len(b.Chapters) {
		if (i - 1) >= 0 {
			b.Chapters[i].Previous = &b.Chapters[i-1]
		}

		if (i + 1) < len(b.Chapters) {
			b.Chapters[i].Next = &b.Chapters[i+1]
		}
	}

	return b, nil
}

// Decode file path with .md extension into a Chapter.
func DecodeChapter(path string, parent *Book) (Chapter, error) {
	if filepath.Ext(path) != ".md" {
		return Chapter{}, fmt.Errorf("chapter %v: missing `.md` (markdown) file extension", filepath.Base(path), path)
	}

	var c Chapter
	c.InitializeDefaults(path, parent)

	rawMarkdown, err := os.ReadFile(path)
	if err != nil {
		return Chapter{}, fmt.Errorf("chapter `%v`: failed to read file at `%v`. %w", c.PageName, path, err)
	}

	c.Content.Raw = string(rawMarkdown)
	contentHTML, metadata, err := convertMarkdownToHTML(rawMarkdown, false)
	if err != nil {
		return c, fmt.Errorf("chapter `%v`: failed to convert markdown to HTML. %w", c.PageName, err)
	}
	c.Content.HTML = contentHTML

	c.Params = metadata
	if err := mapstructure.Decode(c.Params, &c); err != nil {
		return c, fmt.Errorf("chapter `%v`: failed to decode metadata in chapter. %w", c.PageName, err)
	}

	datePubParam, ok := c.Params["published"]
	if ok && c.DatePublished.IsZero() {
		c.DatePublished, err = getTimeFromParam(datePubParam)
		if err != nil {
			return c, fmt.Errorf("chapter `%v`: failed to parse date published: %w", c.PageName, err)
		}
	}

	dateModParam, ok := c.Params["modified"]
	if ok && c.DateModified.IsZero() {
		c.DateModified, err = getTimeFromParam(dateModParam)
		if err != nil {
			return c, fmt.Errorf("chapter `%v`: failed to parse date modified: %w", c.PageName, err)
		}
	}

	return c, nil
}

// Assumes param exists.
func getTimeFromParam(param any) (time.Time, error) {
	switch v := param.(type) {
	case time.Time:
		return v, nil
	case string:
		date, err := stringToTime(v)
		if err != nil {
			return time.Time{}, err
		}
		return date, nil
	default:
		return time.Time{}, fmt.Errorf("incorrect parameter type given for date format. Must be either `string` or `time.Time`.")
	}
}

// Convert a date string input into time. String argument must be in
// the correct format or it will return an error.
func stringToTime(sTime string) (time.Time, error) {
	formats := []string{
		"2006",
		"2006-01",
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04Z07:00",
		"2006-01-02T15:04Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04:05Z07:00",
	}

	var errs error
	for _, format := range formats {
		date, err := time.Parse(format, sTime)
		if err == nil {
			return date, nil
		}
		errs = errors.Join(errs, err)
	}

	return time.Time{}, fmt.Errorf("date string `%v` does not match any of the following formats:\n%w", sTime, errs)
}

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

	metadata := Get(context)

	return template.HTML(buffer.String()), metadata, nil
}
