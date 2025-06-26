package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func RenderCollectionToWebsite(c *Collection, workingDir, outputDir string) error {
	layoutsDir := filepath.Join(workingDir, "layouts")
	collectionTemplatePath := filepath.Join(layoutsDir, "index.html")
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

	if err := os.MkdirAll(outputDir, os.ModeDir); err != nil {
		return fmt.Errorf("failed to create output directory. %w", err)
	}

	// ---
	// Read templates
	// ---
	t, err := template.ParseFiles(
		collectionTemplatePath,
		bookTemplatePath,
		chapterTemplatePath,
	)
	if err != nil {
		return fmt.Errorf("failed to parse template. %w", err)
	}

	// ---
	// Collection index
	// ---
	outIndex, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("failed to create collection index file. %w", err)
	}
	defer outIndex.Close()

	if err := t.ExecuteTemplate(outIndex, "index.html", c); err != nil {
		return fmt.Errorf("failed to write collection index file. %w", err)
	}

	return nil
}
