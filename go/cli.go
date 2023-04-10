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

func parseArguments() (string, string, string) {
	if len(os.Args) != 4 {
		printUsage()
		os.Exit(0)
	}

	if os.Args[1] != "image" && os.Args[1] != "video" {
		printError("type should be either \"image\" or \"video\"")
	}

	return os.Args[1], checkSourceFolder(os.Args[2]), createFolder(os.Args[3])
}

func printUsage() {
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf(`
%s

organize images into folders based on EXIF dates (or fallback to filename parsing)

%s %s %s %s
`,
		green("Image Sorter"),
		green("./sort-files"),
		yellow("TYPE"),
		yellow("SRC_FOLDER"),
		yellow("DEST_FOLDER"),
	)
}
