package mkpub

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/JessebotX/mkpub/other/meta"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"golang.org/x/sync/errgroup"
)

var (
	goldmarkExtensions = goldmark.WithExtensions(
		meta.Meta,
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
	// Static files
	//
	// ---

	if err := copyFilesToDir(layoutsDir, outputDir, []string{
		"index.html",
		"_book.html",
		"_chapter.html",
	}, []string{
		"_template_*.html",
	}); err != nil {
		return fmt.Errorf("failed to copy files to output. %w", err)
	}

	// ---
	//
	// Collection index
	//
	// ---
	fCollection, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("write collection: failed to create collection index file: %w", err)
	}
	defer fCollection.Close()

	if err := collectionTemplate.ExecuteTemplate(fCollection, "index.html", collection); err != nil {
		return fmt.Errorf("write collection: failed to write to collection index file: %w", err)
	}

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

	chaptersOutputDir := filepath.Join(outputDir, "chapters")
	if err := os.MkdirAll(chaptersOutputDir, 0755); err != nil {
		return fmt.Errorf("write book '%s': %w", book.Title, err)
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		fIndex, err := os.Create(filepath.Join(outputDir, "index.html"))
		if err != nil {
			return fmt.Errorf("failed to create index.html: %w", err)
		}
		defer fIndex.Close()

		book.Content.Parsed["html"], err = convertMarkdownToHTML(book.Content.Raw)
		if err != nil {
			return fmt.Errorf("failed to convert index.html markdown to HTML: %w", err)
		}

		if err := bookTemplate.ExecuteTemplate(fIndex, "_book.html", book); err != nil {
			return fmt.Errorf("failed to write index.html: %w", err)
		}

		return nil
	})

	for _, chapter := range book.Chapters {
		g.Go(func() error {
			fChapter, err := os.Create(filepath.Join(chaptersOutputDir, chapter.UniqueID+".html"))
			if err != nil {
				return fmt.Errorf("failed to create chapter '%s': %w", chapter.UniqueID+".html", err)
			}
			defer fChapter.Close()

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
		return fmt.Errorf("write book '%s': %w", book.Title, err)
	}

	return nil
}

func convertMarkdownToHTML(content []byte) (template.HTML, error) {
	var buffer bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(content, &buffer, parser.WithContext(context)); err != nil {
		return template.HTML(""), err
	}

	return template.HTML(buffer.String()), nil
}

func copyFilesToDir(inputDir, outputDir string, excludePaths, excludePatterns []string) error {
	return copyFilesToDirHelper(inputDir, outputDir, inputDir, excludePaths, excludePatterns)
}

func copyFilesToDirHelper(currDir, newDir, rootDir string, excludePaths, excludePatterns []string) error {
	items, err := os.ReadDir(currDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		target := filepath.Join(currDir, item.Name())
		targetFromRoot := strings.TrimLeft(strings.TrimPrefix(target, rootDir), "/\\")
		newPath := filepath.Join(newDir, targetFromRoot)

		fmt.Println(target, targetFromRoot, newPath)

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

		if item.IsDir() {
			if err := os.MkdirAll(newPath, 0755); err != nil {
				return err
			}

			if err := copyFilesToDirHelper(target, newDir, rootDir, excludePaths, excludePatterns); err != nil {
				return err
			}

			continue
		}

		// Copy item from old to new
		if err := os.RemoveAll(newPath); err != nil {
			return err
		}

		if err := os.Link(target, newPath); err != nil {
			return err
		}
	}

	return nil
}
