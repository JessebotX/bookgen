package main

import (
	"fmt"
	"path/filepath"

	"github.com/JessebotX/pub"
	pubhtml "github.com/JessebotX/pub/renderer/html"
)

const (
	LayoutsDirName = "_layouts"
	OutputDirName  = "_output"
)

type BuildCommand struct {
	InputDirectory   string  `name:"input-directory" type:"existingdir" default:"./" arg:"" help:"Directory containing structured source files that will be parsed into distributable output formats"`
	OutputDirectory  *string `name:"output-directory" short:"o" help:"Directory for distributable output formats. By default, directory is relative to the specified input directory"`
	LayoutsDirectory *string `name:"layouts-directory" short:"t" help:"Directory containing formatting instructions for distributable output formats. By default: directory is relative to the specified input directory"`
	Minify           bool    `name:"minify" help:"Optimize file sizes of distributable output formats"`
}

func (b BuildCommand) Run(ctx *Context) error {
	if ctx.Debug && ctx.NoNonEssentialMessages {
		fmt.Println("[NOTE] CLI flag \"debug\" and \"no-non-essential-messages\" are enabled at the same time.")
	}

	inputDir := b.InputDirectory
	var outputDir, layoutsDir string
	if b.OutputDirectory != nil {
		outputDir = *b.OutputDirectory
	} else {
		outputDir = filepath.Join(inputDir, OutputDirName)
	}

	if b.LayoutsDirectory != nil {
		layoutsDir = *b.LayoutsDirectory
	} else {
		layoutsDir = filepath.Join(inputDir, LayoutsDirName)
	}

	if ctx.Debug {
		fmt.Printf("[DEBUG] Input Directory:   %s\n", inputDir)
		fmt.Printf("[DEBUG] Layouts Directory: %s\n", layoutsDir)
		fmt.Printf("[DEBUG] Output Directory:  %s\n", outputDir)
	}

	if !ctx.NoNonEssentialMessages {
		fmt.Println("CREATING NEW BOOK...")
	}

	book, err := pub.NewBook(inputDir)
	if err != nil {
		return err
	}

	if !ctx.NoNonEssentialMessages {
		fmt.Println("Done CREATING NEW BOOK!")
	}

	if ctx.Debug {
		fmt.Println("[DEBUG] --- BEGIN BOOK STRUCTURE ---")

		fmt.Printf("[DEBUG] %#v\n", book)

		fmt.Println("[DEBUG] ---  END BOOK STRUCTURE  ---")
	}

	if !ctx.NoNonEssentialMessages {
		fmt.Println("GENERATING STATIC WEBSITE...")
	}

	if err := pubhtml.RenderBook(&book, inputDir, outputDir, layoutsDir, b.Minify); err != nil {
		return err
	}

	if !ctx.NoNonEssentialMessages {
		fmt.Println("Done GENERATING STATIC WEBSITE...")
	}

	return nil
}
