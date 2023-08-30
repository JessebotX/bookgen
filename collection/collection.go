// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package collection

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/JessebotX/bookgen/book"
	"github.com/JessebotX/bookgen/config"
	"github.com/yuin/goldmark"
)

// Constructs a new collection object from a project root directory
func Create(root string, converter goldmark.Markdown) (*config.Collection, error) {
	config := config.Collection{
		Root:      root,
		BooksDir:  "./src",
		LayoutDir: "./layout",
		OutputDir: "./out",
		StaticDir: "./static",
	}

	configPath := filepath.Join(root, "bookgen.toml")
	source, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(source, &config)
	if err != nil {
		return nil, err
	}

	err = unmarshalBooks(&config, converter)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func unmarshalBooks(collection *config.Collection, converter goldmark.Markdown) error {
	booksDir := filepath.Join(collection.Root, collection.BooksDir)

	items, err := os.ReadDir(booksDir)
	if err != nil {
		return err
	}

	books := make([]config.Book, 0)
	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		book, err := book.Create(filepath.Join(booksDir, item.Name()), collection, converter)
		if err != nil {
			return err
		}

		books = append(books, *book)
	}

	collection.Books = books

	return nil
}
