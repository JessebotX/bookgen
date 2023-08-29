// Copyright 2023 Jesse
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"os"

	"github.com/JessebotX/bookgen/collection"
)

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument.")
	}
}

func main() {
	collection, err := collection.Create(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", *collection)
}
