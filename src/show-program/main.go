package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
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
	// bookingXlsxFilePath := GetGoogleSheetsPath(SHOW_SCHEDULE_SHEET_ID)

	// bookings := ReadBookingsFromFile(bookingXlsxFilePath, SHOW_SCHEDULE_SHEET_NAME)

	showProgramXlsxFilePath := GetGoogleSheetsPath(SHOW_PROGRAM_SHEET_ID)
	showSchedule := ReadShowScheduleFromFile(showProgramXlsxFilePath, SHOW_PROGRAM_SHEET_NAME)

	funcMap := template.FuncMap{
		"GetTeamPhoto":   GetTeamPhoto,
		"formatMonth":    formatMonth,
		"GetShowEndTime": GetShowEndTime,
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

		// screenshot the html file
		screenshotFile := "output/screenshots/" + show.Date.Format("2006-01-02") + " - " + show.Title + ".jpg"
		path, _ := os.Getwd()
		fileUrl := "file://" + filepath.Join(path, outputFile.Name())
		log.Println("fileUrl:", fileUrl)

		// Create a context
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		// Capture screenshot of an entire webpage in JPEG format
		var buf []byte
		if err := chromedp.Run(ctx,
			chromedp.EmulateViewport(1920, 1004),
			chromedp.Navigate(fileUrl),
			chromedp.ActionFunc(func(ctx context.Context) error {
				// Set the zoom level by scaling the CSS
				return chromedp.Evaluate(`document.body.style.zoom = "2"`, nil).Do(ctx)
			}),
			chromedp.FullScreenshot(&buf, 100),
		); err != nil {
			log.Fatal(err)
		}

		// Save the screenshot to a file
		if err := os.WriteFile(screenshotFile, buf, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
