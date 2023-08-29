// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

// Package book manages a book object found in collections.
package book

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/JessebotX/bookgen/config"
	"github.com/JessebotX/bookgen/renderer"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func Create(directory string, collection *config.Collection) (*config.Book, error) {
	book := config.Book{
		ID:          filepath.Base(directory),
		Collection:  collection,
		ChaptersDir: "./chapters",
	}

	configPath := filepath.Join(directory, "bookgen-book.toml")
	source, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(source, &book)
	if err != nil {
		return nil, err
	}

	// get blurb
	converter := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(),
	)

	book.Blurb, err = getBlurb(filepath.Join(directory, "index.md"), converter)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func getBlurb(path string, converter goldmark.Markdown) (template.HTML, error) {
	html, _, err := renderer.MarkdownFileToHTML(path, converter)
	if err != nil {
		return template.HTML(""), err
	}

	return html, nil
}
