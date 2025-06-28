package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func RenderCollectionToWebsite(c *Collection, workingDir, outputDir string) error {
	layoutsDir := filepath.Join(workingDir, "layouts")
	//baseTemplatePath := filepath.Join(layoutsDir, "_base.html")
	collectionTemplatePath := filepath.Join(layoutsDir, "index.html")
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

	if err := os.MkdirAll(outputDir, os.ModeDir); err != nil {
		return fmt.Errorf("failed to create output directory. %w", err)
	}

	// ---
	// Read templates
	// ---
	templateFileNames := []string{collectionTemplatePath}
	fileNames, err := filepath.Glob(filepath.Join(layoutsDir, "_*_t.html"))
	if err == nil {
		templateFileNames = append(templateFileNames, fileNames...)
	}

	// currentParsing := []string{collectionTemplatePath}

	// if _, err := os.Stat(baseTemplatePath); err == nil {
	// 	currentParsing = append(currentParsing, baseTemplatePath)
	// } else if err != nil && !errors.Is(err, os.ErrNotExist) {
	// 	return fmt.Errorf("failed to stat base template `%v`. %w", baseTemplatePath, err)
	// }

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

	for _, book := range c.Books {
		// bookWorkingDir := filepath.Join(workingDir, "books", book.PageName)
		bookOutputDir := filepath.Join(outputDir, book.PageName)
		if err := os.MkdirAll(bookOutputDir, os.ModeDir); err != nil {
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
