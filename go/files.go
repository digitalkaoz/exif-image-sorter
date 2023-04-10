package main

import (
	"errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/tajtiattila/metadata"
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

func iterateSourceFiles(filetype string, folder string, dest string) {
	filter := getFileFilterExpression(filetype)
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

func isImage(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".jpg") || strings.HasSuffix(strings.ToLower(name), ".jpeg")
}

func handleImage(name string, destination string, channel chan processInfo) {
	var creationTime *time.Time
	var folder string
	var err error

	if isImage(name) {
		// fetch image creation date from exif data
		creationTime, err = readDateFromExif(name)

		if creationTime == nil || err != nil {
			// no exif or no date in tags, try to read from filename
			creationTime, err = readDateFromName(name)
		}
	} else {
		// fetch video creation date from metadata
		creationTimeVideo, _ := readDateFromMetadata(name)
		creationTimeFile, _ := readDateFromName(name)
		if creationTimeVideo == nil && creationTimeFile != nil {
			// only date in filename detected
			creationTime = creationTimeFile
		} else if creationTimeVideo != nil && creationTimeFile == nil {
			// only date in metadata detected
			creationTime = creationTimeVideo
		} else if creationTimeVideo != nil && creationTimeFile != nil && creationTimeVideo.Format(time.DateOnly) != creationTimeFile.Format(time.DateOnly) {
			// detected both, but they are different, so take the one from the filename
			creationTime = creationTimeFile
		} else {
			// both are set an equal
			creationTime = creationTimeFile
		}
		// is there an else?
	}

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

func readDateFromMetadata(name string) (*time.Time, error) {
	fd, err := os.Open(name)
	defer fd.Close()

	if err != nil {
		return nil, err
	}

	md, err := metadata.Parse(fd)
	if err != nil {
		return nil, err
	}
	if _, ok := md.Attr["DateTimeCreated"]; ok {
		zeroTime := time.Unix(0, 0)
		if md.DateTimeCreated.Time.Before(zeroTime) {
			// unreasonable date detected
			return nil, errors.New("no metadata extractable")
		}
		creationTime := md.DateTimeCreated.Time
		return &creationTime, nil
	}
	return nil, errors.New("no metadata extractable")
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
			name = strings.Replace(strings.Replace(name, "_", "", 3), "-", "", 3)
			t, err := time.Parse("20060102", name)
			return &t, err
		}
	}

	return nil, nil
}
