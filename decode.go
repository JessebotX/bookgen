package mkpub

import (
	"bufio"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"golang.org/x/sync/errgroup"
)

const (
	CollectionConfigName = "mkpub.yml"
	BookConfigName       = "mkpub-book.yml"
)

// DecodeCollection transforms source files in inputDir into a mkpub
// [Collection].
func DecodeCollection(inputDir string) (Collection, error) {
	var collection Collection
	collection.InitDefaults(inputDir)

	// --- Read collection config ---
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

	// --- Check requirements ---
	if strings.TrimSpace(collection.Title) == "" {
		return collection, fmt.Errorf("decode collection '%s': collection title is required and cannot be empty/only spaces", inputDir)
	}

	if strings.TrimSpace(collection.LanguageCode) == "" {
		return collection, fmt.Errorf("decode collection '%s': collection language code is required and cannot be empty/only spaces", inputDir)
	}

	// --- Further parsing ---
	if collection.FaviconImageName != "" {
		if _, err := os.Stat(filepath.Join(inputDir, collection.FaviconImageName)); errors.Is(err, os.ErrNotExist) {
			return collection, fmt.Errorf("decode collection '%s': failed to find favicon image '%s' in input directory '%s')", inputDir, filepath.Clean(collection.FaviconImageName), inputDir)
		} else if err != nil {
			return collection, fmt.Errorf("decode collection '%s': %w", inputDir, err)
		}
	}

	// --- Decode books ---
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

// DecodeBook transforms source files in inputDir into a mkpub [Book].
func DecodeBook(inputDir string, collection *Collection) (Book, error) {
	id := filepath.Base(inputDir)

	var book Book
	book.InitDefaults(id, inputDir, collection)

	// --- Parse config ---

	yamlData, err := os.ReadFile(filepath.Join(inputDir, BookConfigName))
	if err != nil {
		return book, fmt.Errorf("decode book '%s': failed to read config: %w", inputDir, err)
	}

	if err := yaml.Unmarshal(yamlData, &book.Params); err != nil {
		return book, fmt.Errorf("decode book '%s': failed to parse config: %v", inputDir, yaml.FormatError(err, false, true))
	}

	if err := mapToStruct(book.Params, &book); err != nil {
		return book, fmt.Errorf("decode book '%s': failed to parse config: %w", inputDir, err)
	}

	// --- Check requirements ---

	if strings.TrimSpace(book.UniqueID) == "" {
		return book, fmt.Errorf("decode book '%s': unique ID is required and cannot be empty/only spaces", inputDir)
	}

	if strings.TrimSpace(book.Title) == "" {
		return book, fmt.Errorf("decode book '%s': title is required and cannot be empty/only spaces", inputDir)
	}

	if strings.TrimSpace(book.LanguageCode) == "" {
		return book, fmt.Errorf("decode book '%s': language code is required and cannot be empty/only spaces", inputDir)
	}

	// --- Further parsing ---

	book.Status = strings.ToLower(book.Status)

	if book.CoverImageName != "" {
		if _, err := os.Stat(filepath.Join(inputDir, book.CoverImageName)); errors.Is(err, os.ErrNotExist) {
			return book, fmt.Errorf("decode book '%s': failed to find cover image '%s' in input directory '%s')", inputDir, filepath.Clean(book.CoverImageName), inputDir)
		} else if err != nil {
			return book, fmt.Errorf("decode book '%s': %w", inputDir, err)
		}
	}

	if !slices.Contains(BookStatusValues, strings.ToLower(book.Status)) {
		valid := strings.Join(BookStatusValues, ", ")

		return book, fmt.Errorf("decode book '%s': unrecognized status value '%s'.\nValue must be one of the following: (comma-separated): %s", inputDir, book.Status, valid)
	}

	rawContentPath := filepath.Join(inputDir, "index.md")
	rawContent, err := os.ReadFile(rawContentPath)
	if err != nil {
		return book, fmt.Errorf("decode book '%s': failed to read content file: %w", inputDir, err)
	}
	book.Content.Raw = rawContent

	// --- Decode chapters ---

	chaptersDir := filepath.Join(inputDir, "chapters")
	chapterItems, err := os.ReadDir(chaptersDir)
	if err != nil {
		return book, fmt.Errorf("decode book '%s': failed to read directory '%s': %w", inputDir, chaptersDir, err)
	}

	g := new(errgroup.Group)
	for _, item := range chapterItems {
		if item.IsDir() {
			continue
		}

		if !strings.HasSuffix(item.Name(), ".md") {
			continue
		}

		chapterPath := filepath.Join(chaptersDir, item.Name())
		g.Go(func() error {
			chapter, err := decodeChapter(chapterPath, &book)
			if err != nil {
				return err
			}
			book.Chapters = append(book.Chapters, chapter)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return book, fmt.Errorf("decode book '%s': failed to decode chapter '%s': %w", inputDir, chaptersDir, err)
	}

	// --- Sort decoded chapters ---
	//
	// due to concurrency, unsorted slice may be random

	slices.SortFunc(book.Chapters, func(a, b Chapter) int {
		if n := cmp.Compare(a.Order, b.Order); n != 0 {
			return n
		}

		if !a.DatePublished.Equal(b.DatePublished) {
			if a.DatePublished.After(b.DatePublished) {
				return 1
			} else {
				return -1
			}
		}

		if n := strings.Compare(a.Title, b.Title); n != 0 {
			return n
		}

		return strings.Compare(a.UniqueID, b.UniqueID)
	})

	// Set next and previous chapters
	for i := 0; i < len(book.Chapters); i++ {
		if i-1 >= 0 {
			book.Chapters[i].Previous = &book.Chapters[i-1]
		}

		if i+1 < len(book.Chapters) {
			book.Chapters[i].Next = &book.Chapters[i+1]
		}
	}

	// By default, inherit date and times from chapters if not provided
	// at the book level.
	var youngest, latest time.Time
	for _, chapter := range book.Chapters {
		if chapter.DatePublished.IsZero() {
			continue
		}

		if youngest.IsZero() {
			youngest = chapter.DatePublished
		}

		if latest.IsZero() {
			latest = chapter.DatePublished
		}

		if youngest.After(chapter.DatePublished) {
			youngest = chapter.DatePublished
		}

		if latest.Before(chapter.DatePublished) {
			latest = chapter.DatePublished
		}

		// youngest date also considers the earliest modification
		if chapter.DateModified.IsZero() {
			continue
		}

		if youngest.After(chapter.DateModified) {
			youngest = chapter.DateModified
		}
	}

	if book.DatePublishedStart.IsZero() {
		v, exists := book.Params["published_start"]
		if exists {
			book.DatePublishedStart, err = parseTimeParam(v)
			if err != nil {
				return book, err
			}
		}

		if !exists {
			book.DatePublishedStart = youngest
		}
	}

	if book.Status == BookStatusCompleted || book.Status == BookStatusInactive {
		if book.DatePublishedEnd.IsZero() {
			v, exists := book.Params["published_end"]
			if !exists {
				v, exists = book.Params["date"]
			}

			if exists {
				book.DatePublishedEnd, err = parseTimeParam(v)
				if err != nil {
					return book, err
				}
			}

			if !exists {
				book.DatePublishedEnd = latest
			}
		}

		if book.DatePublishedEnd.IsZero() {
			book.DatePublishedEnd = book.DatePublishedStart
		}

		// Make start == end if one or the other is not set
		if book.DatePublishedStart.IsZero() {
			book.DatePublishedStart = book.DatePublishedEnd
		}
	}

	return book, nil
}

// decodeChapter transforms source file at path into a mkpub [Chapter].
func decodeChapter(path string, book *Book) (Chapter, error) {
	id := strings.TrimSuffix(filepath.Base(path), ".md")

	var chapter Chapter
	chapter.InitDefaults(id, book)

	f, err := os.Open(path)
	if err != nil {
		return chapter, err
	}

	// ---
	//
	// Parse frontmatter at the top
	//
	// ---
	scanner := bufio.NewScanner(f)
	var yamlData string

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return chapter, err
	}
	first := scanner.Text()

	if strings.TrimRight(first, " ") == "---" {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimRight(line, " ") == "---" {
				break
			}

			yamlData += line + "\n"
		}
		if err := scanner.Err(); err != nil {
			return chapter, err
		}

		if err := yaml.Unmarshal([]byte(yamlData), &chapter.Params); err != nil {
			return chapter, err
		}

		if err := mapToStruct(chapter.Params, &chapter); err != nil {
			return chapter, err
		}
	}

	if err := f.Close(); err != nil {
		return chapter, err
	}

	// ---
	//
	// Further parsing
	//
	// ---
	rawContent, err := os.ReadFile(path)
	if err != nil {
		return chapter, err
	}
	chapter.Content.Raw = rawContent

	if chapter.DatePublished.IsZero() {
		v, exists := chapter.Params["published"]
		if !exists {
			v, exists = chapter.Params["date"]
		}

		if exists {
			chapter.DatePublished, err = parseTimeParam(v)
			if err != nil {
				return chapter, err
			}
		}
	}

	if chapter.DateModified.IsZero() {
		v, exists := chapter.Params["modified"]
		if !exists {
			v, exists = chapter.Params["lastmod"]
		}

		if exists {
			chapter.DatePublished, err = parseTimeParam(v)
			if err != nil {
				return chapter, err
			}
		}
	}

	return chapter, nil
}

func mapToStruct(m map[string]any, s any) error {
	// Read yaml into a map and then turn into a struct using
	// encoding/json as intermediary. Can also possibly use mapstructure
	// package instead or a custom made mapToStruct function.

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, s); err != nil {
		return err
	}

	return nil
}

func parseTimeParam(param any) (time.Time, error) {
	switch v := param.(type) {
	case time.Time:
		return v, nil
	case string:
		date, err := guessTimeFormat(v)
		if err != nil {
			return time.Time{}, err
		}
		return date, nil
	}

	return time.Time{}, fmt.Errorf("unrecognized parameter type given for date format. Value type must be either 'string' or 'time.Time'")
}

// Convert a date string input into time. String argument must be in
// the correct format or it will return an error.
func guessTimeFormat(s string) (time.Time, error) {
	formats := []string{
		"2006",
		"2006-01",
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04Z07:00",
		"2006-01-02T15:04Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04:05Z07:00",
	}

	var errs error
	for _, format := range formats {
		date, err := time.Parse(format, s)
		if err == nil {
			return date, nil
		}
		errs = errors.Join(errs, err)
	}

	return time.Time{}, fmt.Errorf("unrecognized format for '%s':\n%w", s, errs)
}
