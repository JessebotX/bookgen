package mkpub

import (
	"html/template"
	"time"
)

type Content struct {
	Raw []byte
}

func (c *Content) HTML() template.HTML {
	return template.HTML(c.Raw) // TODO: convert c.Raw first
}

type Internal struct {
	GenerateEPUB    bool
	GenerateRSS     bool
	LayoutDirectory string
}

type Author struct {
	Params       map[string]any
	ID           string
	Name         string
	About        Content
	EmailAddress string
	Links        []ExternalLink
}

type ExternalLink struct {
	Name        string
	Address     string
	IsHyperlink bool
}

type Series struct {
	ID     string
	Name   string
	Number float64
}

type Collection struct {
	Params           map[string]any
	Internal         Internal
	Books            []Book
	Title            string
	Description      string
	BaseURL          string
	LanguageCode     string
	Content          Content
	FaviconImageName string
}

type Book struct {
	Params           map[string]any
	Internal         Internal
	Parent           *Collection
	UniqueID         string
	DatePublished    time.Time
	DateModified     time.Time
	Chapters         []Chapter
	Title            string
	Subtitle         string
	TitleSort        string
	BaseURL          string
	Description      string
	LanguageCode     string
	Content          Content
	Series           Series
	Authors          []Author
	AuthorsSort      string
	Contributors     []Author
	Mirrors          []ExternalLink
	IDs              []string
	Subjects         []string
	FaviconImageName string
	CoverImageName   string
}

type Chapter struct {
	Params        map[string]any
	Parent        *Book
	DatePublished time.Time
	DateModified  time.Time
	Next          *Chapter
	Previous      *Chapter
	UniqueID      string
	Title         string
	Subtitle      string
	Description   string
	LanguageCode  string
	Content       Content
	Authors       []Author
	Mirrors       []ExternalLink
	Order         int
	Draft         bool
}
