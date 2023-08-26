# Markdown Book Generator
Generate a static website for a collection of books written in Markdown, including RSS feeds and EPUB files for each book.

# Usage and Structure
```sh
ROOT_FOLDER/
  ...
  bookgen.toml
  theme/
    index.html
    book-chapter.html
    book-index.html
  static/
    ... # static files like css, js, etc.
  books/
    BOOK_1/
      bookgen-book.toml
      cover.jpg
      index.md
      chapters/
    BOOK_2/
      ...
    ...
```
