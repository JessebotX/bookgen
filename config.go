package mkpub

import (
	"html/template"
	"net/url"
	"strings"
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

func (i *Internal) Init() {
	i.GenerateEPUB = true
	i.GenerateRSS = true
	i.LayoutDirectory = "layouts"
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

func (c *Collection) Init(title string) {
	c.Internal.Init()

	c.Title = title
	c.LanguageCode = "en"
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

func (b *Book) Init(uniqueID, title string, parent *Collection) {
	b.Internal.Init()

	b.DatePublished = time.Now()
	b.DateModified = time.Now()

	b.UniqueID = uniqueID
	b.Title = title
	b.TitleSort = title
	b.LanguageCode = "en"

	if parent != nil {
		b.Parent = parent
		b.LanguageCode = parent.LanguageCode
		b.Internal.GenerateRSS = parent.Internal.GenerateRSS
		b.Internal.GenerateEPUB = parent.Internal.GenerateEPUB

		if strings.TrimSpace(parent.BaseURL) != "" {
			b.BaseURL, _ = url.JoinPath(parent.BaseURL, "books", uniqueID)
		}
	}
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

func (c *Chapter) Init(uniqueID, title string, parent *Book) {
	c.Parent = parent

	c.UniqueID = uniqueID
	c.Title = title

	c.LanguageCode = parent.LanguageCode

	if !parent.DatePublished.IsZero() {
		c.DatePublished = parent.DatePublished
	}

	if !parent.DateModified.IsZero() {
		c.DateModified = parent.DateModified
	}
}
