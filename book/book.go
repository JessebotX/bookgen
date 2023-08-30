// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

// Package book manages a book object found in collections.
package book

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/JessebotX/bookgen/chapter"
	"github.com/JessebotX/bookgen/config"
	"github.com/JessebotX/bookgen/renderer"
	"github.com/yuin/goldmark"
)

func Create(directory string, collection *config.Collection, converter goldmark.Markdown) (*config.Book, error) {
	book := config.Book{
		Root:        directory,
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
	book.Blurb, err = getBlurb(filepath.Join(directory, "index.md"), converter)
	if err != nil {
		return nil, err
	}

	err = unmarshalChapters(&book, converter)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func unmarshalChapters(book *config.Book, converter goldmark.Markdown) error {
	resolvedChaptersDir := filepath.Join(book.Root, book.ChaptersDir)

	chapters := make([]config.Chapter, 0)
	err := filepath.WalkDir(
		resolvedChaptersDir,
		func(path string, dir fs.DirEntry, err error) error {
			if dir.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".md" {
				book.StaticAssets = append(book.StaticAssets, path)
				return nil
			}

			chapter, err := chapter.Create(path, book, converter)
			if err != nil {
				return err
			}
			chapters = append(chapters, *chapter)

			return nil
		},
	)
	if err != nil {
		return err
	}

	chapter.UnmarshalNextPrev(chapters)
	book.Chapters = chapters

	return nil
}

func getBlurb(path string, converter goldmark.Markdown) (template.HTML, error) {
	html, _, err := renderer.MarkdownFileToHTML(path, converter)
	if err != nil {
		return template.HTML(""), err
	}

	return html, nil
}
