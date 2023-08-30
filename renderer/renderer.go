// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package renderer

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/JessebotX/bookgen/config"
)

func BuildSite(collection *config.Collection) error {
	resolvedOutputDir := filepath.Join(collection.Root, collection.OutputDir)
	resolvedLayoutDir := filepath.Join(collection.Root, collection.LayoutDir)

	// clean out directories
	err := os.RemoveAll(resolvedOutputDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(resolvedOutputDir, 0755)
	if err != nil {
		return err
	}

	// read book and chapter templates
	var bookTemplate *template.Template
	bookTemplatePath := filepath.Join(resolvedLayoutDir, "book.html")
	if exists(bookTemplatePath) {
		bookTemplate, err = template.ParseFiles(bookTemplatePath)
		if err != nil {
			return err
		}
	} else {
		bookTemplate = template.Must(template.New("book").Parse(BookDefaultTemplate))
	}

	rssTemplate := template.Must(template.New("rss").Parse(RSSTemplate))

	// create book index
	for _, bk := range collection.Books {
		bookOutputDir := filepath.Join(resolvedOutputDir, bk.ID)
		err = os.MkdirAll(bookOutputDir, 0755)
		if err != nil {
			return err
		}

		// static assets in book
		for _, assetPath := range bk.StaticAssets {
			resolvedPath := filepath.Join(bookOutputDir, strings.TrimPrefix(assetPath, bk.Root))
			err = os.MkdirAll(filepath.Dir(resolvedPath), 0755)
			if err != nil {
				return err
			}

			err = os.Link(assetPath, resolvedPath)
			if err != nil {
				return err
			}
		}

		// book index
		bookIndexFile, err := os.Create(filepath.Join(bookOutputDir, "index.html"))
		if err != nil {
			return err
		}

		err = bookTemplate.Execute(bookIndexFile, &bk)
		if err != nil {
			return err
		}

		// rss feed
		rssXML, err := os.Create(filepath.Join(bookOutputDir, "rss.xml"))
		if err != nil {
			return err
		}

		err = rssTemplate.Execute(rssXML, &bk)
		if err != nil {
			return err
		}
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
