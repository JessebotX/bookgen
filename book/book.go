package book

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JessebotX/bookgen/common"
)

// Unmarshal a single-book configuration
func UnmarshalBookConfig(path string, config *common.Book) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(source, config)
	if err != nil {
		return err
	}

	return nil
}
