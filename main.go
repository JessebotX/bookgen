package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FlagType int

const (
	Version             = "0.1"
	FlagString FlagType = iota
	FlagBool
)

var (
	SuppressNonEssentialOutput = false
	EnablePlainOutput          = false
	ProgramName                = "bookgen"
)

type Flag struct {
	Name        string
	Type        FlagType
	ShortName   string
	Description string
	Value       string
	IsSet       bool
}

func main() {
	// ---
	// Read CLI arguments
	// ---
	// Home-made cli argument parsing
	ProgramName = os.Args[0]

	inputDirFlag := Flag{
		Name:        "input-directory",
		Type:        FlagString,
		ShortName:   "i",
		Description: "The working/input directory.",
		Value:       "./",
	}
	outputDirFlag := Flag{
		Name:        "output-directory",
		Type:        FlagString,
		ShortName:   "o",
		Description: "The output directory.",
		Value:       "./out",
	}
	plainOutputFlag := Flag{
		Name:        "plain",
		Type:        FlagBool,
		ShortName:   "",
		Description: "Strip text styling/terminal escape codes from terminal output.",
		Value:       "false",
	}
	suppressNonEssentialOutputFlag := Flag{
		Name:        "suppress-non-essential-output",
		Type:        FlagBool,
		ShortName:   "q",
		Description: "Suppress non-essential terminal output.",
		Value:       "false",
	}

	flags := []*Flag{
		&inputDirFlag,
		&outputDirFlag,
		&plainOutputFlag,
		&suppressNonEssentialOutputFlag,
	}
	positionalArgs, err := flagParse(flags)
	if err != nil {
		errorExit(1, err.Error())
	}

	if len(positionalArgs) > 0 && !inputDirFlag.IsSet {
		inputDirFlag.Value = positionalArgs[0]
	}

	EnablePlainOutput, _ = strconv.ParseBool(plainOutputFlag.Value)
	SuppressNonEssentialOutput, _ = strconv.ParseBool(suppressNonEssentialOutputFlag.Value)

	// make default output dir generate relative to working dir if
	// output dir is not provided
	if !outputDirFlag.IsSet {
		outputDirFlag.Value = filepath.Join(inputDirFlag.Value, "out")
	}

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

func flagParse(flags []*Flag) ([]string, error) {
	positionalArgs := make([]string, 0)

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		for _, f := range flags {
			if arg == "--"+f.Name || (f.ShortName != "" && arg == "-"+f.ShortName) {
				if f.Type == FlagString {
					if (i + 1) >= len(os.Args) {
						return positionalArgs, fmt.Errorf("missing value for `%v`. See `%v --help for more information`", arg, ProgramName)
					}

					f.IsSet = true
					f.Value = os.Args[i+1]
					i++
					break
				}

				if f.Type == FlagBool {
					f.IsSet = true
					f.Value = "true"
					break
				}
			}

			positionalArgs = append(positionalArgs, arg)
		}
	}

	return positionalArgs, nil
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
