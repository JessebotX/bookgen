// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package chapter

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/JessebotX/bookgen/config"
	"github.com/JessebotX/bookgen/renderer"
	"github.com/yuin/goldmark"
)

const DateTimeLayout = "2006-01-02T15:04:05-07:00"

func Create(path string, book *config.Book, converter goldmark.Markdown) (*config.Chapter, error) {
	chapter := config.Chapter{
		ID:         strings.TrimSuffix(filepath.Base(path), ".md"),
		Parent:     book,
		Collection: book.Collection,
	}

	content, metadata, err := renderer.MarkdownFileToHTML(path, converter)
	if err != nil {
		return nil, err
	}

	chapter.Content = content

	if metadata["title"] != nil {
		chapter.Title = metadata["title"].(string)
	} else {
		chapter.Title = chapter.ID // title = filename
	}

	if metadata["description"] != nil {
		chapter.Description = metadata["description"].(string)
	}

	if metadata["date"] != nil {
		d := metadata["date"].(string)
		t, err := time.Parse(DateTimeLayout, d)
		if err != nil {
			return nil, err
		}

		chapter.Published = t
	}

	if metadata["lastmod"] != nil {
		d := metadata["lastmod"].(string)
		t, err := time.Parse(DateTimeLayout, d)
		if err != nil {
			return nil, err
		}

		chapter.LastModified = t
	}

	return &chapter, nil
}

// Set next and previous chapters for each chapter object
func UnmarshalNextPrev(chapters []config.Chapter) {
	for i, _ := range chapters {
		if (i - 1) >= 0 {
			chapters[i].Prev = &chapters[i-1]
		}

		if (i + 1) <= (len(chapters) - 1) {
			chapters[i].Next = &chapters[i+1]
		}
	}
}

