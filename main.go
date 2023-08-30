// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"os"

	"github.com/JessebotX/bookgen/collection"
	"github.com/JessebotX/bookgen/renderer"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

const Version = "1.0.0"

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument.")
	}
}

func main() {
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
		goldmark.WithRendererOptions(),
	)

	collection, err := collection.Create(os.Args[1], converter)
	if err != nil {
		log.Fatal(err)
	}

	err = renderer.BuildSite(collection)
	if err != nil {
		log.Fatal(err)
	}
}
