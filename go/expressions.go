package main

import (
	"regexp"
)

var expressions []*regexp.Regexp

var fileFilter *regexp.Regexp

func getFileFilterExpression() *regexp.Regexp {
	if fileFilter == nil {
		f, err := regexp.Compile("^.+\\.(jpeg|jpg|JPG|JPEG)$")

		if err != nil {
			printError(err.Error())
		}
		fileFilter = f
	}

	return fileFilter
}

func getDateFromFilenameExpressions() []*regexp.Regexp {
	if len(expressions) == 0 {
		patterns := []string{
			"[_|-]([\\d]{8})[_|-]",                           //matches e.g. IMG_20221030-foo.jpg
			"[_|-]([\\d]{4}[_|-][\\d]{2}[_|-][\\d]{2})[_|-]", //matches e.g. IMG-2022-10-30_bar.jpg
		}

		for _, pattern := range patterns {
			f, err := regexp.Compile(pattern)

			if err != nil {
				printError(err.Error())
			}
			expressions = append(expressions, f)
		}
	}

	return expressions
}
