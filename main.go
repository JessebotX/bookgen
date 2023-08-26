package main

import (
	"log"
	"path/filepath"

	"github.com/JessebotX/bookgen/common"
	"github.com/JessebotX/bookgen/index"
)

func main() {
	// TODO: remove hardcoding and read from os.Args
	configPath := "./testdata/collection1"

	// set default values in main bookgen config
	// NOTE: paths are relative to the directory of bookgen.toml
	config := &common.Config{
		StaticDir: filepath.Join(configPath, "./static"),
		BooksDir:  filepath.Join(configPath, "./books"),
		ThemeDir:  filepath.Join(configPath, "./theme"),
		OutputDir: filepath.Join(configPath, "./out"),
	}

	err := index.UnmarshalMainConfig(filepath.Join(configPath, "bookgen.toml"), config)
	if err != nil {
		log.Fatal(err)
	}

	printVerbose(config)

	err = index.GenerateHTMLSiteFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
}

func printVerbose(config *common.Config) {
	log.Printf("%#v\n", config)
}
