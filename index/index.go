package index

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"html/template"
	"strings"

	"github.com/JessebotX/bookgen/book"
	"github.com/JessebotX/bookgen/common"
)

// Generate configuration file
func GenerateHTMLSiteFromConfig(config *common.Config) error {
	index := config.Index

	// recreate output directory
	err := os.RemoveAll(config.OutputDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		return err
	}

	// copy static files to output
	err = copyStaticDirToOutputDir(config.StaticDir, config.OutputDir)
	if err != nil {
		return err
	}

	// Get books
	books, err := Books(config)
	if err != nil {
		return err
	}

	// read templates
	bookTemplatePath := filepath.Join(config.ThemeDir, "book.html")
	bookTemplate, err := template.ParseFiles(bookTemplatePath)
	if err != nil {
		return err
	}

	chapterTemplatePath := filepath.Join(config.ThemeDir, "chapter.html")
	chapterTemplate, err := template.ParseFiles(chapterTemplatePath)
	if err != nil {
		return err
	}

	// create book output
	for i, _ := range books {
		bookItem := &books[i]
		// get all the chapters
		err = book.UnmarshalChapters(bookItem)
		if err != nil {
			return err
		}

		bookOutputDir := filepath.Join(config.OutputDir, bookItem.Slug)
		err = os.MkdirAll(bookOutputDir, 0755)
		if err != nil {
			return err
		}

		// copy cover image to output
		baseCoverPath := filepath.Base(bookItem.CoverPath)
		err = os.Link(bookItem.CoverPath, filepath.Join(bookOutputDir, baseCoverPath))
		if err != nil {
			log.Println("WARN", err)
		}

		// generate book index
		newBookIndexOutput, err := os.Create(filepath.Join(bookOutputDir, "index.html"))
		if err != nil {
			return err
		}

		err = bookTemplate.Execute(newBookIndexOutput, bookItem)
		if err != nil {
			return err
		}

		// generate chapters
		for _, chapter := range bookItem.Chapters {
			newChapterOutput, err := os.Create(filepath.Join(bookOutputDir, chapter.Slug + ".html"))
			if err != nil {
				return err
			}

			err = chapterTemplate.Execute(newChapterOutput, &chapter)
			if err != nil {
				return err
			}
		}
	}

	index.Books = books
	// TODO generate site index

	return nil
}

// Retrieve all books
func Books(config *common.Config) ([]common.Book, error) {
	booksDir := config.BooksDir

	dirs, err := os.ReadDir(booksDir)
	if err != nil {
		return nil, err
	}

	books := make([]common.Book, 0)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		bookItemDir := filepath.Join(booksDir, dir.Name())
		bookConfig := &common.Book{
			Config:       config,
			LanguageCode: "en",
			Copyright:    "Copyright (c) " + config.Index.Author,
			License:      "All rights reserved.",
			CoverPath:    "cover.jpg",
			IndexPath:    "index.md",
			ChaptersDir:  "chapters",
			Slug:         dir.Name(),
		}
		configPath := filepath.Join(bookItemDir, "bookgen-book.toml")
		err = book.UnmarshalBookConfig(configPath, bookConfig)
		if err != nil {
			return nil, err
		}

		err := book.UnmarshalIndexBlurb(bookConfig)
		if err != nil {
			return nil, err
		}

		books = append(books, *bookConfig)
	}

	return books, nil
}

// Copy every item in the static files directory to the output
// directory while preserving its structure relative to the static
// directory.
func copyStaticDirToOutputDir(staticDir, outputDir string) error {
	filepath.WalkDir(
		staticDir,
		func(path string, dir fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			base := strings.TrimPrefix(path, staticDir)
			outputPath := filepath.Join(outputDir, base)

			if dir.IsDir() {
				err = os.MkdirAll(outputPath, 0755)
				if err != nil {
					return err
				}

				return nil
			}

			err = os.Link(path, outputPath)
			if err != nil {
				return err
			}

			return nil
		},
	)

	return nil
}
