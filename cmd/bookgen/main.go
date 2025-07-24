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
	Version           = "0.10.3"
	EnablePlainOutput = false
)

type Opts struct {
	Help                 bool      `long:"help" short:"h" desc:"Print help/usage information"`
	Version              bool      `long:"version" short:"v" desc:"Print program version"`
	PlainOutput          bool      `long:"plain" desc:"Remove terminal escape codes from printing into stdout/stderr"`
	NoNonEssentialOutput bool      `long:"no-non-essential-output" short:"q" desc:"Prevent printing non-error messages into stdout/stderr"`
	BuildCommand         BuildOpts `subcommand:"build" desc:"build source files"`
}

type BuildOpts struct {
	Minify          bool   `long:"minify" desc:"Minify output/distributable files"`
	InputDirectory  string `long:"input-directory" short:"i" desc:"Directory containing source files with a bookgen.yml"`
	OutputDirectory string `long:"output-directory" short:"o" desc:"Directory to output distributable files"`
}

func main() {
	// ---
	// Read CLI arguments
	// ---
	var opts Opts
	command, _, err := OptsParse(&opts, os.Args)
	if err != nil {
		errorExit(1, err.Error())
	}

	// ---
	// Set defaults
	// ---
	EnablePlainOutput = opts.PlainOutput

	// ---
	// Parse collection
	// ---
	if opts.Help {
		//optsPrintHelp(&opts, commands)
		fmt.Println("TODO: help")
		os.Exit(0)
	} else if opts.Version {
		fmt.Printf("bookgen version %v %v/%v\n", Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	} else if command == "build" {
		inputDirectory := opts.BuildCommand.InputDirectory
		outputDirectory := opts.BuildCommand.OutputDirectory
		enableMinify := opts.BuildCommand.Minify

		// Default output directory is relative to the input directory.
		if outputDirectory == "" {
			outputDirectory = filepath.Join(inputDirectory, "out")
		}

		// output dir cannot be == input dir
		if filepath.Clean(inputDirectory) == filepath.Clean(outputDirectory) {
			errorExit(1, "output directory cannot be equal to the working/input directory (`%s` and `%s` are the same).", inputDirectory, outputDirectory)
		}

		totalTimeStart := time.Now()
		decodeTimeStart := time.Now()

		collection, err := bookgen.DecodeCollection(inputDirectory)
		if err != nil {
			errorExit(1, err.Error())
		}

		decodeTimeElapsed := time.Since(decodeTimeStart)

		if !opts.NoNonEssentialOutput {
			fmt.Printf("Decoded (%v)\n", decodeTimeElapsed)
		}

		renderTimeStart := time.Now()

		if err := RenderCollectionToWebsite(&collection, inputDirectory, outputDirectory, enableMinify); err != nil {
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
		errorExit(1, "unrecognized command. See `%v --help` for more information", os.Args[0])
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
