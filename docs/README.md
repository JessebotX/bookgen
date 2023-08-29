# Bookgen
## COPYING
Copyright (C) 2023 Free Software Foundation, Inc.

> Permission is granted to copy, distribute and/or modify this document under the terms of the GNU Free Documentation License, Version 1.3 or any later version published by the Free Software Foundation; with no Invariant Sections, with the Front-Cover Texts being “A GNU Manual,” and with the Back-Cover Texts as in (a) below. A copy of the license is included in the section entitled “GNU Free Documentation License.”
>
> (a) The FSF’s Back-Cover Text is: “You have the freedom to copy and modify this GNU manual.”

## Overview
Bookgen is a static site generator designed for authors who want to distribute their _collection_ of Markdown-based written _works_ as a standalone static website. It additionally creates RSS feeds---for readers to keep up to date with new chapters---and ebook (`.epub`) files---allowing readers to download your DRM-free written works for offline reading.

Bookgen is a cross-platform command-line interface application. It has only been tested on Windows and Linux.

## Bookgen Command-Line Interface Overview
```bash
bookgen
```
* Generates a full static website along with `.epub` files and RSS feeds. _User must be currently in the root folder of their bookgen project._

```bash
bookgen new <project_name>
```
* **REQUIRED FIELD**: `<project_name>`
* Bootstrap a new bookgen project, using `<project_name>` as the name of the root directory.

```bash
bookgen help
```

```bash
bookgen -h
```

```bash
bookgen --help
```
* Printing usage information

```bash
bookgen version
```

```bash
bookgen -V
```

```bash
bookgen --version
```
* Print version of application

## Collection Configuration File
A `bookgen.toml` file should be in the root of your bookgen collection. This is where you can customize internals such as the location of certain directories, as well as specifying main site index settings, such as the title of the main site, base URL, etc.

**NOTE: ALL DIRECTORIES ARE RELATIVE TO THE ROOT DIRECTORY OF THE BOOKGEN COLLECTION**

### Collection Configuration File Reference
```toml
### NOTE: commented fields are optional
### NOTE: **ALL DIRECTORIES IN CONFIG RELATIVE TO ROOT OF COLLECTION

# directory storing all books' sources
#booksDir  = "./books"
# directory containing layout files
#themeDir  = "./themes"
# where to output the generated static site
#outputDir = "./out"
#staticDir = "./static"

[index] # main site index settings
# the main site title
title   = "MAIN_SITE_TITLE"
# the base domain name url (important for RSS feeds)
baseURL = "https://example.com"
# language code of the main site collection index page
languageCode = "en"
```

## Book Configuration File
Each book in a collection (subdirectories of the collection's `booksDir`) should contain a `bookgen-book.toml`

### Book Configuration File Reference
```toml
title = "BOOK TITLE NAME"
# One or two sentences describing what the book is about
shortDescription = "SHORT DESCRIPTION"
# Genre of the book
genre = "GENRE"
# The current status of the book, whether it has been completed, currently ongoing, on hiatus, or dropped
status = "Ongoing|On Hiatus|Completed"
# Path to cover image (RELATIVE TO ROOT OF BOOK DIRECTORY)
coverPath = "./cover.png"
# Directory containing chapters and static assets related to the book (RELATIVE TO ROOT OF BOOK DIRECTORY)
#chaptersDir = "./chapters"
# Book language code
languageCode = "en"

# author name
author  = "AUTHOR_NAME"
# a short "about me" of the author
bio     = """
ABOUT THE AUTHOR
"""

# You're able to specify as many ways to donate as you want

[[author.donation]]
#name = "DONATION_SITE_NAME (Patreon, Paypal, etc.)"
#link = "LINK_TO_DONATION_PAGE"

[[author.donation]]
#name = "DONATION_SITE_NAME (Patreon, Paypal, etc.)"
#link = "LINK_TO_DONATION_PAGE"

# You can technically specify cryptocurrency wallet addresses instead
[[author.donation]]
#name = "CRYPTOCURRENY NAME (Bitcoin, Monero, etc.)"
#link = "CRYPTOCURRENCY_WALLET_ADDRESS"
```

## Chapter Configuration Frontmatter
### Chapter Configuration Frontmatter Reference
```yaml
# title of the chapter
title:  "CHAPTER_TITLE"
# short chapter description
description: "DESCRIPTION HERE"
# the published date (MUST BE IN THIS FORMAT)
# - Format: YYYY-mm-ddTHH:MM:SS[timezone offset]
date:    2006-01-02T15:04:05-07:00
# last modified date (MUST BE IN THIS FORMAT)
# - Format: YYYY-mm-ddTHH:MM:SS[timezone offset]
lastmod: 2006-01-02T15:04:05-07:00
```

## Future
Things that are nice to have possibly in the future.

* [ ] Epub CSS styling and other features
* [ ] Generate [Gemini capsule](https://gemini.circumlunar.space/)
  * [ ] Generate [gempub](https://codeberg.org/oppenlab/gempub)
