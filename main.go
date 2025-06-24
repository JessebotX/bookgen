package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// ---
	// Read CLI arguments
	// ---
	// Home-made cli argument parsing
	var workingDirFlag, outputDirFlag string
	setWorkingDirFlag, setOutputDirFlag := false, false

	// cli arguments are optional
	if len(os.Args) >= 2 {
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]

			if arg == "--" {
				break
			} else if arg == "-i" || arg == "--input-dir" || arg == "--input-directory" {
				// Flag requires arg
				if (i + 1) >= len(os.Args) {
					errorExit(1, "missing value for flag '%s", arg)
				}

				workingDirFlag = os.Args[i+1]
				setWorkingDirFlag = true
				i++ // next arg is not a flag
			} else if arg == "-o" || arg == "--output-dir" || arg == "--output-directory" {
				// Flag requires arg
				if (i + 1) >= len(os.Args) {
					errorExit(1, "missing value for flag '%s'", arg)
				}

				outputDirFlag = os.Args[i+1]
				setOutputDirFlag = true
				i++ // next arg is not a flag
			} else {
				errorExit(1, "unknown argument '%s'", arg)
			}
		}
	}

	// set default values for flags
	if !setWorkingDirFlag {
		workingDirFlag = "./"
	}

	// make default output dir generate relative to working dir if
	// output dir is not provided
	if !setOutputDirFlag {
		outputDirFlag = filepath.Join(workingDirFlag, "out")
	}

	// output dir cannot be == working dir
	if filepath.Clean(outputDirFlag) == filepath.Clean(workingDirFlag) {
		errorExit(1, "cannot set output directory to be the same as working/input directory")
	}

	// ---
	// Parse collection
	// ---

	tomlBody := `
tiTle = "Hello world"
title = 1

[hello]
exTra = "Hi"
`

	collection, err := DecodeCollection([]byte(tomlBody))
	if err != nil {
		errorExit(1, err.Error())
	}

	fmt.Printf("%#v\n", collection)
}

func errorExit(code int, format string, a ...any) {
	fmt.Fprintf(os.Stderr, "bookgen error: "+format+"\n", a...)
	os.Exit(code)
}
