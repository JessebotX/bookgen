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
)

type Bookgen struct {
	BooksDir   string
	ThemeDir   string
	OutputDir  string
	StaticDir  string
	Collection *Collection
}

type Collection struct {
	Title        string `toml:"title"`
	BaseURL      string `toml:"baseURL,omitempty"`
	LanguageCode string `toml:"languageCode,omitempty"`
	Books        []Book `toml:-`
}

type Book struct {
	Title            string      `toml:"title"`
	Author           Author      `toml:"author,omitempty"`
	Mirrors          []string    `toml:"mirrors,omitempty"`
	ShortDescription string      `toml:"shortDescription,omitempty"`
	Genre            string      `toml:"genre,omitempty"`
	Status           string      `toml:"status,omitempty"`
	CoverPath        string      `toml:"coverPath,omitempty"`
	LanguageCode     string      `toml:"languageCode,omitempty"`
	Copyright        string      `toml:"copyright,omitempty"`
	License          string      `toml:"license,omitempty"`
	Blurb            string      `toml-`
	Chapters         []Chapter   `toml:-`
	Collection       *Collection `toml:-`
	ChaptersDir      string      `toml:-`
}

type Author struct {
	Name     string `toml:"name"`
	Bio      string `toml:"bio,omitempty"`
	Donation []Donation
}

type Donation struct {
	Name        string `toml:"name,omitempty"`
	Link        string `toml:"link"`
	NonLinkItem bool   `toml:"nonLinkItem,omitempty"`
}

type Chapter struct {
	Title        string
	Description  string
	Published    time.Time
	LastModified time.Time
	Content      string
	Parent       *Book
	Collection   *Collection
	ID           string
}

func (c Chapter) ContentHTML() template.HTML {
	return template.HTML(c.Content)
}
