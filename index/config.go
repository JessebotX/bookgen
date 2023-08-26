package index

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JessebotX/bookgen/common"
)

// Unmarshal bookgen's main toml configuration file
func UnmarshalMainConfig(path string, config *common.Config) error {
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
