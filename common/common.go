package common

import (
	"time"
)

// Main bookgen configuration file
type Config struct {
	Title     string
	Author    string
	BaseURL   string
	OutputDir string
	ThemeDir  string
	BooksDir  string
	StaticDir string
}

// A serialized work with chapters
type Book struct {
	Config           *Config
	Title            string
	ShortDescription string
	CoverPath        string
	LanguageCode     string
	Copyright        string
	License          string
	Slug             string    `toml:-`
	Chapters         []Chapter `toml-`
}

// A chapter of a book
type Chapter struct {
	Config       *Config
	Parent       Book
	Title        string
	Content      string
	PublishDate  time.Time
	LastModified time.Time
	Slug         string `toml:-`
}
