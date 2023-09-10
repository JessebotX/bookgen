// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JessebotX/bookgen/collection"
	"github.com/JessebotX/bookgen/renderer"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// Current build version
const Version = "v1.0.0-beta.1"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument.")
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	if cmd == "build" {
		if len(args) <= 0 {
			log.Fatal("Missing collection path for ", cmd)
		}

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

		collection, err := collection.Create(args[0], converter)
		if err != nil {
			log.Fatal(err)
		}

		err = renderer.BuildSite(collection)
		if err != nil {
			log.Fatal(err)
		}
	} else if cmd == "version" || cmd == "-V" {
		fmt.Println("bookgen", Version)
		os.Exit(0)
	} else {
		log.Fatal("Invalid command ", cmd)
	}
}
