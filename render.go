package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	DirPerms = 0755
)

func RenderCollectionToWebsite(c *Collection, workingDir, outputDir string) error {
	layoutsDir := filepath.Join(workingDir, "layouts")
	collectionTemplatePath := filepath.Join(layoutsDir, "index.html")
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

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
	outIndex, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("failed to create collection index file. %w", err)
	}

	if err := collectionTemplate.ExecuteTemplate(outIndex, "index.html", c); err != nil {
		return fmt.Errorf("failed to write collection index file. %w", err)
	}

	if err := outIndex.Close(); err != nil {
		return fmt.Errorf("failed to close collection index file. %w", err)
	}

	// ---
	// Book index & chapters
	// ---

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

		if err := renderBookChapters(book.Chapters, chapterTemplate, bookOutputDir); err != nil {
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

func renderBookChapters(chapters []Chapter, chapterTemplate *template.Template, bookOutputDir string) error {
	for _, chapter := range chapters {
		fChapter, err := os.Create(filepath.Join(bookOutputDir, chapter.PageName+".html"))
		if err != nil {
			return err
		}
		defer fChapter.Close()

		if err := chapterTemplate.ExecuteTemplate(fChapter, "_chapter.html", chapter); err != nil {
			return err
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
