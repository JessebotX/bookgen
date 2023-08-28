package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JessebotX/bookgen/common"
	"github.com/JessebotX/bookgen/index"
)

const Version = "0.1.3"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide directory containing a bookgen.toml")
	}

	if os.Args[1] == "-V" {
		fmt.Println("bookgen", Version)
		return
	}

	configPath := os.Args[1]

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

	err = index.GenerateHTMLSiteFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
}

func printVerbose(config *common.Config) {
	log.Printf("%#v\n", config)
}
