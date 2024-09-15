package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	//SHOW_PROGRAM_SHEET_ID string = "1ejEDxQJIwQ1ougcpWIKTqauT-05PDVT1" // Test Sheet
	SHOW_PROGRAM_SHEET_ID    string = "1BYucz1R4IoH5whYe4goRbk_kO8LosrZ2" // Live Sheet
	SHOW_PROGRAM_SHEET_NAME  string = "ShowProgram"
	SHOW_SCHEDULE_SHEET_ID   string = "15cDopxkZDbFwIcIU5tuqAUCM4U0GXf7O"
	SHOW_SCHEDULE_SHEET_NAME string = "Yesplan"
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
	outputFile, err := os.Create("output/" + fileName + ".html")
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
	screenshotFile := "output/screenshots/" + fileName + ".jpg"
	path, _ := os.Getwd()
	fileUrl := "file://" + filepath.Join(path, outputFile.Name())
	log.Println("fileUrl:", fileUrl)

	// Create a context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Capture screenshot of an entire webpage in JPEG format
	imageWidth := int64(1920)
	imageHeight := int64(1004)

	if tmplType == "meetup" || tmplType == "sio" {
		imageHeight = int64(1080)
	}
	if tmplType == "insta" {
		imageWidth = int64(2160)
		imageHeight = int64(2160)
	}

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(imageWidth, imageHeight),
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

func CreateIndex(shows []Show) {
	// Create or open the index.html file
	indexFile, err := os.Create("output/screenshots/index.html")
	if err != nil {
		panic(err)
	}
	defer indexFile.Close()

	// Start writing HTML content
	indexFile.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	indexFile.WriteString("<title>Shows</title>\n")
	indexFile.WriteString("<style>body { font-family: Arial, sans-serif; }</style>\n")
	indexFile.WriteString("</head>\n<body>\n")
	indexFile.WriteString("<h1>Shows</h1>\n")
	indexFile.WriteString("<ul>\n")

	timestamp := time.Now().Unix()
	for _, show := range shows {
		if len(show.Teams) == 0 {
			continue
		}
		dateStr := show.Date.Format("2006-01-02")
		indexFile.WriteString(fmt.Sprintf("<li><a href=\"%s.html?%d\">%s - %s</a></li>\n", dateStr, timestamp, dateStr, show.Title))
	}

	indexFile.WriteString("</ul>\n")
	indexFile.WriteString("</body>\n</html>")
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

	// Start writing HTML content
	showFile.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	showFile.WriteString(fmt.Sprintf("<title>%s - %s</title>\n", dateStr, show.Title))
	showFile.WriteString("<style>body { font-family: Arial, sans-serif; }</style>\n")
	showFile.WriteString("</head>\n<body>\n")
	showFile.WriteString(fmt.Sprintf("<h1>%s - %s</h1>\n", dateStr, show.Title))
	showFile.WriteString("<a href='index.html'>back</h1>\n")
	showFile.WriteString("<ul>\n")

	// List of types

	timestamp := time.Now().Unix()
	for _, tmplType := range POST_TYPES {
		imageFileName := fmt.Sprintf("%s - %s - %s.jpg?%d", dateStr, show.Title, tmplType, timestamp)
		showFile.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>\n", imageFileName, tmplType))
	}

	showFile.WriteString("</ul>\n")
	showFile.WriteString("</body>\n</html>")
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

	var shows []Show

	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

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

		// Increment the WaitGroup counter
		wg.Add(1)

		// Run screenshot and show page generation in a goroutine
		go func(show Show) {
			defer wg.Done() // Decrement counter when goroutine completes

			SaveScreenshot(tmplFb, show, POST_TYPE_FB)
			SaveScreenshot(tmplInsta, show, POST_TYPE_INSTA)
			SaveScreenshot(tmplSio, show, POST_TYPE_SIO)
			SaveScreenshot(tmplMeetup, show, POST_TYPE_MEETUP)
			SaveScreenshot(tmplStory, show, POST_TYPE_STORY)

			// Generate date-specific HTML file with links
			CreateShowPage(show)
		}(show)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Generate the index.html
	CreateIndex(shows)
}
