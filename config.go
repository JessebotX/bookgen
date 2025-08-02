package mkpub

import (
	"net/url"
	"strings"
	"time"
)

const (
	BookStatusCompleted = "completed"
	BookStatusInactive  = "inactive"
	BookStatusOngoing   = "ongoing"
	BookStatusHiatus    = "hiatus"
)

var BookStatusValues = []string{
	BookStatusCompleted,
	BookStatusHiatus,
	BookStatusInactive,
	BookStatusOngoing,
}

type Content struct {
	Raw    []byte
	Parsed map[string]any
}

func (c *Content) Init() {
	c.Parsed = make(map[string]any, 1)
}

func (c *Content) Format(key string) any {
	return c.Parsed[key]
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
	Params  map[string]any
	Content Content

	ID           string
	Name         string
	About        string
	EmailAddress string
	Role         string
	Links        []ExternalLink
}

func (a *Author) Init() {
	a.Content.Init()
}

type ExternalLink struct {
	Name           string
	Address        string
	IsAddressPlain bool
}

type Series struct {
	ID     string
	Name   string
	Number float64
}

type Collection struct {
	Params         map[string]any
	DateLastBuild  time.Time
	InputDirectory string
	Content        Content

	Format           string
	Internal         Internal
	Books            []Book
	Title            string
	Description      string
	BaseURL          string
	LanguageCode     string
	FaviconImageName string
}

func (c *Collection) InitDefaults(inputDir string) {
	c.Internal.Init()

	c.LanguageCode = "en"
	c.DateLastBuild = time.Now()
	c.InputDirectory = inputDir
}

type Book struct {
	Params         map[string]any
	Internal       Internal
	Parent         *Collection
	UniqueID       string
	DateLastBuild  time.Time
	Chapters       []Chapter
	InputDirectory string
	Content        Content

	Format             string
	DatePublishedEnd   time.Time
	DatePublishedStart time.Time
	Title              string
	Subtitle           string
	TitleSort          string
	Status             string
	BaseURL            string
	Description        string
	LanguageCode       string
	Series             Series
	Authors            []Author
	AuthorsSort        string
	Contributors       []Author
	Mirrors            []ExternalLink
	IDs                []string
	Tags               []string
	FaviconImageName   string
	CoverImageName     string
}

func (b *Book) InitDefaults(uniqueID, inputDir string, parent *Collection) {
	// Defaults
	b.Internal.Init()
	b.Content.Init()

	b.DateLastBuild = time.Now()
	b.UniqueID = uniqueID
	b.Status = "completed"
	b.InputDirectory = inputDir

	// Inherited from parent
	if parent != nil {
		b.DateLastBuild = parent.DateLastBuild
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
	DateLastBuild time.Time
	Next          *Chapter
	Previous      *Chapter
	UniqueID      string
	Content       Content

	DatePublished time.Time
	DateModified  time.Time
	Title         string
	Subtitle      string
	Description   string
	LanguageCode  string
	Authors       []Author
	Mirrors       []ExternalLink
	Order         int
	Draft         bool
}

func (c *Chapter) InitDefaults(uniqueID string, parent *Book) {
	c.Content.Init()

	// Required
	c.Parent = parent
	c.UniqueID = uniqueID
	c.Title = uniqueID

	// Inherited from parent
	c.DateLastBuild = parent.DateLastBuild
	c.LanguageCode = parent.LanguageCode
	c.Authors = parent.Authors
}
