package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type processInfo struct {
	path    string
	newPath string
	moved   bool
}

func iterateSourceFiles(folder string, dest string) {
	filter := getFileFilterExpression()
	exif.RegisterParsers(mknote.All...)

	foundFiles := 0
	channel := make(chan processInfo)

	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err == nil && filter.MatchString(info.Name()) {
			foundFiles++
			go handleImage(path, dest, channel)
			return nil
		}
		return err
	})

	if err != nil {
		printError(err.Error())
	}

	for i := 0; i < foundFiles; i++ {
		info := <-channel
		if info.moved {
			printSuccess("moved " + info.path + " to " + info.newPath)
		} else {
			printWarning("not moved " + info.path)
		}
	}
}

func handleImage(name string, destination string, channel chan processInfo) {
	creationTime, err := readDateFromExif(name)

	if creationTime == nil || err != nil {
		// no exif or no date in tags, try to read from filename
		creationTime, err = readDateFromName(name)
	}

	var folder string

	if creationTime != nil {
		month := strconv.Itoa(int(creationTime.Month()))
		if len(month) == 1 {
			month = "0" + month
		}

		folder = filepath.Join(destination, strconv.Itoa(creationTime.Year()), month)
		folder = createFolder(folder)

		err = moveFile(name, folder, *creationTime)

		if err != nil {
			creationTime = nil
		}
	}

	channel <- processInfo{path: name, moved: creationTime != nil, newPath: folder}
}

func moveFile(filename string, folder string, t time.Time) error {
	target := filepath.Join(folder, filepath.Base(filename))

	err := os.Rename(filename, target)
	if err != nil {
		return err
	}

	return os.Chtimes(target, t, t)
}

func checkSourceFolder(name string) string {
	_, srcError := os.Stat(name)
	if srcError != nil {
		printError("source folder doesnt exists")
	}

	return name
}

func createFolder(name string) string {
	dstError := os.MkdirAll(name, 0755)

	if dstError != nil {
		printError("destination folder could not be created: " + name)
	}

	return name
}

func readDateFromExif(name string) (*time.Time, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	decoded, err := exif.Decode(file)
	if err != nil || decoded == nil {
		return nil, err
	}
	creationTime, err := decoded.DateTime()

	if err != nil {
		return nil, err
	}

	return &creationTime, nil
}

func readDateFromName(name string) (*time.Time, error) {
	patterns := getDateFromFilenameExpressions()
	name = filepath.Base(name)

	for _, p := range patterns {
		if p.MatchString(name) {
			name = p.FindString(name)
			name = strings.Replace(name, "_", "-", 3)
			t, err := time.Parse("YYYY-MM-DD", name)
			return &t, err
		}
	}

	return nil, nil
}
