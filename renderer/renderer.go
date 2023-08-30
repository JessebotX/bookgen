// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package renderer

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JessebotX/bookgen/config"
	"github.com/bmaupin/go-epub"
)

func BuildSite(collection *config.Collection) error {
	resolvedOutputDir := filepath.Join(collection.Root, collection.OutputDir)
	resolvedLayoutDir := filepath.Join(collection.Root, collection.LayoutDir)

	// clean out directories
	err := os.RemoveAll(resolvedOutputDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(resolvedOutputDir, 0755)
	if err != nil {
		return err
	}

	// read book and chapter templates
	bookTemplatePath := filepath.Join(resolvedLayoutDir, "book.html")

	bookTemplate := template.Must(template.New("book").Parse(BookDefaultTemplate))
	if exists(bookTemplatePath) {
		bookTemplate, err = template.ParseFiles(bookTemplatePath)
		if err != nil {
			return err
		}
	}

	chapterTemplatePath := filepath.Join(resolvedLayoutDir, "chapter.html")
	chapterTemplate := template.Must(template.New("chapter").Parse(ChapterDefaultTemplate))
	if exists(chapterTemplatePath) {
		bookTemplate, err = template.ParseFiles(chapterTemplatePath)
		if err != nil {
			return err
		}
	}

	rssTemplate := template.Must(template.New("rss").Parse(RSSTemplate))

	// create book index
	for _, bk := range collection.Books {
		bookOutputDir := filepath.Join(resolvedOutputDir, bk.ID)
		err = os.MkdirAll(bookOutputDir, 0755)
		if err != nil {
			return err
		}

		// init epub building
		e := epub.NewEpub(bk.Title)
		e.SetAuthor(bk.Author.Name)

		blurbBody := "<h1>" + bk.Title + "</h1>" + string(bk.Blurb)
		e.AddSection(blurbBody, bk.Title, "", "")

		// cover image
		resolvedCoverPath := filepath.Join(bk.Root, bk.CoverPath)
		if bk.CoverPath != "" && exists(resolvedCoverPath) {
			log.Println("Found cover")
			outputCoverPath := filepath.Join(bookOutputDir, bk.CoverPath)
			err = os.Link(resolvedCoverPath, outputCoverPath)
			if err != nil {
				return err
			}

			coverImageEpubPath, err := e.AddImage(resolvedCoverPath, "")
			if err != nil {
				return err
			}
			e.SetCover(coverImageEpubPath, "")
		}

		// static assets in book
		for _, assetPath := range bk.StaticAssets {
			_, err = e.AddImage(assetPath, "")
			if err != nil {
				return err
			}

			resolvedPath := filepath.Join(bookOutputDir, strings.TrimPrefix(assetPath, filepath.Join(bk.Root, bk.ChaptersDir)))
			err = os.MkdirAll(filepath.Dir(resolvedPath), 0755)
			if err != nil {
				return err
			}

			err = os.Link(assetPath, resolvedPath)
			if err != nil {
				return err
			}
		}

		// book index
		bookIndexFile, err := os.Create(filepath.Join(bookOutputDir, "index.html"))
		if err != nil {
			return err
		}

		err = bookTemplate.Execute(bookIndexFile, &bk)
		if err != nil {
			return err
		}

		// rss feed
		rssXML, err := os.Create(filepath.Join(bookOutputDir, "rss.xml"))
		if err != nil {
			return err
		}

		err = rssTemplate.Execute(rssXML, &bk)
		if err != nil {
			return err
		}

		// chapters
		for _, chapter := range bk.Chapters {
			chapterFile, err := os.Create(filepath.Join(bookOutputDir, chapter.SlugHTML()))
			if err != nil {
				return err
			}

			err = chapterTemplate.Execute(chapterFile, &chapter)
			if err != nil {
				return err
			}

			sectionBody := "<h1>" + chapter.Title + "</h1>" + string(chapter.Content)
			e.AddSection(sectionBody, chapter.Title, "", "")
		}

		err = e.Write(filepath.Join(bookOutputDir, bk.ID+".epub"))
		if err != nil {
			return err
		}
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
