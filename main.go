package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const Version = "0.1"

var (
	SuppressNonEssentialOutput = false
	EnablePlainOutput          = false
	ProgramName                = "bookgen"
)

func main() {
	// ---
	// Read CLI arguments
	// ---
	// Home-made cli argument parsing
	ProgramName = os.Args[0]

	var workingDirFlag, outputDirFlag string
	setWorkingDirFlag := false
	setOutputDirFlag := false

	// cli arguments are optional
	if len(os.Args) >= 2 {
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]

			if arg == "--" {
				break
			} else if arg == "--plain" {
				EnablePlainOutput = true
			} else if arg == "-q" || arg == "--suppress-non-essential-output" {
				SuppressNonEssentialOutput = true
			} else if arg == "-i" || arg == "--input-directory" {
				// Flag requires arg
				if (i + 1) >= len(os.Args) {
					errorExit(1, "missing value for flag `%v`. See `%v` --help` for more info.", arg, ProgramName)
				}

				workingDirFlag = os.Args[i+1]
				setWorkingDirFlag = true
				i++ // next arg is not a flag
			} else if arg == "-o" || arg == "--output-directory" {
				// Flag requires arg
				if (i + 1) >= len(os.Args) {
					errorExit(1, "missing value for flag `%v`. See `%v --help` for more info.", arg, ProgramName)
				}

				outputDirFlag = os.Args[i+1]
				setOutputDirFlag = true
				i++ // next arg is not a flag
			} else {
				errorExit(1, "unknown flag `%v`. See `%v --help` for more info.", arg, ProgramName)
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
		errorExit(1, "output directory cannot be equal to the working/input directory (`%s` and `%s` reference the same path).", workingDirFlag, outputDirFlag)
	}

	// ---
	// Parse collection
	// ---

	timeStart := time.Now()

	collection, err := DecodeCollection(workingDirFlag)
	if err != nil {
		errorExit(1, err.Error())
	}

	if err := RenderCollectionToWebsite(&collection, workingDirFlag, outputDirFlag); err != nil {
		errorExit(1, err.Error())
	}

	// totalFiles := 0
	// for _, b := range collection.Books {
	// 	for _ = range b.Chapters {
	// 		totalFiles++
	// 	}
	// 	totalFiles++
	// }

	timeElapsed := time.Since(timeStart)

	if !SuppressNonEssentialOutput {
		fmt.Printf(terminalPrintBold("Done")+" (%v)\n", timeElapsed)
	}
}

func errorExit(code int, format string, a ...any) {
	fmt.Fprintf(os.Stderr, terminalPrintBold("bookgen error: ")+format+"\n", a...)
	os.Exit(code)
}

func terminalPrintBold(s string) string {
	if EnablePlainOutput {
		return s
	}

	return "\033[1m" + s + "\033[0m"
}
