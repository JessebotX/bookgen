package main

import (
	"fmt"
)

type Opts struct {
	BuildOpts          BuildOpts `subcommand:"build" desc:"Compile source files into distributable output formats."`
	Version            bool      `long:"version" short:"v" desc:"Print application version."`
	NonEssentialOutput bool      `long:"non-essential-output" short:"q=no" desc:"Include non-essential messages (e.g. compilation states) when printing to terminal output."`
	PlainOutput        bool      `long:"plain" desc:"Strip terminal escape codes (e.g. colors, bold fonts) from terminal output." env:"TERM==dumb"`
}

type BuildOpts struct {
	InputDirectory   string `long:"input-directory" short:"i" desc:"Path to directory containing source files."`
	OutputDirectory  string `long:"output-directory" short:"o" desc:"Path to directory containing compiled output files/formats for distribution."`
	LayoutsDirectory string `long:"layouts-directory" desc:"Path to directory containing files that lay out output formats."`
	Minify           bool   `long:"minify" desc:"Minify output of supported file formats."`
	JSON             bool   `long:"json" desc:"Output distributable contents as JSON instead of static files."`
}

func main() {
	var opts Opts
	fmt.Println("Hello, world!")
	fmt.Printf("%#v\n", opts)
}
