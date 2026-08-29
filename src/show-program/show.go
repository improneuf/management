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

type ShowTypeData struct {
	Type          string
	ImageFileName string
}

type ShowPageData struct {
	DateStr string
	Title   string
	Teams   string
	Types   []ShowTypeData
	IsPast  bool
}

type IndexPageData struct {
	Shows []ShowPageData
}

func GetShowFromRow(row []string) Show {
	var teams []string
	var date time.Time
	var day string
	var showLanguages []ShowLanguage
	var showTypes []ShowType

	// Skip empty rows
	if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
		return Show{}
	}

	// Team 1 through Team 3 are columns F-H.
	for i := 5; i <= 7 && i < len(row); i++ {
		if team := strings.TrimSpace(row[i]); team != "" {
			teams = append(teams, team)
		}
	}
	if len(teams) == 0 {
		return Show{}
	}

	fmt.Printf("Processing row: %v\n", row)

	// Date is column A in DD.MM.YY format.
	dateStr := strings.TrimSpace(row[0])
	date, err := time.Parse("02.01.06", dateStr)
	if err != nil {
		log.Printf("Skipping row with invalid date %q: %v", dateStr, err)
		return Show{}
	}

	fmt.Printf("Parsed date: %v\n", date)

	// Get day of week
	if len(row) > 1 {
		day = strings.TrimSpace(row[1])
	}

	// Room is column D; column C is the booking time.
	venue := "Lillesalen, Chateau Neuf" // Default venue
	if len(row) > 3 && strings.TrimSpace(row[3]) != "" {
		venue = strings.TrimSpace(row[3])
	}

	fmt.Printf("Found teams: %v\n", teams)

	// Determine languages based on day of week
	if strings.Contains(strings.ToLower(day), "wednesday") {
		showLanguages = append(showLanguages, English)
	} else {
		showLanguages = append(showLanguages, Norwegian)
	}

	// If no show type was specified, default to regular show
	if len(showTypes) == 0 {
		showTypes = append(showTypes, ShowTypeRegular)
	}

	showTypes = removeDuplicateShowTypes(showTypes)
	showTitle, showSubtitle := getShowTitleAndSubtitle(showTypes[0], date)

	fmt.Printf("Generated title: %s, subtitle: %s\n", showTitle, showSubtitle)

	return Show{
		Date:         date,
		Day:          day,
		CrewSjefTeam: strings.TrimSpace(row[4]),
		Teams:        teams,
		Languages:    showLanguages,
		Types:        showTypes,
		Title:        showTitle,
		Subtitle:     showSubtitle,
		Venue:        venue,
		Price:        GetPriceFromShowType(showTypes[0]),
	}
}

// reads and prints the contents of the given Excel file.
func ReadShowScheduleFromFile(filePath string, sheetName string) []Show {
	shows := make([]Show, 0)

	fmt.Printf("Reading file: %s, sheet: %s\n", filePath, sheetName)

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	// Get the names of all the sheets
	sheets := f.GetSheetList()
	fmt.Printf("Found sheets: %v\n", sheets)

	// Iterate through each sheet and print its contents
	for _, sheet := range sheets {
		fmt.Printf("Processing sheet: %s\n", sheet)
		if sheet != sheetName {
			fmt.Printf("Skipping sheet %s (looking for %s)\n", sheet, sheetName)
			continue
		}

		// Get all rows
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Fatalf("Failed to get rows for sheet %s: %v", sheet, err)
		}

		fmt.Printf("Found %d rows in sheet %s\n", len(rows), sheet)

		// Skip header row
		for r, row := range rows {
			if r == 0 {
				fmt.Printf("Skipping header row: %v\n", row)
				continue
			}

			// Get the actual cell values
			var processedRow []string
			for c := range row {
				cellAddr := fmt.Sprintf("%c%d", 'A'+c, r+1)

				// Set date format for first column
				if c == 0 {
					style, err := f.NewStyle(&excelize.Style{
						NumFmt: 14, // 14 is the format code for "m/d/yyyy"
					})
					if err != nil {
						log.Printf("Warning: Failed to create style: %v", err)
					} else {
						err = f.SetCellStyle(sheet, cellAddr, cellAddr, style)
						if err != nil {
							log.Printf("Warning: Failed to set cell style: %v", err)
						}
					}
				}

				cell, err := f.GetCellValue(sheet, cellAddr)
				if err != nil {
					log.Printf("Warning: Failed to get cell value at %s: %v", cellAddr, err)
					cell = ""
				}

				processedRow = append(processedRow, cell)
			}

			// Print raw cell values
			fmt.Printf("Raw row %d: ", r)
			for c, cell := range processedRow {
				fmt.Printf("[%d:%q] ", c, cell)
			}
			fmt.Println()

			show := GetShowFromRow(processedRow)
			if show.Date.IsZero() {
				fmt.Printf("Skipping empty row: %v\n", processedRow)
				continue
			}
			shows = append(shows, show)
		}
	}
	fmt.Printf("Processed %d shows\n", len(shows))
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
	case "The Sound of Neuf":
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
