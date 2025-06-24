package main

import (
	"html/template"
	"time"
)

// InternalSettings represents the app's settings that may be useful
// for themes to know about.
type InternalSettings struct {
	GenerateEPUB bool
}

// Collection represents a list/index of one or more books.
type Collection struct {
	Params              map[string]any
	Internal            InternalSettings
	Title               string
	Description         string
	BaseURL             string
	LanguageCode        string
	Books               []Book
	FaviconImageName    string
	ConfigFormatVersion string
}

// Close properly deallocates elements in the Collection object such
// as maps. Also closes all Books as well.
func (c Collection) Close() {
	clear(c.Params)

	for _, b := range c.Books {
		b.Close()
	}
}

// Book represents an ordered list of chapters.
type Book struct {
	Params           map[string]any
	Parent           *Collection
	Internal         InternalSettings
	UniqueID         string
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
	Chapters         []Chapter
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
	UniqueID         string
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
