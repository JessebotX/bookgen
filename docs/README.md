# Bookgen
## COPYING
Copyright (C) 2023 Free Software Foundation, Inc.

> Permission is granted to copy, distribute and/or modify this document under the terms of the GNU Free Documentation License, Version 1.3 or any later version published by the Free Software Foundation; with no Invariant Sections, with the Front-Cover Texts being “A GNU Manual,” and with the Back-Cover Texts as in (a) below. A copy of the license is included in the section entitled “GNU Free Documentation License.”
>
> (a) The FSF’s Back-Cover Text is: “You have the freedom to copy and modify this GNU manual.”

## Overview
Bookgen is a static site generator designed for authors who want to distribute their _collection_ of Markdown-based written _works_ as a standalone static website. It additionally creates RSS feeds---for readers to keep up to date with new chapters---and ebook (`.epub`) files---allowing readers to download your DRM-free written works for offline reading.

## Features
* Write in Markdown and generate a full static website ready to be deployed
* Themes: Fully customize web page generation utilizing [Go's builtin text/html templating engine](https://pkg.go.dev/html/template)
* Works out of the box with a default theme
* Generates RSS feeds for each book, allowing users to keep up to date with the latest releases
* Generates EPUB file downloads, allowing users to read all published chapters offline on any device
* Reading time and word count estimates for each chapter

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

### Examples
For further clarification, see the examples in [`testdata/`](../testdata), as well as the default templates seen in [`renderer/templates.go`](../renderer/templates.go)

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

## Concepts
Bookgen is structured in the following way:

* **Chapter**: a chapter is a section in a book
  * HTML generation controlled by `layout/chapter.html`
* **Book**: a book is a written work containing 1 or more chapters
  * HTML generation controlled by `layout/book.html`
* **Index (or "collection")**: an index (or "collection") refers to the top level object that contains one or more books
  * HTML generation controlled by `layout/index.html`

## Reference
### `bookgen.toml`
```toml
# Main title of the collection
title        = "John Doe's Book Collection"
# Base URL of this site, used primarily for RSS feed generation so it can point to the correct chapter addresses
baseURL      = "https://johndoe.xyz/books"
# Language of the main landing page (NOT ALWAYS the language of all books). Specify as ISO 639-1 Language code. See <https://www.w3schools.com/tags/ref_language_codes.asp> for a list of them.
languageCode = "en"
```

### `bookgen-book.toml`
```toml
# Title of the book
title            = "The Adventure's of John Doe"
# A short one/two sentence description of what the book is about
shortDescription = """
Follow John Doe as he hitches a ride on a spaceship and explores the universe.
"""
# Genre[s] the book can be categorized into
genre            = ["Science Fiction", "Adventure"]
# The status of the written work, whether serialization is "On Hiatus", "Ongoing", "Completed", etc.
status           = "Completed"
# File name of the cover (NOTE: there shouldn't be directories in this path as the cover is in the book's root directory)
coverPath        = "cover.png"
# Language of the book. Specify as ISO 639-1 Language code. See <https://www.w3schools.com/tags/ref_language_codes.asp> for a list of them.
languageCode     = "en"
# A copyright notice
copyright        = "Copyright John Doe"
# Book's content license for copyright purposes
license          = "CC0"
# Other places where you can read this book
mirrors          = ["https://anotherwebsite.com/johndoebook", "https://mirrorwebsite.net/johndoe/1"]

[author]
# Name of the author
name = "John Doe"
# About the author
bio = """
John Doe is a bestselling author, astrophysicist, astronaut, mathematician, computer scientist, biochemist, psychologist, engineer, clothing designer, and a entrepreneur known for his contributions in adventure sci-fi with his debut work \"The Adventure's of John Doe\".
"""

# You're able to specify as many ways to donate as you want using [[author.donate]]
[[author.donate]]
# Name of the site
name = "Patreon"
# Link to where you can donate to the author
link = "https://patreon.com/johndoeuser11111"

[[author.donate]]
# Name of the site
name = "Paypal"
# Link to where you can donate to the author
link = "https://paypal.com/johndoeuser11111"

[[author.donate]] # Supports non link items such as cryptocurrency addresses
# Name of donation method
name = "Bitcoin"
# Since `nonLinkItem = true`, this is no longer treated as a link
link = "1234567890qwertyuiopasdfghjklzxcvbnm"
# Set to true for non-link donations
nonLinkItem = true
```

### Chapter File YAML metadata
In each chapter markdown file, you can specify front matter in yaml format

```yaml
# the title of the chapter
title: Chapter 1 - Where It All Began
# a short description of the chapter
description: "Well, this is how it all began..."
# date and time of published, must be in this format (YYYY-mm-ddTHH:MM:SSZ)
published: 2006-01-02T15:04:05-07:00
# date and time file was changed (YYYY-mm-ddTHH:MM:SSZ)
lastmod: 2023-01-01T12:00:00-07:00
```

### Go Templating Reference
See [`config/config.go`](../config/config.go)

* `Collection` struct is passed into template `index.html`
* `Book` struct is passed into template `book.html`
* `Chapter` struct is passed into template `chapter.html`

## Examples
See [testdata/](../testdata) for example collection projects.
