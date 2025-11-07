package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/goccy/go-yaml"
)

type Context struct {
	Plain                  bool
	NoNonEssentialMessages bool
	Debug                  bool
}

var CLI struct {
	Book                   BookCommand    `cmd:"" help:"Create/manage a book project"`
	Build                  BuildCommand   `cmd:"" default:"withargs" help:"Build structured source files"`
	Version                VersionCommand `cmd:"" help:"Print program version"`
	Plain                  bool           `name:"plain" env:"NO_COLOR" help:"Disable escape codes such as colors and font styling from being printed to terminal output"`
	NoNonEssentialMessages bool           `name:"no-non-essential-messages" short:"q" help:"Disable non-error and non-warning messages from being printed to terminal output"`
	Debug                  bool           `name:"debug" help:"Print extra information [such as all config structures and contents] to terminal output for debugging purposes"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if err := ctx.Run(&Context{
		Debug:                  CLI.Debug,
		Plain:                  CLI.Plain,
		NoNonEssentialMessages: CLI.NoNonEssentialMessages,
	}); err != nil {
		ctx.FatalIfErrorf(fmt.Errorf("%s", yaml.FormatError(err, CLI.Plain, true)))
	}
}
