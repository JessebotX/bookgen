package pub

import (
	"errors"
	"fmt"
)

var (
	ErrSeriesMissingTitle = errors.New("series: missing title")
)

// Series describes a [Book]'s relation to a set of other [Book] objects (i.e. prequels, sequels, side stories, sharing the same world/universe, etc.)
type Series struct {
	Title           string      `json:"title"`
	TitlesAlternate []string    `json:"titles_alternate"`
	Number          float64     `json:"number"`
	Description     string      `json:"description"`
	Content         Content     `json:"content"`
	External        []Reference `json:"external"`
}

func (s Series) EnsureValid() error {
	if s.Title == "" {
		return ErrSeriesMissingTitle
	}

	for _, e := range s.External {
		if err := e.EnsureValid(); err != nil {
			return fmt.Errorf("series \"%s\": %w", s.Title, err)
		}
	}

	return nil
}
