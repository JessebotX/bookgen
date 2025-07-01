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
	// Valid fields for Book.Status (case-insensitive).
	BookStatusValidValues = []string{"completed", "hiatus", "ongoing"}
)

// Author represents an individual writer or contributor of an original work.
type Author struct {
	Name  string
	About string
	Links []SocialLink
}

// Content represents unparsed and parsed text content found in books.
type Content struct {
	Raw   string
	HTML  template.HTML
	XHTML template.HTML
}

// Internal represents the app's settings that may be useful
// for themes to know about.
type Internal struct {
	GenerateEPUB bool
}

// Series represent a set of books that are related to each other,
// such as sequels, prequels, side stories, etc.
type Series struct {
	Name   string
	Number float32
}

// SocialLink represents usually represents a link that directs the
// user to another website/social media/contact/donation pages
// associated with the author and/or the written work itself.
type SocialLink struct {
	Name        string
	Address     string
	IsHyperlink bool
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

func (c *Collection) InitializeDefaults() {
	c.Title = "My Writing"
	c.ConfigFormatVersion = 0
	c.Internal.GenerateEPUB = true
}

// Close properly deallocates elements in the Collection object such
// as maps, and calls Book.Close for each Book in Collection.Books.
func (c *Collection) Close() {
	clear(c.Params)

	for _, b := range c.Books {
		b.Close()
	}
}

// CheckRequirementsForParsing checks if required fields have valid
// values for future parsing (e.g. Collection.Title is not empty). It
// does not check for things such as the existence of file
// contents/paths that the user may have specified, and it assumes
// that the Collection has been initialized with correct defaults.
func (c *Collection) CheckRequirementsForParsing() error {
	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("missing/empty required field `title`")
	}

	return nil
}

// Book represents an ordered list of chapters.
//
// NOTE: if Book.Parent exists, then Book.PageName must be unique
// within the Collection in Collection.Books.
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
	Mirrors          []SocialLink
	DatePublished    time.Time
	DateModified     time.Time
	Content          Content
	IsStub           bool
	Chapters         []Chapter
}

func (b *Book) InitializeDefaults(workingDir string, parent *Collection) {
	b.PageName = filepath.Base(workingDir)
	b.Parent = parent
	b.IsStub = false
	b.Status = "completed"
	b.Internal.GenerateEPUB = true

	if parent != nil {
		b.Internal.GenerateEPUB = parent.Internal.GenerateEPUB

		if strings.TrimSpace(parent.LanguageCode) != "" {
			b.LanguageCode = parent.LanguageCode
		}
	}
}

// CheckRequirementsForParsing checks if required fields have valid
// values for future parsing (e.g. Book.Title is not empty). It does
// not check for things such as the existence of file contents/paths
// that the user may have specified, and it assumes that the Book has
// been initialized with correct defaults (i.e. assuming Book.PageName is
// unique in a Collection).
func (b *Book) CheckRequirementsForParsing(workingDir string) error {
	// check if user accidentally set PageName
	if b.PageName != filepath.Base(workingDir) {
		return fmt.Errorf("field `PageName` must equal to the base name of the working directory.")
	}

	if strings.TrimSpace(b.Title) == "" {
		return fmt.Errorf("missing/empty required field `title`.")
	}

	if strings.TrimSpace(b.Status) != "" {
		if !slices.Contains(BookStatusValidValues, strings.ToLower(b.Status)) {
			return fmt.Errorf("invalid value for field `status`. Must be one of the following options (case-insensitive): %v.", strings.Join(BookStatusValidValues[:], " | "))
		}
	}

	return nil
}

// Close properly deallocates elements in the Book object such as
// maps, and calls Chapter.Close for each Chapter in Book.Chapters.
func (b *Book) Close() {
	clear(b.Params)

	for _, c := range b.Chapters {
		c.Close()
	}
}

// Chapter represents a division in a Book, primarily containing the
// book's text content.
//
// NOTE: Chapter.PageName must be unique within a Book in Book.Chapters
type Chapter struct {
	Params        map[string]any
	Parent        *Book
	Previous      *Chapter
	Next          *Chapter
	PageName      string
	Title         string
	Subtitle      string
	Description   string
	Order         int
	Authors       []Author
	Copyright     string
	LanguageCode  string
	DatePublished time.Time
	DateModified  time.Time
	Content       Content
}

func (c *Chapter) InitializeDefaults(workingDir string, parent *Book) {
	c.Parent = parent
	c.PageName = strings.TrimSuffix(filepath.Base(workingDir), ".md")
	c.Order = 1

	if parent != nil {
		if strings.TrimSpace(parent.LanguageCode) != "" {
			c.LanguageCode = parent.LanguageCode
		}

		if strings.TrimSpace(parent.Copyright) != "" {
			c.Copyright = parent.Copyright
		}
	}
}

// Close properly deallocates any elements in the Chapter object such
// as maps.
func (c *Chapter) Close() {
	clear(c.Params)
}
