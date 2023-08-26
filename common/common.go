package common

import (
	"html/template"
	"time"
)

// Main bookgen configuration file
type Config struct {
	Index     Index
	OutputDir string
	ThemeDir  string
	BooksDir  string
	StaticDir string
}

// The full index of
type Index struct {
	Title   string
	Author  string
	BaseURL string
	Books   []Book `toml:-`
}

// A serialized work with chapters
type Book struct {
	Config           *Config
	Title            string
	ShortDescription string
	CoverPath        string
	IndexPath        string
	LanguageCode     string
	Copyright        string
	License          string
	ChaptersDir      string
	Blurb            template.HTML `toml:-`
	Slug             string        `toml:-`
	Chapters         []Chapter     `toml:-`
}

// A chapter of a book
type Chapter struct {
	Config       *Config
	Parent       *Book
	Title        string
	Content      template.HTML
	PublishDate  time.Time `toml:",omitempty"`
	LastModified time.Time `toml:",omitempty"`
	Slug         string    `toml:-`
	Next         *Chapter  `toml:-`
	Prev         *Chapter  `toml:-`
}
