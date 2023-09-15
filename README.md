# Bookgen
A powerful static site generator geared towards creatives who want to distribute their written works. Fully customizable and works out of the box.

## Features
* Write in Markdown and generate a full static website ready to be deployed
* Themes: Fully customize web page generation utilizing [Go's builtin text/html templating engine](https://pkg.go.dev/html/template)
* Works out of the box with a default theme
* Generates RSS feeds for each book, allowing users to keep up to date with the latest releases
* Generates EPUB file downloads, allowing users to read all published chapters offline on any device
* Reading time and word count estimates for each chapter

## Usage
See [docs/README.md](docs/README.md) for the complete documentation

### Synopsis
```bash
bookgen help|-h
bookgen version|-V
bookgen build <PATH>
```

## License
See [LICENSE.txt](LICENSE.txt)
