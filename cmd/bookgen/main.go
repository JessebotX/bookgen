package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/JessebotX/bookgen"
)

type FlagType int

const (
	FlagString FlagType = iota
	FlagBool
)

var (
	Version                    = "0.8.1"
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
		Description: "The working/input directory that contains a bookgen.yml file.",
		Value:       "./",
	}
	outputDirFlag := Flag{
		Name:        "output-directory",
		Type:        FlagString,
		ShortName:   "o",
		Description: "The output directory where the distributable contents will be generated in.",
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
	minifyFlag := Flag{
		Name:        "minify",
		Type:        FlagBool,
		ShortName:   "",
		Description: "Minify output files.",
		Value:       "false",
	}

	positionalArgs, err := flagParse([]*Flag{
		&inputDirFlag,
		&outputDirFlag,
		&plainOutputFlag,
		&suppressNonEssentialOutputFlag,
		&minifyFlag,
	})
	if err != nil {
		errorExit(1, err.Error())
	}

	if len(positionalArgs) > 0 && !inputDirFlag.IsSet {
		inputDirFlag.Value = positionalArgs[0]
	}

	EnablePlainOutput, _ = strconv.ParseBool(plainOutputFlag.Value)
	SuppressNonEssentialOutput, _ = strconv.ParseBool(suppressNonEssentialOutputFlag.Value)
	enableMinify, _ := strconv.ParseBool(minifyFlag.Value)

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

	totalTimeStart := time.Now()

	decodeTimeStart := time.Now()

	collection, err := bookgen.DecodeCollection(inputDirFlag.Value)
	if err != nil {
		errorExit(1, err.Error())
	}

	decodeTimeElapsed := time.Since(decodeTimeStart)

	if !SuppressNonEssentialOutput {
		fmt.Printf("Decoded (%v)\n", decodeTimeElapsed)
	}

	renderTimeStart := time.Now()

	if err := RenderCollectionToWebsite(&collection, inputDirFlag.Value, outputDirFlag.Value, enableMinify); err != nil {
		errorExit(1, err.Error())
	}

	renderTimeElapsed := time.Since(renderTimeStart)

	if !SuppressNonEssentialOutput {
		fmt.Printf("Generated website (%v)\n", renderTimeElapsed)
	}

	totalTimeElapsed := time.Since(totalTimeStart)
	if !SuppressNonEssentialOutput {
		fmt.Printf(terminalPrintBold("Done")+" (%v)\n", totalTimeElapsed)
	}
}

func flagParse(flags []*Flag) ([]string, error) {
	positionalArgs := make([]string, 0)

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if arg == "--" {
			break
		}

		if arg == "--help" || arg == "-h" || arg == "-?" {
			printHelp(flags)

			os.Exit(0)
		}

		// TODO: support help and version flags here.
		if arg == "--version" || arg == "-v" || arg == "-V" {
			fmt.Printf("%v version %v %v/%v\n", ProgramName, Version, runtime.GOOS, runtime.GOARCH)

			os.Exit(0)
		}

		flagSet := false
		for _, f := range flags {
			if arg == "--"+f.Name || (f.ShortName != "" && arg == "-"+f.ShortName) {
				if f.Type == FlagString {
					if (i + 1) >= len(os.Args) {
						return positionalArgs, fmt.Errorf("missing value for `%v`. See `%v --help for more information`", arg, ProgramName)
					}

					f.IsSet = true
					flagSet = true
					f.Value = os.Args[i+1]
					i++
					break
				}

				if f.Type == FlagBool {
					f.IsSet = true
					flagSet = true
					f.Value = "true"
					break
				}
			}
		}

		if !flagSet {
			positionalArgs = append(positionalArgs, arg)
		}
	}

	return positionalArgs, nil
}

func printHelp(flags []*Flag) {
	indentSize := 4
	fmt.Println("USAGE")

	for range indentSize {
		fmt.Printf(" ")
	}

	fmt.Printf("%v [FLAGS...] [/path/to/input-directory]\n\n", ProgramName)
	fmt.Println("FLAGS")

	// Help flag
	for range indentSize {
		fmt.Printf(" ")
	}
	fmt.Println("-h, -?, --help")
	for range indentSize * 3 {
		fmt.Printf(" ")
	}
	fmt.Println("Get help information on how to use this program.")

	// Version flag
	for range indentSize {
		fmt.Printf(" ")
	}
	fmt.Println("-v, -V, --version")
	for range indentSize * 3 {
		fmt.Printf(" ")
	}
	fmt.Println("Get program version.")

	// Other defined flags
	for _, f := range flags {
		for range indentSize {
			fmt.Printf(" ")
		}

		if f.ShortName != "" {
			fmt.Printf("-%v, --%v", f.ShortName, f.Name)
		} else {
			fmt.Printf("    --%v", f.Name)
		}

		if f.Type == FlagString {
			fmt.Printf(" <string>")
		}

		fmt.Printf("\n")

		for range indentSize * 3 {
			fmt.Printf(" ")
		}

		fmt.Printf("%v (default: %v)\n", f.Description, f.Value)
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
