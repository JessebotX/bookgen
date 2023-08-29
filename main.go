// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"os"

	"github.com/JessebotX/bookgen/config"
)

func init() {
	log.Println("Hello, world!")
}

func main() {
	config := config.Bookgen{
		BooksDir:  "./src",
		StaticDir: "./static",
		ThemeDir:  "./theme",
		OutputDir: "./out",
	}

	
}
