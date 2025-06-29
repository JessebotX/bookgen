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

type Flag struct {
	Name        string
	Type        string
	ShortName   string
	Description string
	Value       string
	IsSet       bool
}

func flagParse(flags []*Flag) ([]string, error) {
	positionalArguments := make([]string, 0)
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		for _, f := range flags {
			if arg == "--"+f.Name || arg == "-"+f.ShortName {
				if f.Type == "string" {
					if (i + 1) >= len(os.Args) {
						return positionalArguments, fmt.Errorf("missing value for `%v`. See `%v --help for more information`", arg, ProgramName)
					}

					f.IsSet = true
					f.Value = os.Args[i+1]
					i++
					break
				}
			}

			positionalArguments = append(positionalArguments, arg)
		}
	}

	return positionalArguments, nil
}

func main() {
	// ---
	// Read CLI arguments
	// ---
	// Home-made cli argument parsing
	ProgramName = os.Args[0]

	inputDirFlag := Flag{
		Name:        "input-directory",
		Type:        "string",
		ShortName:   "i",
		Description: "The working/input directory.",
		Value:       "./",
	}
	outputDirFlag := Flag{
		Name:        "output-directory",
		Type:        "string",
		ShortName:   "o",
		Description: "The output directory.",
		Value:       "./out",
	}

	flags := []*Flag{&inputDirFlag, &outputDirFlag}
	if _, err := flagParse(flags); err != nil {
		errorExit(1, err.Error())
	}

	// make default output dir generate relative to working dir if
	// output dir is not provided
	if !outputDirFlag.IsSet {
		outputDirFlag.Value = filepath.Join(inputDirFlag.Value, "out")
	}

	fmt.Println(inputDirFlag.Value, outputDirFlag.Value)

	// output dir cannot be == input dir
	if filepath.Clean(inputDirFlag.Value) == filepath.Clean(outputDirFlag.Value) {
		errorExit(1, "output directory cannot be equal to the working/input directory (`%s` and `%s` reference the same path).", inputDirFlag.Value, outputDirFlag.Value)
	}

	// ---
	// Parse collection
	// ---

	timeStart := time.Now()

	collection, err := DecodeCollection(inputDirFlag.Value)
	if err != nil {
		errorExit(1, err.Error())
	}

	if err := RenderCollectionToWebsite(&collection, inputDirFlag.Value, outputDirFlag.Value); err != nil {
		errorExit(1, err.Error())
	}

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
