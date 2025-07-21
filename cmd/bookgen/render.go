package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/JessebotX/bookgen"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"

	"golang.org/x/sync/errgroup"
)

const (
	DirPerms = 0755
)

var (
	globalMinifier = minify.New()
)

func RenderCollectionToWebsite(c *bookgen.Collection, workingDir, outputDir string, enableMinify bool) error {
	globalMinifier.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepSpecialComments: true,
		KeepQuotes:          true,
	})
	globalMinifier.AddFunc("text/css", css.Minify)
	globalMinifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	globalMinifier.AddFunc("image/svg+xml", svg.Minify)

	layoutsDir := filepath.Join(workingDir, c.Internal.LayoutsDirectory)
	collectionTemplatePath := filepath.Join(layoutsDir, "index.html")
	bookTemplatePath := filepath.Join(layoutsDir, "_book.html")
	chapterTemplatePath := filepath.Join(layoutsDir, "_chapter.html")

	if err := os.MkdirAll(outputDir, DirPerms); err != nil {
		return fmt.Errorf("failed to create output directory. %w", err)
	}

	// ---
	// Copy global static items into output
	// ---
	if err := copyStaticFilesToDir(layoutsDir, outputDir, layoutsDir, []string{
		"index.html",
		"_book.html",
		"_chapter.html",
	}, []string{
		"_template_*.html",
	}); err != nil {
		return fmt.Errorf("failed to copy files to output. %w", err)
	}

	// ---
	// Read templates
	// ---
	templateFileNames := []string{collectionTemplatePath}
	fileNames, err := filepath.Glob(filepath.Join(layoutsDir, "_template_*.html"))
	if err == nil {
		templateFileNames = append(templateFileNames, fileNames...)
	}

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
	outputIndexPath := filepath.Join(outputDir, "index.html")

	outIndex, err := os.Create(outputIndexPath)
	if err != nil {
		return fmt.Errorf("failed to create collection index file. %w", err)
	}

	if err := collectionTemplate.ExecuteTemplate(outIndex, "index.html", c); err != nil {
		return fmt.Errorf("failed to write collection index file. %w", err)
	}

	if enableMinify {
		if err := minifyFileHTML(outputIndexPath, outIndex); err != nil {
			return fmt.Errorf("failed to minify collection index output file. %w", err)
		}
	}

	if err := outIndex.Close(); err != nil {
		return fmt.Errorf("failed to close collection index file. %w", err)
	}

	// TODO: epub generation
	for _, book := range c.Books {
		bookWorkingDir := filepath.Join(workingDir, "books", book.PageName)
		bookOutputDir := filepath.Join(outputDir, "books", book.PageName)
		if err := os.MkdirAll(bookOutputDir, DirPerms); err != nil {
			return fmt.Errorf("failed to create book `%v` directory. %w", book.PageName, err)
		}

		bookOutputPath := filepath.Join(bookOutputDir, "index.html")
		outBook, err := os.Create(bookOutputPath)
		if err != nil {
			return fmt.Errorf("failed to create book `%v` index file. %w", book.PageName, err)
		}
		defer outBook.Close()

		if err := bookTemplate.ExecuteTemplate(outBook, "_book.html", book); err != nil {
			return fmt.Errorf("failed to write book `%v` index file. %w", book.PageName, err)
		}

		if enableMinify {
			if err := minifyFileHTML(bookOutputPath, outBook); err != nil {
				return fmt.Errorf("failed to minify book `%v` index output file. %w", book.PageName, err)
			}
		}

		g := new(errgroup.Group)
		g.Go(func() error {
			if err := renderBookChapters(book.Chapters, chapterTemplate, chapterTemplatePath, bookOutputDir, enableMinify); err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return fmt.Errorf("failed to write book `%v` chapter file. %w", book.PageName, err)
		}

		// Add cover image to output
		if strings.TrimSpace(book.CoverImageName) != "" {
			coverPathOld := filepath.Join(bookWorkingDir, book.CoverImageName)
			coverPathNew := filepath.Join(bookOutputDir, book.CoverImageName)

			if err := os.RemoveAll(coverPathNew); err != nil {
				return fmt.Errorf("failed to remove cover path in output directory `%v`. %w", err)
			}

			if err := os.Link(coverPathOld, coverPathNew); err != nil {
				return err
			}
		}
	}

	return nil
}

func renderBookChapters(chapters []bookgen.Chapter, chapterTemplate *template.Template, chapterTemplatePath, bookOutputDir string, enableMinify bool) error {
	for _, chapter := range chapters {
		chapterOutputPath := filepath.Join(bookOutputDir, chapter.PageName+".html")

		fChapter, err := os.Create(chapterOutputPath)
		if err != nil {
			return err
		}
		defer fChapter.Close()

		if err := chapterTemplate.ExecuteTemplate(fChapter, "_chapter.html", chapter); err != nil {
			return err
		}

		if enableMinify {
			if err := minifyFileHTML(chapterOutputPath, fChapter); err != nil {
				return fmt.Errorf("failed to minify book `%v` chapter `%v` index output file. %w", chapter.Parent.PageName, chapter.PageName, err)
			}
		}
	}

	return nil
}

func copyStaticFilesToDir(currDir, newDir, rootDir string, relExcludes, relExcludesPatterns []string) error {
	items, err := os.ReadDir(currDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		oldPath := filepath.Join(currDir, item.Name())
		oldPathFromRoot := strings.TrimLeft(strings.TrimPrefix(oldPath, rootDir), "/\\")
		newPath := filepath.Join(newDir, oldPathFromRoot)

		// Check against exclusions
		if slices.Contains(relExcludes, oldPathFromRoot) {
			continue
		}

		matching := false
		for _, pattern := range relExcludesPatterns {
			matching, err = filepath.Match(pattern, oldPathFromRoot)
			if err != nil {
				return err
			}
		}
		if matching {
			continue
		}

		if item.IsDir() {
			if err := os.MkdirAll(newPath, DirPerms); err != nil {
				return err
			}

			if err := copyStaticFilesToDir(oldPath, newDir, rootDir, relExcludes, relExcludesPatterns); err != nil {
				return err
			}

			continue
		}

		// Copy item from old to new path (re-copy if already exists)
		if err := os.RemoveAll(newPath); err != nil {
			return err
		}

		if err := os.Link(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

func minifyFileHTML(path string, f *os.File) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	mb, err := globalMinifier.Bytes("text/html", b)
	if err != nil {
		return err
	}

	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	if _, err := f.Write(mb); err != nil {
		return err
	}

	return nil
}
