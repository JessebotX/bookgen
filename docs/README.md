# Bookgen
## COPYING
Copyright (C) 2023 Free Software Foundation, Inc.

> Permission is granted to copy, distribute and/or modify this document under the terms of the GNU Free Documentation License, Version 1.3 or any later version published by the Free Software Foundation; with no Invariant Sections, with the Front-Cover Texts being “A GNU Manual,” and with the Back-Cover Texts as in (a) below. A copy of the license is included in the section entitled “GNU Free Documentation License.”
>
> (a) The FSF’s Back-Cover Text is: “You have the freedom to copy and modify this GNU manual.”

## Overview
Bookgen is a static site generator designed for authors who want to distribute their _collection_ of Markdown-based written _works_ as a standalone static website. It additionally creates RSS feeds---for readers to keep up to date with new chapters---and ebook (`.epub`) files---allowing readers to download your DRM-free written works for offline reading.

## Installation
Bookgen is a cross-platform command-line interface application. The easiest way to install it would be to use the `go` CLI.

```bash
go install github.com/JessebotX/bookgen@latest
```

## Getting Started
### Command Line Interface Overview
```bash
bookgen build <ROOT PATH>
```
Compile a project at the root directory path `<ROOT PATH>`

```bash
bookgen help|-h
```
Print usage information

```bash
bookgen version|-V
```
Print the current version of `bookgen`

### File Structure
Begin by creating the following file/folder structure:

```bash
<ROOT> (aka. collection, project)/
  bookgen.toml
  static/
    # ...static assets here...
  layout/
    book.html
    index.html
    chapter.html
    static/
      style.css
      # ...
  src/
    book 1 identifier/
      bookgen-book.toml
      cover.png
      index.md
      chapters/
        chapter-1.md
        chapter-2.md
        # ...
  # ...
```

* **Root/collection/project/**: the root project directory that contains all necessary items to generate a complete static website
* `bookgen.toml`: Main collection configuration file. See the bookgen.toml reference
* `static/`: any assets to be included here (e.g. include a favicon). If they are a part of the layout/theme (i.e. css files), put them in `layout/static/` instead
  * Upon building the website, all static files and folders will be placed in the root folder of the output directory, preserving its source-form file structure
* `layout/`: the website theme that changes the way HTML is generated, utilizing Go's builtin [text/html templating engine](https://pkg.go.dev/html/template)
* `layout/chapter.html`: HTML template controlling the generation of a chapter of a book
* `layout/book.html`: HTML template controlling the generation of the book index page which allows users to navigate into the chapters
* `layout/index.html`: HTML template controlling the generation of the main index/collection homepage, which displays a list of books available on the website
* `layout/static/`: folder containing static assets such as CSS that are an essential part of the theme. Without it, the website would be incorrectly generated
  * Upon building the website, all static files and folders will be placed in the root folder of the output directory, preserving its source-form file structure
* `src/`: folder containing 1 or more subdirectories of books/written works
* `src/<book identifier>.../`: folder containing the necessary items to generate a single book that is a part of the full collection. Book identifier can be any valid folder name supported by your operating system and able to be parsed as a URL. Create as many of these subdirectories as you need.
* `src/<book identifier>.../bookgen-book.toml`: Book metadata file. See the bookgen-book.toml reference
* `src/<book identifier>.../<cover image>`: A cover image for your book. The filetype must be supported by web browsers (e.g. webp, png, jpg, gif, etc.)
* `src/<book identifier>.../index.md`: A markdown file that contains the full book blurb and nothing more
* `src/<book identifier>.../chapters/`: directory containing a list of markdown files and any static assets
  * Chapter files should not be further nested in subdirectories of `chapters/`
  * Static assets should have unique names, even if they are in different directories
  * Limitation: static assets must not have an `.md` file extension as it is reserved for chapter files
* `src/<book identifier>.../chapters/<chapter files>.md`: Markdown files containing chapter content. Each markdown file should have a yaml metadata section at the top. See the reference for more information
* `out/` (AUTOMATICALLY CREATED): the folder that contains a fully generated static website with RSS feeds and EPUB files included. Created after running `bookgen build <ROOT PATH>`

### Concepts
Bookgen is structured in the following way:

* **Chapter**: a chapter is a section in a book
* **Book**: a book is a written work containing 1 or more chapters
* **Index (or "collection")**: an index (or "collection") refers to the top level object that contains one or more books
