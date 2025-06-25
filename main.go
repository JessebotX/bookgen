package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const Version = "0.1"

var (
	EnablePlainOutput = false
	ProgramName       = "bookgen"
)

func main() {
	// ---
	// Read CLI arguments
	// ---
	// Home-made cli argument parsing
	ProgramName = os.Args[0]

	var workingDirFlag, outputDirFlag string
	setWorkingDirFlag, setOutputDirFlag, setPlainFlag := false, false, false

	// cli arguments are optional
	if len(os.Args) >= 2 {
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]

			if arg == "--" {
				break
			} else if arg == "--plain" {
				if (i + 1) < len(os.Args) {
					val := os.Args[i+1]
					if val == "none" {
						EnablePlainOutput = false
						setPlainFlag = true
						i++
					}
				}

				if !setPlainFlag {
					EnablePlainOutput = true
				}
				setPlainFlag = true
			} else if arg == "-i" || arg == "--input-directory" {
				// Flag requires arg
				if (i + 1) >= len(os.Args) {
					errorExit(1, "missing value for flag '%v. See '%v --help' for more info.", arg, ProgramName)
				}

				workingDirFlag = os.Args[i+1]
				setWorkingDirFlag = true
				i++ // next arg is not a flag
			} else if arg == "-o" || arg == "--output-directory" {
				// Flag requires arg
				if (i + 1) >= len(os.Args) {
					errorExit(1, "missing value for flag '%v'. See '%v --help' for more info.", arg, ProgramName)
				}

				outputDirFlag = os.Args[i+1]
				setOutputDirFlag = true
				i++ // next arg is not a flag
			} else {
				errorExit(1, "unknown flag '%v'.", arg)
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
		errorExit(1, "output directory cannot be equal to the working/input directory ('%s' and '%s' reference the same path).", workingDirFlag, outputDirFlag)
	}

	// ---
	// Parse collection
	// ---

	collection, err := DecodeCollection(workingDirFlag)
	if err != nil {
		errorExit(1, "%w", err)
	}

	fmt.Printf("%#v\n", collection)
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
