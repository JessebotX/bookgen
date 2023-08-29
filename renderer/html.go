// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package renderer

import (
	"bytes"
	"html/template"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func MarkdownFileToHTML(filepath string, converter goldmark.Markdown) (template.HTML, map[string]interface{}, error) {
	source, err := os.ReadFile(filepath)
	if err != nil {
		return template.HTML(""), nil, err
	}

	var buffer bytes.Buffer
	context := parser.NewContext()
	err = converter.Convert(source, &buffer, parser.WithContext(context))
	if err != nil {
		return template.HTML(""), nil, err
	}

	return template.HTML(buffer.String()), meta.Get(context), nil
}
