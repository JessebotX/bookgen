package main

import (
	"fmt"
	"os"
)

type Opts struct {
	BuildOpts            BuildOpts `subcommand:"build" desc:"Compile source files into distributable output formats."`
	InitOpts             InitOpts  `subcommand:"init" desc:"Generate directory structure"`
	Version              bool      `long:"version" short:"v" desc:"Print application version."`
	NoNonEssentialOutput bool      `long:"no-non-essential-output" short:"q" desc:"Include non-essential messages (e.g. compilation states) when printing to terminal output."`
	PlainOutput          bool      `long:"plain" desc:"Strip terminal escape codes (e.g. colors, bold fonts) from terminal output." env:"TERM==dumb,NO_COLOR==1"`
}

type BuildOpts struct {
	InputDirectory   string `long:"input-directory" short:"i" desc:"Path to directory containing source files."`
	OutputDirectory  string `long:"output-directory" short:"o" desc:"Path to directory containing compiled output files/formats for distribution."`
	LayoutsDirectory string `long:"layouts-directory" desc:"Path to directory containing files that lay out output formats."`
	Minify           bool   `long:"minify" desc:"Minify output of supported file formats."`
	JSON             bool   `long:"json" desc:"Output distributable contents as JSON instead of static files."`
}

type InitOpts struct {
	DryRun bool `long:"dry-run" desc:"Show changes that will happen without actually performing the actions."`
}

func main() {
	var opts Opts
	command, posArgs, err := OptsParse(&opts, os.Args)
	if err != nil {
		return
	}

	fmt.Println("command:", command)
	fmt.Println("posArgs:", posArgs)
	fmt.Printf("%#v\n", opts)
}
