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
	SHOW_PROGRAM_SHEET_ID    string = "1nYxraiz2lu4YXmxFsSmhxpBNLy5yuz1tE3nybKznFgA" // Live Sheet
	SHOW_PROGRAM_SHEET_NAME  string = "Sheet1"
	SERVICE_ACCOUNT_KEY_FILE string = "impro-neuf-management-99d59b5d3102.json"

	POST_TYPE_FB     = "fb"
	POST_TYPE_SIO    = "sio"
	POST_TYPE_MEETUP = "meetup"
	POST_TYPE_INSTA  = "insta"
	POST_TYPE_STORY  = "story"
)

var POST_TYPES = []string{
	POST_TYPE_FB,
	POST_TYPE_INSTA,
	POST_TYPE_MEETUP,
	POST_TYPE_SIO,
	POST_TYPE_STORY,
}

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

func SaveScreenshot(tmpl *template.Template, show Show, tmplType string) {
	// Create output file
	fileName := show.Date.Format("2006-01-02") + " - " + show.Title + " - " + tmplType
	outputFilePath := "output/" + fileName + ".html"
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Execute the template
	err = tmpl.Execute(outputFile, show)
	if err != nil {
		panic(err)
	}

	// Prepare screenshot file path
	screenshotFile := "output/screenshots/" + fileName + ".jpg"
	path, _ := os.Getwd()
	fileUrl := "file://" + filepath.Join(path, outputFile.Name())
	log.Println("fileUrl:", fileUrl)

	// Create a context with a timeout to prevent hanging
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set image dimensions
	imageWidth := int64(1920)
	imageHeight := int64(1080)

	switch tmplType {
	case POST_TYPE_MEETUP, POST_TYPE_SIO:
		imageWidth = int64(1920)
		imageHeight = int64(1080)
	case POST_TYPE_INSTA:
		imageWidth = int64(1080)
		imageHeight = int64(1080)
	case POST_TYPE_STORY:
		imageWidth = int64(1080)
		imageHeight = int64(1920)
	case POST_TYPE_FB:
		imageWidth = int64(1920)
		imageHeight = int64(1004)
	}

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(imageWidth*2, imageHeight*2),
		chromedp.Navigate(fileUrl),
		// Wait until window.layoutAdjusted === true
		chromedp.ActionFunc(func(ctx context.Context) error {
			var isAdjusted bool
			for i := 0; i < 100; i++ { // Adjust the number of iterations as needed
				err := chromedp.Evaluate(`window.layoutAdjusted === true`, &isAdjusted).Do(ctx)
				if err != nil {
					return err
				}
				if isAdjusted {
					return nil
				}
				time.Sleep(50 * time.Millisecond) // Adjust the sleep duration as needed
			}
			return fmt.Errorf("timeout waiting for window.layoutAdjusted to be true")
		}),
		// Optional: Adjust zoom for higher resolution
		chromedp.ActionFunc(func(ctx context.Context) error {
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

func CreateIndex(shows []Show) {
	// Filter shows that have at least one team
	var validShows []Show
	today := TruncateToDate(time.Now())
	for _, show := range shows {
		if len(show.Teams) == 0 {
			continue
		}
		validShows = append(validShows, show)
	}

	// Prepare data for the template
	timestamp := time.Now().Unix()
	var showsData []ShowPageData
	for _, show := range validShows {
		dateStr := show.Date.Format("2006-01-02")
		showDate := TruncateToDate(show.Date)

		var types []ShowTypeData
		for _, tmplType := range POST_TYPES {
			imageFileName := fmt.Sprintf("%s - %s - %s.jpg?%d", dateStr, show.Title, tmplType, timestamp)
			types = append(types, ShowTypeData{
				Type:          tmplType,
				ImageFileName: imageFileName,
			})
		}

		showData := ShowPageData{
			DateStr: dateStr,
			Title:   show.Title,
			Types:   types,
			IsPast:  showDate.Before(today),
		}

		showsData = append(showsData, showData)
	}

	data := IndexPageData{
		Shows: showsData,
	}

	// Define the path to the template file
	templatePath := "index.tmpl"

	// Parse the template from the external file
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}

	// Create or open the index.html file
	indexFile, err := os.Create("output/screenshots/index.html")
	if err != nil {
		panic(err)
	}
	defer indexFile.Close()

	// Execute the template, writing the output to the indexFile
	err = t.Execute(indexFile, data)
	if err != nil {
		panic(err)
	}

	fmt.Println("index.html has been successfully created.")
}

// Helper function to truncate time to midnight
func TruncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func CreateShowPage(show Show) {
	dateStr := show.Date.Format("2006-01-02")
	fileName := "output/screenshots/" + dateStr + ".html"

	// Create or open the date-specific HTML file
	showFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer showFile.Close()

	// Prepare data for the template
	timestamp := time.Now().Unix()
	var types []ShowTypeData
	for _, tmplType := range POST_TYPES {
		imageFileName := fmt.Sprintf("%s - %s - %s.jpg?%d", dateStr, show.Title, tmplType, timestamp)
		types = append(types, ShowTypeData{
			Type:          tmplType,
			ImageFileName: imageFileName,
		})
	}

	data := ShowPageData{
		DateStr: dateStr,
		Title:   show.Title,
		Types:   types,
	}

	// Define the path to the template file
	templatePath := "show-page.tmpl"

	// Parse the template from the external file
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}

	// Execute the template, writing the output to the showFile
	err = t.Execute(showFile, data)
	if err != nil {
		panic(err)
	}
}

func GetFreeText(show Show) string {
	var hasEnglish, hasNorwegian bool

	// Check which languages are present in show.Languages
	for _, language := range show.Languages {
		switch language {
		case English:
			hasEnglish = true
		case Norwegian:
			hasNorwegian = true
		}
	}

	// Return text based on language presence
	if hasEnglish && hasNorwegian {
		return "FREE (students) / GRATIS (studenter)"
	} else if hasEnglish {
		return "FREE (students)"
	} else if hasNorwegian {
		return "GRATIS (studenter)"
	}

	return ""
}

func main() {
	// Create output directories if they don't exist
	if err := os.MkdirAll("output/screenshots", 0755); err != nil {
		log.Fatalf("Unable to create output directories: %v", err)
	}

	// Get the path to the Excel file
	xlsxFilePath := GetGoogleSheetsPath(SHOW_PROGRAM_SHEET_ID)
	showSchedule := ReadShowScheduleFromFile(xlsxFilePath, SHOW_PROGRAM_SHEET_NAME)

	funcMap := template.FuncMap{
		"GetTeamPhoto":   GetTeamPhoto,
		"formatMonth":    formatMonth,
		"GetShowEndTime": GetShowEndTime,
		"GetFreeText":    GetFreeText,
	}

	var shows []Show

	for _, show := range showSchedule {
		if show.Types[0] != ShowTypeRegular {
			continue
		}

		// Parse the template file
		tmplFb, err := template.New("regular-fb.tmpl").Funcs(funcMap).ParseFiles("regular-fb.tmpl")
		if err != nil {
			panic(err)
		}
		tmplInsta, err := template.New("regular-insta.tmpl").Funcs(funcMap).ParseFiles("regular-insta.tmpl")
		if err != nil {
			panic(err)
		}
		tmplSio, err := template.New("regular-sio-meetup.tmpl").Funcs(funcMap).ParseFiles("regular-sio-meetup.tmpl")
		if err != nil {
			panic(err)
		}
		tmplMeetup, err := template.New("regular-sio-meetup.tmpl").Funcs(funcMap).ParseFiles("regular-sio-meetup.tmpl")
		if err != nil {
			panic(err)
		}
		tmplStory, err := template.New("regular-story.tmpl").Funcs(funcMap).ParseFiles("regular-story.tmpl")
		if err != nil {
			panic(err)
		}
		//show.Teams = deduplicateStrings(show.Teams)

		// Save the show for index generation
		shows = append(shows, show)

		SaveScreenshot(tmplFb, show, POST_TYPE_FB)
		SaveScreenshot(tmplInsta, show, POST_TYPE_INSTA)
		SaveScreenshot(tmplSio, show, POST_TYPE_SIO)
		SaveScreenshot(tmplMeetup, show, POST_TYPE_MEETUP)
		SaveScreenshot(tmplStory, show, POST_TYPE_STORY)

		// Generate date-specific HTML file with links
		CreateShowPage(show)

	}

	// Generate the index.html
	CreateIndex(shows)
}
