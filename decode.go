package mkpub

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	CollectionConfigName = "mkpub.yml"
	BookConfigName       = "mkpub-book.yml"
)

// DecodeCollection transforms source files in inputDir into a mkpub
// Collection.
func DecodeCollection(inputDir string) (Collection, error) {
	// ---
	// Init defaults
	// ---
	var collection Collection
	collection.InitDefaults()

	// ---
	// Read collection config
	// ---

	// Read yaml into a map and then turn into a struct using
	// encoding/json as intermediary. Can also possibly use mapstructure
	// package instead or a custom made mapToStruct function.

	yamlData, err := os.ReadFile(filepath.Join(inputDir, CollectionConfigName))
	if err != nil {
		return collection, fmt.Errorf("decode collection '%s': %w", inputDir, err)
	}

	if err := yaml.Unmarshal(yamlData, &collection.Params); err != nil {
		return collection, fmt.Errorf("decode collection '%s': failed to parse config: %v", inputDir, yaml.FormatError(err, false, true))
	}

	if err := mapToStruct(collection.Params, &collection); err != nil {
		return collection, fmt.Errorf("decode collection '%s': failed to parse config: %w", inputDir, err)
	}

	// ---
	// Check requirements
	// ---
	if strings.TrimSpace(collection.Title) == "" {
		return collection, fmt.Errorf("decode collection '%s': collection title is required and cannot be empty/only spaces", inputDir)
	}

	if strings.TrimSpace(collection.LanguageCode) == "" {
		return collection, fmt.Errorf("decode collection '%s': collection language code is required and cannot be empty/only spaces", inputDir)
	}

	// ---
	// Decode books
	// ---
	booksDir := filepath.Join(inputDir, "books")
	bookItems, err := os.ReadDir(booksDir)
	if err != nil {
		return collection, fmt.Errorf("decode collection '%s': failed to read directory '%s': %w", inputDir, booksDir, err)
	}

	for _, item := range bookItems {
		if !item.IsDir() {
			continue
		}

		bookDir := filepath.Join(inputDir, "books", item.Name())
		book, err := DecodeBook(bookDir, &collection)
		if err != nil {
			return collection, err
		}

		collection.Books = append(collection.Books, book)
	}

	return collection, nil
}

// DecodeBook transforms source files in inputDir into a mkpub Book.
func DecodeBook(inputDir string, collection *Collection) (Book, error) {
	id := filepath.Base(inputDir)

	var book Book
	book.InitDefaults(id, collection)

	// ---
	// Decode chapters
	// ---

	chaptersDir := filepath.Join(inputDir, "chapters")
	chapterItems, err := os.ReadDir(chaptersDir)
	if err != nil {
		return book, fmt.Errorf("decode book '%s': failed to read directory '%s': %w", inputDir, chaptersDir, err)
	}

	for _, item := range chapterItems {
		if item.IsDir() {
			continue
		}

		if !strings.HasSuffix(item.Name(), ".md") {
			continue
		}

		chapterPath := filepath.Join(chaptersDir, item.Name())
		chapter, err := decodeChapter(chapterPath, &book)
		if err != nil {
			return book, fmt.Errorf("decode book '%s': failed to decode chapter '%s': %w", inputDir, chaptersDir, err)
		}

		book.Chapters = append(book.Chapters, chapter)
	}

	return book, nil
}

// decodeChapter transforms source file at path into a mkpub Chapter.
func decodeChapter(path string, book *Book) (Chapter, error) {
	id := strings.TrimSuffix(filepath.Base(path), ".md")

	var chapter Chapter
	chapter.InitDefaults(id, book)

	return chapter, nil
}

func mapToStruct(m map[string]any, s any) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, s); err != nil {
		return err
	}

	return nil
}
