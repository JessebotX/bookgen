// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

/*
Package config provides the structure of toml and yaml metadata found
in collections (main site index), books (book index) and chapters.
*/
package config

import (
	"html/template"
	"time"

	readingtime "github.com/begmaroman/reading-time"
)

type Collection struct {
	Root         string `toml:-`
	BooksDir     string `toml:"booksDir"`
	LayoutDir    string `toml:"layoutDir"`
	OutputDir    string `toml:"outputDir"`
	StaticDir    string `toml:"staticDir"`
	Title        string `toml:"title"`
	BaseURL      string `toml:"baseURL,omitempty"`
	LanguageCode string `toml:"languageCode,omitempty"`
	Books        []Book `toml:-`
}

type Book struct {
	Root             string        `toml:-`
	ID               string        `toml:"id,omitempty"`
	Title            string        `toml:"title"`
	Author           Author        `toml:"author,omitempty"`
	Mirrors          []string      `toml:"mirrors,omitempty"`
	ShortDescription string        `toml:"shortDescription,omitempty"`
	Genre            string        `toml:"genre,omitempty"`
	Status           string        `toml:"status,omitempty"`
	CoverPath        string        `toml:"coverPath,omitempty"`
	LanguageCode     string        `toml:"languageCode,omitempty"`
	Copyright        string        `toml:"copyright,omitempty"`
	License          string        `toml:"license,omitempty"`
	ChaptersDir      string        `toml:"chaptersDir"`
	Blurb            template.HTML `toml:-`
	Chapters         []Chapter     `toml:-`
	StaticAssets     []string      `toml:-`
	Collection       *Collection   `toml:-`
}

type Author struct {
	Name   string `toml:"name"`
	Bio    string `toml:"bio,omitempty"`
	Donate []Donation
}

type Donation struct {
	Name        string `toml:"name,omitempty"`
	Link        string `toml:"link"`
	NonLinkItem bool   `toml:"nonLinkItem,omitempty"`
}

type Chapter struct {
	ID           string
	Title        string
	Description  string
	Published    time.Time
	LastModified time.Time
	Content      template.HTML
	Parent       *Book
	Collection   *Collection
	Next         *Chapter
	Prev         *Chapter
}

func (c Chapter) SlugHTML() string {
	return c.ID + ".html"
}

func (c Chapter) EstimatedReadingTime() *readingtime.Result {
	estimate := readingtime.Estimate(string(c.Content))
	return estimate
}
