package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var (
	StatusValidValues = []string{"completed", "hiatus", "ongoing"}
)

// Internal represents the app's settings that may be useful
// for themes to know about.
type Internal struct {
	GenerateEPUB bool
}

// Collection represents a list/index of one or more books.
type Collection struct {
	Params              map[string]any
	Internal            Internal
	Title               string
	Description         string
	BaseURL             string
	LanguageCode        string
	Books               []Book
	FaviconImageName    string
	ConfigFormatVersion int
}

// Close properly deallocates elements in the Collection object such
// as maps. Also closes all Books as well.
func (c Collection) Close() {
	clear(c.Params)

	for _, b := range c.Books {
		b.Close()
	}
}

// ValidateFields returns nil if the collection fields are
// well-formed, otherwise it errors. Note that it does not validate
// the books found in Collection.Books, nor does it verify the
// existence of any files specified in the config.
func (c Collection) ValidateFields() error {
	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("collection: missing/empty required field 'title'")
	}

	if strings.TrimSpace(c.LanguageCode) == "" {
		return fmt.Errorf("collection: missing/empty required field 'languageCode'")
	}

	return nil
}

// Book represents an ordered list of chapters.
type Book struct {
	Params           map[string]any
	Parent           *Collection
	Internal         Internal
	PageName         string
	Title            string
	Subtitle         string
	TitleSort        string
	Authors          []Author
	AuthorsSort      string
	Series           Series
	Description      string
	Copyright        string
	IDs              []string
	Tags             []string
	CoverImageName   string
	FaviconImageName string
	Status           string
	LanguageCode     string
	DatePublished    time.Time
	DateLastModified time.Time
	Content          Content
	IsStub           bool
	Chapters         []Chapter
}

// ValidateFields returns nil if the book fields are well-formed,
// otherwise it errors. Note that it does not validate the chapters
// found in Book.Chapters, nor does it verify the existence of any
// files specified in the config.
func (b Book) ValidateFields(workingDir string) error {
	if b.PageName != filepath.Base(workingDir) {
		return fmt.Errorf("book: cannot set field 'pageName' to anything other than the base name of the working directory")
	}

	if strings.TrimSpace(b.Title) == "" {
		return fmt.Errorf("book: missing/empty required field 'title'")
	}

	if strings.TrimSpace(b.Status) != "" {
		if !slices.Contains(StatusValidValues, strings.ToLower(b.Status)) {
			return fmt.Errorf("book: invalid value for 'status' field. Must be one of the following: %v", StatusValidValues)
		}
	}

	return nil
}

// Close properly deallocates elements in the Book object such as
// maps. Also closes all Chapters as well.
func (b Book) Close() {
	clear(b.Params)

	for _, c := range b.Chapters {
		c.Close()
	}
}

// Chapter represents a division in a Book, primarily containing the
// book's text content.
type Chapter struct {
	Params           map[string]any
	Parent           *Book
	Previous         *Chapter
	Next             *Chapter
	PageName         string
	Title            string
	Subtitle         string
	Description      string
	Order            int
	Authors          []Author
	Copyright        string
	LanguageCode     string
	DatePublished    time.Time
	DateLastModified time.Time
	Content          Content
}

// Close properly deallocates any elements in the Chapter object such
// as maps.
func (c Chapter) Close() {
	clear(c.Params)
}

// Author represents an individual writer or contributor of an original work.
type Author struct {
	Name  string
	About string
	Links []SocialLink
}

// Content represents unparsed and parsed text content found in books.
type Content struct {
	Raw    string
	Parsed string
	HTML   template.HTML
}

// Series represent a set of books that are related to each other,
// such as sequels, prequels, side stories, etc.
type Series struct {
	Name   string
	Number float32
}

// SocialLink represents a link that directs the user to social
// media/contact/donation pages associated with the author.
type SocialLink struct {
	Name        string
	Address     string
	IsHyperlink bool
}
