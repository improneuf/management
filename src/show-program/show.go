package main

import (
	"time"
)

type Show struct {
	Date          time.Time
	Day           string
	CrewSjefTeam  string
	Teams         []string
	ShowLanguages []ShowLanguage
	ShowTypes     []ShowType
}
