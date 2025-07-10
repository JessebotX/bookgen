package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/JessebotX/bookgen"
)

var (
	Version           = "0.9.1"
	EnablePlainOutput = false
)

type FlagOpts struct {
	Help                 bool   `long:"help" short:"h" desc:"Print help/usage information"`
	Version              bool   `long:"version" short:"v" desc:"Print program version"`
	InputDirectory       string `long:"input-directory" short:"i" desc:"Directory containing source files with a bookgen.yml"`
	OutputDirectory      string `long:"output-directory" short:"o" desc:"Directory to output distributable files"`
	Minify               bool   `long:"minify" desc:"Minify output/distributable files"`
	PlainOutput          bool   `long:"plain" desc:"Remove terminal escape codes from printing into stdout/stderr"`
	NoNonEssentialOutput bool   `long:"no-non-essential-output" short:"q" desc:"Prevent printing non-error messages into stdout/stderr"`
}

func main() {
	// ---
	// Read CLI arguments
	// ---
	// Home-made cli argument parsing
	helpCommand := Command{
		Name:        "help",
		Description: "Print help/usage information",
	}
	versionCommand := Command{
		Name:        "version",
		Description: "Print program version",
	}
	buildCommand := Command{
		Name:        "build",
		Description: "Build source files for distribution",
	}

	commands := []Command{buildCommand, helpCommand, versionCommand}
	var opts FlagOpts

	positionalArgs, err := optsParse(os.Args, &opts)
	if err != nil {
		errorExit(1, err.Error())
	}

	// set default command
	if len(positionalArgs) == 0 {
		positionalArgs = append(positionalArgs, buildCommand.Name)
	}

	// help and version commands/flags
	if opts.Help || (len(positionalArgs) > 0 && positionalArgs[0] == helpCommand.Name) {
		optsPrintHelp(&opts, commands)

		os.Exit(0)
	}

	if opts.Version || (len(positionalArgs) > 0 && positionalArgs[0] == versionCommand.Name) {
		fmt.Printf("bookgen version %v %v/%v\n", Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// set flags
	EnablePlainOutput = opts.PlainOutput

	// Default output directory is relative to the input directory.
	if opts.OutputDirectory == "" {
		opts.OutputDirectory = filepath.Join(opts.InputDirectory, "out")
	}

	// output dir cannot be == input dir
	if filepath.Clean(opts.InputDirectory) == filepath.Clean(opts.OutputDirectory) {
		errorExit(1, "output directory cannot be equal to the working/input directory (`%s` and `%s` are the same).", opts.InputDirectory, opts.OutputDirectory)
	}

	// ---
	// Parse collection
	// ---

	if positionalArgs[0] == buildCommand.Name {
		totalTimeStart := time.Now()

		decodeTimeStart := time.Now()

		collection, err := bookgen.DecodeCollection(opts.InputDirectory)
		if err != nil {
			errorExit(1, err.Error())
		}

		decodeTimeElapsed := time.Since(decodeTimeStart)

		if !opts.NoNonEssentialOutput {
			fmt.Printf("Decoded (%v)\n", decodeTimeElapsed)
		}

		renderTimeStart := time.Now()

		if err := RenderCollectionToWebsite(&collection, opts.InputDirectory, opts.OutputDirectory, opts.Minify); err != nil {
			errorExit(1, err.Error())
		}

		renderTimeElapsed := time.Since(renderTimeStart)

		if !opts.NoNonEssentialOutput {
			fmt.Printf("Generated website (%v)\n", renderTimeElapsed)
		}

		totalTimeElapsed := time.Since(totalTimeStart)
		if !opts.NoNonEssentialOutput {
			fmt.Printf(terminalPrintBold("Done")+" (%v)\n", totalTimeElapsed)
		}
	} else {
		errorExit(1, "unrecognized command `%v`. See `%v --help` for more information", positionalArgs[0], os.Args[0])
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
