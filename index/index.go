package index

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JessebotX/bookgen/common"
	"github.com/JessebotX/bookgen/book"
)

// Generate configuration file
func GenerateHTMLSiteFromConfig(config *common.Config) error {
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

	books, err := Books(config)
	if err != nil {
		return err
	}

	log.Println(books)

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

		config := common.Book{
			Config: config,
			LanguageCode: "en",
			Copyright: "Copyright (c) " + config.Author,
			License: "All rights reserved.",
			CoverPath: filepath.Join(booksDir, dir.Name(), "cover.jpg"),
			Slug: dir.Name(),
		}
		configPath := filepath.Join(booksDir, dir.Name(), "bookgen-book.toml")
		err = book.UnmarshalBookConfig(configPath, &config)
		if err != nil {
			return nil, err
		}

		books = append(books, config)
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
