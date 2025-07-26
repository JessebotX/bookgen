package mkpub

import (
	"html/template"
	"net/url"
	"strings"
	"time"
)

var BookValidStatusValues = []string{"completed", "hiatus", "inactive", "ongoing"}

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
	ConfigVersion    int
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

func (c *Collection) Init(title, lang string) {
	c.Internal.Init()

	// Required
	c.Title = title
	c.LanguageCode = lang
}

type Book struct {
	ConfigVersion    int
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
	Status           string
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

func (b *Book) Init(parent *Collection, uniqueID, title, lang string) {
	// Required
	b.UniqueID = uniqueID
	b.Title = title
	b.LanguageCode = lang

	// Defaults
	b.Internal.Init()

	b.TitleSort = title
	b.DatePublished = time.Now()
	b.DateModified = b.DatePublished

	// Inherited from parent
	if parent != nil {
		b.Parent = parent
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

func (c *Chapter) Init(parent *Book, uniqueID, title string) {
	// Required
	c.Parent = parent
	c.UniqueID = uniqueID
	c.Title = title

	// Inherited from parent
	c.LanguageCode = parent.LanguageCode
}
