package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

const (
	DirPerms = 0755
)

func RenderCollectionToWebsite(c *Collection, workingDir, outputDir string, enableMinify bool) error {
	layoutsDir := filepath.Join(workingDir, "layouts")
	collectionTemplatePath := filepath.Join(layoutsDir, "index.html")
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

	minifier := minify.New()
	minifier.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepSpecialComments: true,
		KeepQuotes:          true,
	})
	minifier.AddFunc("text/css", css.Minify)

	if err := os.MkdirAll(outputDir, DirPerms); err != nil {
		return fmt.Errorf("failed to create output directory. %w", err)
	}

	// ---
	// Copy global static items into output
	// ---
	if err := copyStaticFilesToDir(layoutsDir, outputDir, layoutsDir, []string{
		"index.html",
		"_book.html",
		"_chapter.html",
	}, []string{
		"_*_t.html",
	}); err != nil {
		return fmt.Errorf("failed to copy files to output. %w", err)
	}

	// ---
	// Read templates
	// ---
	templateFileNames := []string{collectionTemplatePath}
	fileNames, err := filepath.Glob(filepath.Join(layoutsDir, "_*_t.html"))
	if err == nil {
		templateFileNames = append(templateFileNames, fileNames...)
	}

	templateFileNames[0] = collectionTemplatePath
	collectionTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("failed to parse collection template. %w", err)
	}

	templateFileNames[0] = bookTemplatePath
	bookTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("failed to parse book template. %w", err)
	}

	templateFileNames[0] = chapterTemplatePath
	chapterTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("failed to parse chapter template. %w", err)
	}

	// ---
	// Collection index
	// ---
	outputIndexPath := filepath.Join(outputDir, "index.html")

	outIndex, err := os.Create(outputIndexPath)
	if err != nil {
		return fmt.Errorf("failed to create collection index file. %w", err)
	}

	if err := collectionTemplate.ExecuteTemplate(outIndex, "index.html", c); err != nil {
		return fmt.Errorf("failed to write collection index file. %w", err)
	}

	if enableMinify {
		if err := minifyFileHTML(outputIndexPath, outIndex, minifier); err != nil {
			return fmt.Errorf("failed to minify collection index output file. %w", err)
		}
	}

	if err := outIndex.Close(); err != nil {
		return fmt.Errorf("failed to close collection index file. %w", err)
	}

	// TODO: epub generation
	for _, book := range c.Books {
		bookWorkingDir := filepath.Join(workingDir, "books", book.PageName)
		bookOutputDir := filepath.Join(outputDir, book.PageName)
		if err := os.MkdirAll(bookOutputDir, DirPerms); err != nil {
			return fmt.Errorf("failed to create book `%v` directory. %w", book.PageName, err)
		}

		outBook, err := os.Create(filepath.Join(bookOutputDir, "index.html"))
		if err != nil {
			return fmt.Errorf("failed to create book `%v` index file. %w", book.PageName, err)
		}
		defer outBook.Close()

		if err := bookTemplate.ExecuteTemplate(outBook, "_book.html", book); err != nil {
			return fmt.Errorf("failed to write book `%v` index file. %w", book.PageName, err)
		}

		if enableMinify {
			if err := minifyFileHTML(bookTemplatePath, outBook, minifier); err != nil {
				return fmt.Errorf("failed to minify book `%v` index output file. %w", book.PageName, err)
			}
		}

		if err := renderBookChapters(book.Chapters, chapterTemplate, chapterTemplatePath, bookOutputDir, minifier, enableMinify); err != nil {
			return fmt.Errorf("failed to write book `%v` chapter file. %w", book.PageName, err)
		}

		// Add cover image to output
		if strings.TrimSpace(book.CoverImageName) != "" {
			coverPathOld := filepath.Join(bookWorkingDir, book.CoverImageName)
			coverPathNew := filepath.Join(bookOutputDir, book.CoverImageName)

			if err := os.RemoveAll(coverPathNew); err != nil {
				return fmt.Errorf("failed to remove cover path in output directory `%v`. %w", err)
			}

			if err := os.Link(coverPathOld, coverPathNew); err != nil {
				return err
			}
		}
	}

	return nil
}

func renderBookChapters(chapters []Chapter, chapterTemplate *template.Template, chapterTemplatePath, bookOutputDir string, minifier *minify.M, enableMinify bool) error {
	for _, chapter := range chapters {
		fChapter, err := os.Create(filepath.Join(bookOutputDir, chapter.PageName+".html"))
		if err != nil {
			return err
		}
		defer fChapter.Close()

		if err := chapterTemplate.ExecuteTemplate(fChapter, "_chapter.html", chapter); err != nil {
			return err
		}

		if enableMinify {
			if err := minifyFileHTML(chapterTemplatePath, fChapter, minifier); err != nil {
				return fmt.Errorf("failed to minify book `%v` chapter `%v` index output file. %w", chapter.Parent.PageName, chapter.PageName, err)
			}
		}
	}

	return nil
}

func copyStaticFilesToDir(currDir, newDir, rootDir string, relExcludes, relExcludesPatterns []string) error {
	items, err := os.ReadDir(currDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		oldPath := filepath.Join(currDir, item.Name())
		oldPathFromRoot := strings.TrimLeft(strings.TrimPrefix(oldPath, rootDir), "/\\")
		newPath := filepath.Join(newDir, oldPathFromRoot)

		// Check against exclusions
		if slices.Contains(relExcludes, oldPathFromRoot) {
			continue
		}

		matching := false
		for _, pattern := range relExcludesPatterns {
			matching, err = filepath.Match(pattern, oldPathFromRoot)
			if err != nil {
				return err
			}
		}
		if matching {
			continue
		}

		if item.IsDir() {
			if err := os.MkdirAll(newPath, DirPerms); err != nil {
				return err
			}

			return copyStaticFilesToDir(oldPath, newDir, rootDir, relExcludes, relExcludesPatterns)
		}

		// Copy item from old to new path (re-copy if already exists
		if err := os.RemoveAll(newPath); err != nil {
			return err
		}

		if err := os.Link(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

func minifyFileHTML(path string, f *os.File, minifier *minify.M) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	mb, err := minifier.Bytes("text/html", b)
	if err != nil {
		return err
	}

	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	if _, err := f.Write(mb); err != nil {
		return err
	}

	return nil
}
