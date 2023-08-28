# Bookgen
Bookgen is a static site generator designed for authors who want to distribute their collection of Markdown-based written works as a standalone static website. It additionally creates RSS feeds---for readers to keep up to date with new chapters---and ebook (`.epub`) files---allowing readers to download your DRM-free written works for offline reading.

Bookgen is a cross-platform command-line interface application. It has only been tested on Windows and Linux.

## COPYING
Copyright (C) 2022-2023 Free Software Foundation, Inc.

    Permission is granted to copy, distribute and/or modify this document under the terms of the GNU Free Documentation License, Version 1.3 or any later version published by the Free Software Foundation; with no Invariant Sections, with the Front-Cover Texts being “A GNU Manual,” and with the Back-Cover Texts as in (a) below. A copy of the license is included in the section entitled “GNU Free Documentation License.”

    (a) The FSF’s Back-Cover Text is: “You have the freedom to copy and modify this GNU manual.”

## Getting Started
1. Install Bookgen for your preferred operating system, or by using the `go` command line tool
2. [ ] TODO

## Bookgen Command-Line Interface Overview
```bash
bookgen
```
Generates a full static website along with `.epub` files and RSS feeds. _User must be currently in the root folder of their bookgen project._

```bash
bookgen new <project_name>
```
* `<project_name>` (REQUIRED)

Bootstrap a new bookgen project, using `<project_name>` as the name of the root directory.

```bash
bookgen help
# OR
bookgen -h
# OR
bookgen --help
```
Print usage information.

```bash
bookgen version
# OR
bookgen -V
# OR
bookgen --version
```
Print the current version of the application
