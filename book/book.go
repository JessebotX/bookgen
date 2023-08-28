package book

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/JessebotX/bookgen/common"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type MarkdownHTMLResult struct {
	HTML template.HTML
	Meta map[string]interface{}
}

// Unmarshal a single-book configuration
func Unmarshal(path string, config *common.Book) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(source, config)
	if err != nil {
		return err
	}

	// resolve paths
	if config.IndexPath != "" {
		config.IndexPath = filepath.Join(filepath.Dir(path), config.IndexPath)
	}

	if config.ChaptersDir != "" {
		config.ChaptersDir = filepath.Join(filepath.Dir(path), config.ChaptersDir)
	}

	return nil
}

func UnmarshalBlurb(book *common.Book) error {
	converter := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
	)
	result, err := markdownFileToHTML(book.IndexPath, converter)
	if err != nil {
		return err
	}

	book.Blurb = result.HTML

	return nil
}

func UnmarshalChapters(book *common.Book) error {
	chaptersDir := book.ChaptersDir
	converter := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
	)

	items, err := os.ReadDir(chaptersDir)
	if err != nil {
		return err
	}

	chapters := make([]common.Chapter, 0)
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		chapter := common.Chapter{
			Title:  "[CHAPTER_TITLE]",
			Config: book.Config,
			Parent: book,
			Slug:   strings.TrimSuffix(item.Name(), ".md"),
		}

		output, err := markdownFileToHTML(filepath.Join(chaptersDir, item.Name()), converter)
		if err != nil {
			return err
		}
		chapter.Content = output.HTML

		// retrieve metadata
		if output.Meta["title"] != nil {
			chapter.Title = output.Meta["title"].(string)
		}

		timeLayout := "2006-01-02T15:04:05-07:00"
		if output.Meta["date"] != nil {
			t, err := time.Parse(timeLayout, output.Meta["date"].(string))
			if err != nil {
				return err
			}
			chapter.PublishDate = t
		}

		if output.Meta["lastmod"] != nil {
			t, err := time.Parse(timeLayout, output.Meta["lastmod"].(string))
			if err != nil {
				return err
			}
			chapter.LastModified = t
		}

		chapters = append(chapters, chapter)
	}

	// set next and prev chapters
	totalChapters := len(chapters)
	for i := 0; i < totalChapters; i++ {
		if (i - 1) >= 0 {
			chapters[i].Prev = &chapters[i-1]
		}

		if (i + 1) <= (totalChapters - 1) {
			chapters[i].Next = &chapters[i+1]
		}
	}

	book.Chapters = chapters

	return nil
}

func markdownFileToHTML(filepath string, converter goldmark.Markdown) (*MarkdownHTMLResult, error) {
	source, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	context := parser.NewContext()
	err = converter.Convert(source, &buffer, parser.WithContext(context))
	if err != nil {
		return nil, err
	}
	metadata := meta.Get(context)

	return &MarkdownHTMLResult{
		HTML: template.HTML(buffer.String()),
		Meta: metadata,
	}, nil
}
