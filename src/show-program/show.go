package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type Show struct {
	Date         time.Time
	Day          string
	CrewSjefTeam string
	Teams        []string
	Languages    []ShowLanguage
	Types        []ShowType
	Title        string
	Subtitle     string
	Venue        string
	Price        Price
}

func GetShowFromRow(row []string) Show {
	var teams []string
	var date time.Time
	var day string
	var crewSjefTeam string
	var showLanguages []ShowLanguage
	var showTypes []ShowType
	var err error

	for c, cell := range row {
		if c > 9 {
			break
		}
		cellTrimmed := strings.TrimSpace(cell)
		fmt.Printf("|%s|", cellTrimmed)
		switch c {
		case 0:
			date, err = time.Parse("2 Jan 2006", cellTrimmed)
			if err != nil {
				date2, err2 := time.Parse("2-Jan-06", cellTrimmed)
				if err2 != nil {
					log.Fatalf("Failed to parse date: %v", cellTrimmed)
				}
				date = date2
			}
		case 1:
			day = cellTrimmed
		case 2:
			crewSjefTeam = cellTrimmed
		case 4, 5, 6, 7, 8:
			if cellTrimmed != "" {
				teams = append(teams, cellTrimmed)
				showTypes = append(showTypes, getShowType(cellTrimmed))
			}
		case 9:
			showLanguagesStr := strings.Split(cellTrimmed, "/")
			for _, languageStr := range showLanguagesStr {
				languageStr = strings.TrimSpace(languageStr)
				showLanguages = append(showLanguages, getShowLanguage(languageStr))
			}
		}
	}

	// special handling for Problemfixers
	for i, team := range teams {
		if team == "Problemfikserne/Problemfixers" {
			found := false
			for _, lang := range showLanguages {
				if lang == Norwegian {
					found = true
					break
				}
			}
			if found {
				teams[i] = "Problemfikserne" // Norwegian is in showLanguages
			} else {
				teams[i] = "The Problem Fixers" // Norwegian is not in showLanguages
			}
		}
	}

	showTypes = removeDuplicateShowTypes(showTypes)
	showTitle, showSubtitle := getShowTitleAndSubtitle(showTypes[0], date)
	return Show{
		Date:         date,
		Day:          day,
		CrewSjefTeam: crewSjefTeam,
		Teams:        teams,
		Languages:    showLanguages,
		Types:        showTypes,
		Title:        showTitle,
		Subtitle:     showSubtitle,
		Venue:        "Lillesalen, Chateau Neuf",
		Price:        GetPriceFromShowType(showTypes[0]),
	}
}

// reads and prints the contents of the given Excel file.
func ReadShowScheduleFromFile(filePath string, sheetName string) []Show {
	shows := make([]Show, 0)

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	// Get the names of all the sheets
	sheets := f.GetSheetList()

	// Iterate through each sheet and print its contents
	for _, sheet := range sheets {
		fmt.Printf("Sheet: %s\n", sheet)
		if sheet != sheetName {
			continue
		}
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Fatalf("Failed to get rows for sheet %s: %v", sheet, err)
		}

		for r, row := range rows {
			if r < 10 {
				continue
			}
			shows = append(shows, GetShowFromRow(row))
		}
	}
	return shows
}

func deduplicateStrings(slice []string) []string {
	seen := make(map[string]struct{})
	unique := []string{}

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			unique = append(unique, item)
		}
	}

	return unique
}

func formatMonth(date time.Time) string {
	month := date.Month().String() // Get the full month name
	if len(month) >= 8 {
		return date.Format("Jan")
	}
	return date.Format("January")
}

func GetTeamShowDuration(teamName string) int {
	switch teamName {
	case "The sound of Neuf":
		return 40
	default:
		return 20
	}
}

func GetShowEndTime(startTime string, show Show) (string, error) {
	// Parse the start time using the same format as the input
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return "", err
	}

	// Calculate duration based on teams
	var duration int
	// iterate through the teams and add the extra time
	for _, team := range show.Teams {
		duration += GetTeamShowDuration(team)
	}

	var extraTime int
	switch duration / 20 {
	case 1:
		extraTime = 10
	case 2:
		extraTime = 5
	case 3:
		extraTime = 15
	case 4:
		extraTime = 25
	case 5:
		extraTime = 25
	}

	duration += extraTime

	// Add the duration to the start time to get the end time
	endTime := start.Add(time.Duration(duration) * time.Minute)

	// Format the end time as a string in the same format as the input
	return endTime.Format("15:04"), nil
}
