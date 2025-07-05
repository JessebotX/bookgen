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

type Opts struct {
	Help                 bool   `long:"help" short:"h" desc:"print help/usage information"`
	Version              bool   `long:"version" short:"v" desc:"print program version"`
	InputDirectory       string `long:"input-directory" short:"i" desc:"directory containing source files with a bookgen.yml"`
	OutputDirectory      string `long:"output-directory" short:"o" desc:"directory to output distributable files"`
	PlainOutput          bool   `long:"plain" desc:"remove terminal escape codes from printing into stdout/stderr"`
	NoNonEssentialOutput bool   `long:"no-non-essential-output" short:"q" desc:"prevent printing non-error messages into stdout/stderr"`
	Minify               bool   `long:"minify" desc:"Minify output/distributable files"`
}

func main() {
	// ---
	// Read CLI arguments
	// ---
	// Home-made cli argument parsing
	var opts Opts

	_, err := flagParse(os.Args, &opts)
	if err != nil {
		errorExit(1, err.Error())
	}

	if opts.Help {
		flagPrintHelp(&opts)

		os.Exit(0)
	}

	if opts.Version {
		fmt.Printf("bookgen version %v %v/%v\n", Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

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
