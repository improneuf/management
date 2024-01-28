package main

import (
	"strings"
)

type ShowLanguage string

const (
	English   ShowLanguage = "English"
	Norwegian ShowLanguage = "Norwegian"
)

func getShowLanguage(language string) ShowLanguage {
	language = strings.ToUpper(language)
	switch language {
	case "ENGLISH", "EN":
		return English
	case "NORWEGIAN", "NORSK", "NO", "NB":
		return Norwegian
	default:
		return English
	}
}
