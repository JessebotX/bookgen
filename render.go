package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func RenderCollectionToWebsite(c *Collection, workingDir, outputDir string) error {
	layoutsDir := filepath.Join(workingDir, "layouts")
	baseTemplatePath := filepath.Join(layoutsDir, "_base.html")
	collectionTemplatePath := filepath.Join(layoutsDir, "index.html")
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

	if err := os.MkdirAll(outputDir, os.ModeDir); err != nil {
		return fmt.Errorf("failed to create output directory. %w", err)
	}

	// ---
	// Read templates
	// ---
	currentParsing := []string{collectionTemplatePath}

	if _, err := os.Stat(baseTemplatePath); err == nil {
		currentParsing = append(currentParsing, baseTemplatePath)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat base template `%v`. %w", baseTemplatePath, err)
	}

	collectionTemplate, err := template.ParseFiles(currentParsing...)
	if err != nil {
		return fmt.Errorf("failed to parse collection template. %w", err)
	}

	currentParsing[0] = bookTemplatePath
	bookTemplate, err := template.ParseFiles(currentParsing...)
	if err != nil {
		return fmt.Errorf("failed to parse book template. %w", err)
	}

	currentParsing[0] = chapterTemplatePath
	chapterTemplate, err := template.ParseFiles(currentParsing...)
	if err != nil {
		return fmt.Errorf("failed to parse chapter template. %w", err)
	}

	_, _ = bookTemplate, chapterTemplate

	// ---
	// Collection index
	// ---
	outIndex, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("failed to create collection index file. %w", err)
	}
	defer outIndex.Close()

	if err := collectionTemplate.ExecuteTemplate(outIndex, "index.html", c); err != nil {
		return fmt.Errorf("failed to write collection index file. %w", err)
	}

	return nil
}
