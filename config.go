package mkpub

import (
	"html/template"
	"net/url"
	"strings"
	"time"
)

var BookValidStatusValues = []string{"completed", "hiatus", "inactive", "ongoing"}

type Content struct {
	Raw  []byte
	HTML template.HTML
}

type Internal struct {
	GenerateEPUB bool
	GenerateRSS  bool
}

func (i *Internal) Init() {
	i.GenerateEPUB = true
	i.GenerateRSS = true
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
	ConfigFormat     string
	Params           map[string]any
	Internal         Internal
	Books            []Book
	Title            string
	Description      string
	BaseURL          string
	LanguageCode     string
	Content          Content
	FaviconImageName string
	LastBuildDate    time.Time
}

func (c *Collection) InitDefaults() {
	c.Internal.Init()

	c.LanguageCode = "en"
	c.LastBuildDate = time.Now()
}

type Book struct {
	ConfigFormat     string
	Params           map[string]any
	Internal         Internal
	Parent           *Collection
	UniqueID         string
	LastBuildDate    time.Time
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

func (b *Book) InitDefaults(uniqueID string, parent *Collection) {
	// Defaults
	b.Internal.Init()

	b.UniqueID = uniqueID
	b.LastBuildDate = time.Now()

	// Inherited from parent
	if parent != nil {
		b.Parent = parent
		b.Internal.GenerateRSS = parent.Internal.GenerateRSS
		b.Internal.GenerateEPUB = parent.Internal.GenerateEPUB

		if strings.TrimSpace(parent.BaseURL) != "" {
			b.BaseURL, _ = url.JoinPath(parent.BaseURL, "books", uniqueID)
		}

		if strings.TrimSpace(b.LanguageCode) == "" {
			b.LanguageCode = parent.LanguageCode
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

func (c *Chapter) InitDefaults(uniqueID string, parent *Book) {
	// Required
	c.Parent = parent
	c.UniqueID = uniqueID
	c.Title = uniqueID

	// Inherited from parent
	c.LanguageCode = parent.LanguageCode
}
