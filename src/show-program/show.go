package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type Show struct {
	Date          time.Time
	Day           string
	CrewSjefTeam  string
	Teams         []string
	ShowLanguages []ShowLanguage
	ShowTypes     []ShowType
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
	showTypes = removeDuplicateShowTypes(showTypes)
	return Show{
		Date:          date,
		Day:           day,
		CrewSjefTeam:  crewSjefTeam,
		Teams:         teams,
		ShowLanguages: showLanguages,
		ShowTypes:     showTypes,
	}
}

// ReadExcelFile reads and prints the contents of the given Excel file.
func ReadShowScheduleFromFile(filePath string) []Show {
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
		if sheet != SHOW_PROGRAM_SHEET_NAME {
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
