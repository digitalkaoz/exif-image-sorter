package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
)

func printError(msg string) {
	red := color.New(color.FgHiRed).SprintFunc()
	log.Fatal(red(msg))
}

func printSuccess(msg string) {
	green := color.New(color.FgGreen).SprintFunc()
	log.Print(green(msg))
}

func printWarning(msg string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	log.Print(yellow(msg))
}

func parseArguments() (string, string) {
	if len(os.Args) != 3 {
		printUsage()
		os.Exit(0)
	}

	return checkSourceFolder(os.Args[1]), createFolder(os.Args[2])
}

func printUsage() {
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf(`
%s

organize images into folders based on EXIF dates (or fallback to filename parsing)

%s %s %s
`,
		green("Image Sorter"),
		green("./sort-files"),
		yellow("SRC_FOLDER"),
		yellow("DEST_FOLDER"),
	)
}
