package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

const (
	//SHOW_PROGRAM_SHEET_ID string = "1ejEDxQJIwQ1ougcpWIKTqauT-05PDVT1" // Test Sheet
	SHOW_PROGRAM_SHEET_ID    string = "167cJAqP9fON3ExyLnJLFaJ0MHdu5K--z" // Live Sheet
	SHOW_PROGRAM_SHEET_NAME  string = "ShowProgram"
	SHOW_SCHEDULE_SHEET_ID   string = "15cDopxkZDbFwIcIU5tuqAUCM4U0GXf7O"
	SHOW_SCHEDULE_SHEET_NAME string = "Yesplan"
	SERVICE_ACCOUNT_KEY_FILE string = "impro-neuf-management-99d59b5d3102.json"
)

// GetLocalFileModifiedDate returns the last modified date of the file at the given filePath.
func GetLocalFileModifiedTime(filePath string) (time.Time, error) {
	// Get file information
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// Return a zero time and the error if there's an issue accessing the file
		return time.Time{}, err
	}

	modifiedTime := fileInfo.ModTime()

	fmt.Println("Parsed Local File Modified Time:", modifiedTime)

	return modifiedTime, nil
}

func main() {
	showProgramXlsxFilePath := GetGoogleSheetsPath(SHOW_PROGRAM_SHEET_ID)
	showSchedule := ReadShowScheduleFromFile(showProgramXlsxFilePath)

	funcMap := template.FuncMap{
		"GetTeamPhoto": GetTeamPhoto,
		"formatMonth":  formatMonth,
	}

	for _, show := range showSchedule {
		if show.Types[0] != ShowTypeRegular {
			continue
		}

		// Parse the template file
		tmpl, err := template.New("regular.tmpl").Funcs(funcMap).ParseFiles("regular.tmpl")
		if err != nil {
			panic(err)
		}

		show.Teams = deduplicateStrings(show.Teams)

		// Create output file
		outputFile, err := os.Create("output/" + show.Date.Format("2006-01-02") + " - " + show.Title + ".html")
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		// Execute the template
		err = tmpl.Execute(outputFile, show)
		if err != nil {
			panic(err)
		}
	}
}
