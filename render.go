package mkpub

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"golang.org/x/sync/errgroup"
)

var (
	goldmarkExtensions = goldmark.WithExtensions(
		extension.GFM,
		extension.Footnote,
		extension.Typographer,
	)
	md = goldmark.New(
		goldmarkExtensions,
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
	)
)

func WriteCollectionToHTML(collection *Collection, outputDir, layoutsDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// --- Parse templates ---
	otherTemplatesPath := filepath.Join(layoutsDir, "_template_*.html")
	collectionTemplate, err := parseTemplate(filepath.Join(layoutsDir, "index.html"), otherTemplatesPath)
	if err != nil {
		return fmt.Errorf("write collection: failed to parse collection template: %w", err)
	}

	bookTemplate, err := parseTemplate(filepath.Join(layoutsDir, "_book.html"), otherTemplatesPath)
	if err != nil {
		return fmt.Errorf("write collection: failed to parse book template: %w", err)
	}

	chapterTemplate, err := parseTemplate(filepath.Join(layoutsDir, "_chapter.html"), otherTemplatesPath)
	if err != nil {
		return fmt.Errorf("write collection: failed to parse chapter template: %w", err)
	}

	// --- Static files ---
	if err := copyDirectory(layoutsDir, outputDir, []string{
		"index.html",
		"_book.html",
		"_chapter.html",
	}, []string{
		"_template_*.html",
	}); err != nil {
		return fmt.Errorf("failed to copy files to output. %w", err)
	}

	// --- Collection index ---
	fCollection, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("write collection: failed to create collection index file: %w", err)
	}
	defer fCollection.Close()

	if err := collectionTemplate.ExecuteTemplate(fCollection, "index.html", collection); err != nil {
		return fmt.Errorf("write collection: failed to write to collection index file: %w", err)
	}

	if collection.FaviconImageName != "" {
		oldFaviconPath := filepath.Join(collection.InputDirectory, collection.FaviconImageName)
		newFaviconPath := filepath.Join(outputDir, collection.FaviconImageName)
		if err := copyFile(oldFaviconPath, newFaviconPath); err != nil {
			return fmt.Errorf("write collection: failed to add favicon image to output: %w", err)
		}
	}

	// --- Write books ---
	g := new(errgroup.Group)
	for _, book := range collection.Books {
		bookOutputDir := filepath.Join(outputDir, "books", book.UniqueID)

		g.Go(func() error {
			if err := writeBookToHTML(&book, bookOutputDir, "", bookTemplate, chapterTemplate); err != nil {
				return err
			}

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func WriteBookToHTML(book *Book, outputDir, layoutsDir string) error {
	otherTemplatesPath := filepath.Join(layoutsDir, "_template_*.html")

	bookTemplate, err := parseTemplate(filepath.Join(layoutsDir, "_book.html"), otherTemplatesPath)
	if err != nil {
		return fmt.Errorf("write book '%s': failed to parse book template: %w", book.UniqueID, err)
	}

	chapterTemplate, err := parseTemplate(filepath.Join(layoutsDir, "_chapter.html"), otherTemplatesPath)
	if err != nil {
		return fmt.Errorf("write book '%s': failed to parse chapter template: %w", book.UniqueID, err)
	}

	return writeBookToHTML(book, outputDir, layoutsDir, bookTemplate, chapterTemplate)
}

func writeBookToHTML(book *Book, outputDir, layoutsDir string, bookTemplate *template.Template, chapterTemplate *template.Template) error {
	chaptersOutputDir := filepath.Join(outputDir, "chapters")
	if err := os.MkdirAll(chaptersOutputDir, 0755); err != nil {
		return fmt.Errorf("write book '%s': %w", book.UniqueID, err)
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		fIndex, err := os.Create(filepath.Join(outputDir, "index.html"))
		if err != nil {
			return fmt.Errorf("failed to create index.html: %w", err)
		}
		defer fIndex.Close()

		if book.Content.Parsed == nil {
			book.Content.Init()
		}

		book.Content.Parsed["html"], err = convertMarkdownToHTML(book.Content.Raw)
		if err != nil {
			return fmt.Errorf("failed to convert index.html markdown to HTML: %w", err)
		}

		if err := bookTemplate.ExecuteTemplate(fIndex, "_book.html", book); err != nil {
			return fmt.Errorf("failed to write index.html: %w", err)
		}

		return nil
	})

	// --- Cover/favicon images ---
	if book.CoverImageName != "" {
		oldCoverPath := filepath.Join(book.InputDirectory, book.CoverImageName)
		newCoverPath := filepath.Join(outputDir, book.CoverImageName)
		if err := copyFile(oldCoverPath, newCoverPath); err != nil {
			return fmt.Errorf("write book '%s': failed to add cover image to output: %w", book.UniqueID, err)
		}
	}

	if book.FaviconImageName != "" {
		oldFaviconPath := filepath.Join(book.InputDirectory, book.FaviconImageName)
		newFaviconPath := filepath.Join(outputDir, book.FaviconImageName)
		if err := copyFile(oldFaviconPath, newFaviconPath); err != nil {
			return fmt.Errorf("write book '%s': failed to add favicon image to output: %w", book.UniqueID, err)
		}
	}

	// --- Chapters ---
	for _, chapter := range book.Chapters {
		g.Go(func() error {
			fChapter, err := os.Create(filepath.Join(chaptersOutputDir, chapter.UniqueID+".html"))
			if err != nil {
				return fmt.Errorf("failed to create chapter '%s': %w", chapter.UniqueID+".html", err)
			}
			defer fChapter.Close()

			if chapter.Content.Parsed == nil {
				chapter.Content.Init()
			}

			chapter.Content.Parsed["html"], err = convertMarkdownToHTML(chapter.Content.Raw)
			if err != nil {
				return fmt.Errorf("failed to convert chapter '%s' markdown to HTML: %w", chapter.UniqueID+".html", err)
			}

			if err := chapterTemplate.ExecuteTemplate(fChapter, "_chapter.html", &chapter); err != nil {
				return fmt.Errorf("failed to write chapter '%s': %w", chapter.UniqueID+".html", err)
			}

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("write book '%s': %w", book.UniqueID, err)
	}

	return nil
}

func parseTemplate(path, templatesGlob string) (*template.Template, error) {
	templateFileNames := []string{path}
	fileNames, err := filepath.Glob(templatesGlob)
	if err == nil {
		templateFileNames = append(templateFileNames, fileNames...)
	}

	template, err := template.ParseFiles(templateFileNames...)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func convertMarkdownToHTML(data []byte) (template.HTML, error) {
	var buffer bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(data, &buffer, parser.WithContext(context)); err != nil {
		return template.HTML(""), err
	}

	return template.HTML(buffer.String()), nil
}

func copyDirectory(sourceDir, destinationDir string, excludePaths, excludePatterns []string) error {
	return copyDirectoryHelper(sourceDir, destinationDir, sourceDir, excludePaths, excludePatterns)
}

func copyDirectoryHelper(currDir, destinationDir, sourceDir string, excludePaths, excludePatterns []string) error {
	items, err := os.ReadDir(currDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		target := filepath.Join(currDir, item.Name())
		targetFromRoot := strings.TrimLeft(strings.TrimPrefix(target, sourceDir), "/\\")
		newPath := filepath.Join(destinationDir, targetFromRoot)

		// Check exclusions
		if slices.Contains(excludePaths, targetFromRoot) {
			continue
		}

		matching := false
		for _, pattern := range excludePatterns {
			matching, err = filepath.Match(pattern, targetFromRoot)
			if err != nil {
				return err
			}
		}
		if matching {
			continue
		}

		// Copy directories
		if item.IsDir() {
			if err := os.MkdirAll(newPath, 0755); err != nil {
				return err
			}

			if err := copyDirectoryHelper(target, destinationDir, sourceDir, excludePaths, excludePatterns); err != nil {
				return err
			}

			continue
		}

		// Copy item from old to new
		if err := copyFile(target, newPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destination)
	if err != nil {
		return nil
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
