package pub

import (
	"errors"
	"fmt"
)

var (
	ErrProfileMissingName = errors.New("profile: missing Name")
)

// Profile may represent an individual or an organization that is credited as either an author, contributor or publisher affliated with a [Book].
type Profile struct {
	Name           string      `json:"name"`
	NamesAlternate []string    `json:"names_alternate"`
	Content        Content     `json:"content"`
	External       []Reference `json:"external"`
}

func (p *Profile) EnsureValid() error {
	if p.Name == "" {
		return ErrProfileMissingName
	}

	for _, e := range p.External {
		if err := e.EnsureValid(); err != nil {
			return fmt.Errorf("profile \"%s\": %w", p.Name, err)
		}
	}

	return nil
}
