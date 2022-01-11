package main

import (
	"log"
	"regexp"
)

var expressions []*regexp.Regexp

var fileFilter *regexp.Regexp

func getFileFilterExpression() *regexp.Regexp {
	if fileFilter == nil {
		f, err := regexp.Compile("^.+\\.(jpeg|jpg|JPG|JPEG)$")

		if err != nil {
			log.Fatal(err)
		}
		fileFilter = f
	}

	return fileFilter
}

func getDateFromFilenameExpressions() []*regexp.Regexp {
	if len(expressions) == 0 {
		patterns := []string{
			"([\\d]{8})[_|-]",
			"^IMG-([\\d]{8})[_|-]",
			"^IMG_([\\d]{8})[_|-]",
			"^IMG_([\\d]{4}_[\\d]{2}_[\\d]{2})[_|-]",
		}

		for _, pattern := range patterns {
			f, err := regexp.Compile(pattern)

			if err != nil {
				log.Fatal(err)
			}
			expressions = append(expressions, f)
		}
	}

	return expressions
}
