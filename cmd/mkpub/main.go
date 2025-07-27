package main

import (
	"fmt"
	"os"
)

type GlobalOpts struct {
	HelpOpts             HelpOpts    `subcommand:"help" desc:"Print help/usage information."`
	VersionOpts          VersionOpts `subcommand:"version" desc:"desc:"Print application version."`
	BuildOpts            BuildOpts   `subcommand:"build" desc:"Compile source files into distributable output formats."`
	InitOpts             InitOpts    `subcommand:"init" desc:"Generate directory structure."`
	Help                 bool        `long:"help" short:"h" desc:"Print help/usage information."`
	Version              bool        `long:"version" short:"v" desc:"Print application version."`
	NoNonEssentialOutput bool        `long:"no-non-essential-output" short:"q" desc:"Include non-essential messages (e.g. compilation states) when printing to terminal output."`
	PlainOutput          bool        `long:"plain" desc:"Strip terminal escape codes (e.g. colors, bold fonts, etc.) from terminal output." env:"TERM==dumb,NO_COLOR"`
}

type HelpOpts struct{}

type VersionOpts struct{}

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

const (
	TerminalClear      = "\033[0m"
	TerminalTextBold   = "1"
	TerminalTextRed    = "31"
	TerminalTextGreen  = "32"
	TerminalTextYellow = "33"
	TerminalTextWhite  = "37"
)

var (
	Program = ProgramInfo{
		Name:          "mkpub",
		UsageSynopsis: "mkpub <command> [flags...]",
		Description:   "mkpub transforms your Markdown-based source files into a distributable text publication.",
		Version:       "0.11.0",
	}
	Opts GlobalOpts
)

func main() {
	command, posArgs, err := OptsParse(&Opts, os.Args)
	if err != nil {
		errExit(1, err.Error())
	}

	if command == "help" || Opts.Help {
		if len(posArgs) > 1 {
			errExit(1, "command 'help' accepts at most 1 argument.")
		}

		if len(posArgs) == 1 && posArgs[0] != "build" && posArgs[0] != "init" {
			errExit(1, "command 'help' has no information for '%s'.", posArgs[0])
		}

		OptsWriteHelp(os.Stderr, &Opts, Program)
		os.Exit(0)
	} else if command == "version" || Opts.Version {
		if len(posArgs) > 0 {
			errExit(1, "command 'version' does not accept any arguments.")
		}

		fmt.Printf("%s v%s\n", Program.Name, Program.Version)

		os.Exit(0)
	} else {
		if len(posArgs) == 0 {
			errExit(1, "command/argument not found.")
		}

		errExit(1, "command/argument not recognized.")
	}
}

func errExit(exitCode int, format string, a ...any) {
	fmt.Fprintf(
		os.Stderr,
		terminalStyle(Program.Name+" error: ", TerminalTextBold+";"+TerminalTextRed)+format+"\n",
		a...,
	)
	os.Exit(exitCode)
}

func terminalStyle(s, code string) string {
	if Opts.PlainOutput {
		return s
	}
	return "\033[" + code + "m" + s + TerminalClear
}
