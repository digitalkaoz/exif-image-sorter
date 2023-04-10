package main

func main() {
	filetype, src, dst := parseArguments()
	printSuccess("Reading Source Folder " + src)
	printSuccess("moving to Folder " + dst)
	iterateSourceFiles(filetype, src, dst)
}
