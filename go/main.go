package main

func main() {
	src, dst := parseArguments()
	printSuccess("Reading Source Folder " + src)
	printSuccess("moving to Folder " + dst)
	iterateSourceFiles(src, dst)
}
