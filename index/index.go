package index

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"html/template"
	"strings"

	"github.com/JessebotX/bookgen/book"
	"github.com/JessebotX/bookgen/common"
	"github.com/bmaupin/go-epub"
)

const RSSTemplateSource = `
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
<title>{{ .Config.Index.Title }}</title>
<link>{{ .Config.Index.BaseURL }}</link>
<description>Recent content for {{ .Title }}</description>
<generator>Bookgen -- github.com/JessebotX/bookgen</generator>
{{ with .LanguageCode }}
<language>{{.}}</language>
{{ end }}
{{ with .Copyright }}
<copyright>{{.}}</copyright>
{{ end }}
{{ range .Chapters }}
<item>

<title>{{ .Title }}</title>
<link>{{ .Config.Index.BaseURL }}/{{ .Parent.Slug }}/{{ .Slug }}.html</link>
<pubDate>{{ .PublishDate.Format "Mon, 02 Jan 2006 15:04:05 -0700" }}</pubDate>
<guid>{{ .Config.Index.BaseURL }}/{{ .Parent.Slug }}/{{ .Slug }}.html</guid>
<description>
{{- with .Description }}
{{ . }}
{{- else -}}
{{ .Title }}
{{ end -}}</description>
</item>
{{ end }}
</channel>
</rss>
`

// Generate configuration file
func GenerateHTMLSiteFromConfig(config *common.Config) error {
	index := config.Index

	// recreate output directory
	err := os.RemoveAll(config.OutputDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		return err
	}

	// copy static files to output
	err = copyStaticDirToOutputDir(config.StaticDir, config.OutputDir)
	if err != nil {
		return err
	}

	// Get books
	books, err := Books(config)
	if err != nil {
		return err
	}

	// read templates
	bookTemplatePath := filepath.Join(config.ThemeDir, "book.html")
	bookTemplate, err := template.ParseFiles(bookTemplatePath)
	if err != nil {
		return err
	}

	rssTemplate, err := template.New("rss").Parse(RSSTemplateSource)
	if err != nil {
		return err
	}

	chapterTemplatePath := filepath.Join(config.ThemeDir, "chapter.html")
	chapterTemplate, err := template.ParseFiles(chapterTemplatePath)
	if err != nil {
		return err
	}

	// create book output
	for i, _ := range books {
		bookItem := &books[i]
		// get all the chapters
		err = book.UnmarshalChapters(bookItem)
		if err != nil {
			return err
		}

		bookOutputDir := filepath.Join(config.OutputDir, bookItem.Slug)
		err = os.MkdirAll(bookOutputDir, 0755)
		if err != nil {
			return err
		}

		// copy cover image to output
		baseCoverPath := filepath.Base(bookItem.CoverPath)
		err = os.Link(bookItem.CoverPath, filepath.Join(bookOutputDir, baseCoverPath))
		if err != nil {
			log.Println("WARN", err)
		}

		// generate book index
		newBookIndexOutput, err := os.Create(filepath.Join(bookOutputDir, "index.html"))
		if err != nil {
			return err
		}

		err = bookTemplate.Execute(newBookIndexOutput, bookItem)
		if err != nil {
			return err
		}

		// Create rss feed
		newRssOutput, err := os.Create(filepath.Join(bookOutputDir, "rss.xml"))
		if err != nil {
			return err
		}

		err = rssTemplate.Execute(newRssOutput, bookItem)
		if err != nil {
			return err
		}

		// generate chapters and epub
		e := epub.NewEpub(bookItem.Title)
		e.SetAuthor(config.Index.Author)

		if bookItem.CoverPath != "" {
			coverImage, err := e.AddImage(bookItem.CoverPath, "")
			if err != nil {
				return err
			} else {
				e.SetCover(coverImage, "")
			}
		}

		indexContent := "<h1>" + bookItem.Title + "</h1>" + string(bookItem.Blurb)
		e.AddSection(indexContent, bookItem.Title, "", "")

		for _, chapter := range bookItem.Chapters {
			newChapterOutput, err := os.Create(filepath.Join(bookOutputDir, chapter.Slug + ".html"))
			if err != nil {
				return err
			}

			err = chapterTemplate.Execute(newChapterOutput, &chapter)
			if err != nil {
				return err
			}

			chapterBody := "<h1>" + string(chapter.Title) + "</h1>" + string(chapter.Content)
			e.AddSection(chapterBody, chapter.Title, "", "")
		}

		err = e.Write(filepath.Join(bookOutputDir, bookItem.Slug + ".epub"))
		if err != nil {
			return err
		}
	}

	index.Books = books
	// TODO generate site index

	return nil
}

// Retrieve all books
func Books(config *common.Config) ([]common.Book, error) {
	booksDir := config.BooksDir

	dirs, err := os.ReadDir(booksDir)
	if err != nil {
		return nil, err
	}

	books := make([]common.Book, 0)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		bookItemDir := filepath.Join(booksDir, dir.Name())
		bookConfig := &common.Book{
			Config:       config,
			LanguageCode: "en",
			IndexPath:    "index.md",
			ChaptersDir:  "chapters",
			Slug:         dir.Name(),
		}
		configPath := filepath.Join(bookItemDir, "bookgen-book.toml")
		err = book.Unmarshal(configPath, bookConfig)
		if err != nil {
			return nil, err
		}

		err := book.UnmarshalBlurb(bookConfig)
		if err != nil {
			return nil, err
		}

		books = append(books, *bookConfig)
	}

	return books, nil
}

// Copy every item in the static files directory to the output
// directory while preserving its structure relative to the static
// directory.
func copyStaticDirToOutputDir(staticDir, outputDir string) error {
	filepath.WalkDir(
		staticDir,
		func(path string, dir fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			base := strings.TrimPrefix(path, staticDir)
			outputPath := filepath.Join(outputDir, base)

			if dir.IsDir() {
				err = os.MkdirAll(outputPath, 0755)
				if err != nil {
					return err
				}

				return nil
			}

			err = os.Link(path, outputPath)
			if err != nil {
				return err
			}

			return nil
		},
	)

	return nil
}
