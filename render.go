package mkpub

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func WriteCollectionToHTML(collection *Collection, outputDir, layoutsDir string) error {
	_ = collection
	_ = layoutsDir

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// ---
	//
	// Parse templates
	//
	// ---
	collectionTemplatePath := filepath.Join(layoutsDir, "index.html")
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

	templateFileNames := []string{collectionTemplatePath}
	fileNames, err := filepath.Glob(filepath.Join(layoutsDir, "_template_*.html"))
	if err == nil {
		templateFileNames = append(templateFileNames, fileNames...)
	}

	templateFileNames[0] = collectionTemplatePath
	collectionTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("write collection: failed to parse collection template: %w", err)
	}

	templateFileNames[0] = bookTemplatePath
	bookTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("write collection: failed to parse book template: %w", err)
	}

	templateFileNames[0] = chapterTemplatePath
	chapterTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("write collection: failed to parse chapter template: %w", err)
	}

	// ---
	//
	// Collection index
	//
	// ---
	fCollection, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("write collection: failed to create collection html output: %w", err)
	}
	defer fCollection.Close()

	if err := collectionTemplate.ExecuteTemplate(fCollection, "index.html", collection); err != nil {
		return fmt.Errorf("write collection: failed to write collection index file: %w", err)
	}

	g := new(errgroup.Group)
	for _, book := range collection.Books {
		bookOutputDir := filepath.Join(outputDir, "books")

		g.Go(func() error {
			if err := writeBookToHTML(&book, bookOutputDir, "", bookTemplate, chapterTemplate); err != nil {
				return err
			}

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("write collection: %w", err)
	}

	return nil
}

func WriteBookToHTML(book *Book, outputDir, layoutsDir string) error {
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

	templateFileNames := []string{bookTemplatePath}
	fileNames, err := filepath.Glob(filepath.Join(layoutsDir, "_template_*.html"))
	if err == nil {
		templateFileNames = append(templateFileNames, fileNames...)
	}

	templateFileNames[0] = bookTemplatePath
	bookTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("write book '%s': failed to parse book template: %w", book.UniqueID, err)
	}

	templateFileNames[0] = chapterTemplatePath
	chapterTemplate, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return fmt.Errorf("write book '%s': failed to parse chapter template: %w", book.UniqueID, err)
	}

	return writeBookToHTML(book, outputDir, layoutsDir, bookTemplate, chapterTemplate)
}

func writeBookToHTML(book *Book, outputDir, layoutsDir string, bookTemplate *template.Template, chapterTemplate *template.Template) error {
	_ = book
	_ = outputDir
	_ = layoutsDir
	_ = bookTemplate
	_ = chapterTemplate

	return nil
}
