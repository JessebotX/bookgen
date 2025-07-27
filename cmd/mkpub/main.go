package main

import (
	"fmt"
	"os"
)

type GlobalOpts struct {
	HelpOpts             HelpOpts    `command:"help" desc:"Print help/usage information."`
	VersionOpts          VersionOpts `command:"version" desc:"Print application version."`
	BuildOpts            BuildOpts   `command:"build" desc:"Compile source files into distributable output formats."`
	InitOpts             InitOpts    `command:"init" desc:"Generate directory structure."`
	Help                 bool        `long:"help" short:"h" desc:"Print help/usage information."`
	Version              bool        `long:"version" short:"v" desc:"Print application version."`
	NoNonEssentialOutput bool        `long:"no-non-essential-output" short:"q" desc:"Include non-essential messages (e.g. compilation states) when printing to terminal output."`
	PlainOutput          bool        `long:"plain" desc:"Strip terminal escape codes (e.g. colors, bold fonts, etc.) from terminal output." env:"TERM==dumb,NO_COLOR"`
	TestInt              int         `long:"test-int" desc:"Test accepting int64 CLI"`
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
	TerminalClear       = "\033[0m"
	TerminalTextBold    = "1"
	TerminalTextRed     = "31"
	TerminalTextGreen   = "32"
	TerminalTextYellow  = "33"
	TerminalTextBlue    = "34"
	TerminalTextMagenta = "35"
	TerminalTextCyan    = "36"
	TerminalTextWhite   = "37"
)

var (
	Opts    GlobalOpts
	Program = ProgramInfo{
		Name:          "mkpub",
		UsageSynopsis: "mkpub <command> [flags...]",
		Description:   "mkpub transforms your Markdown-based source files into a distributable text publication.",
		Version:       "0.11.0",
	}
	GlobalHelpExamples = []HelpExample{
		HelpExample{
			Usage:       Program.Name + " init",
			Description: "Initialize default project structure for source files.",
		},
		HelpExample{
			Usage:       Program.Name,
			Description: "Convert sources in current directory into distributable output formats. Shorthand for 'build' command.",
		},
		HelpExample{
			Usage:       Program.Name + " build --minify",
			Description: "Convert sources in current directory to *minified* distributable output formats (minification support depends on format).",
		},
		HelpExample{
			Usage:       Program.Name + " build -i /path/to/source/dir -o /path/to/output/dir",
			Description: "Convert sources found in /path/to/source/dir into distributable output formats generated at /path/to/output/dir.",
		},
		HelpExample{
			Usage:       Program.Name + " help build",
			Description: "Print help/usage information for subcommand 'build'.",
		},
	}
)

func main() {
	command, posArgs, err := OptsParse(&Opts, os.Args)
	if err != nil {
		errExit(1, err.Error())
	}

	if command == "help" || Opts.Help {
		if len(posArgs) == 1 {
			if err := OptsWriteHelp(os.Stderr, &Opts, Program.Name, posArgs[0], GlobalHelpExamples...); err != nil {
				errExit(1, err.Error())
			}
		} else {
			OptsWriteHelp(os.Stderr, &Opts, Program.Name, "", GlobalHelpExamples...)
		}
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
