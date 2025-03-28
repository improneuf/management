package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/chromedp/chromedp"
)

func main() {
	// Check if the filename is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <html_file_path>")
		os.Exit(1)
	}

	// Get the HTML file path from the command-line argument
	fileURL := os.Args[1]

	// Extract the file name (e.g., "julie-cole.html") and then the base name (e.g., "julie-cole")
	htmlFilename := path.Base(fileURL)
	base := strings.TrimSuffix(htmlFilename, path.Ext(htmlFilename))
	fbFileName := fmt.Sprintf("%s-fb.jpg", base)
	meetupFileName := fmt.Sprintf("%s-meetup.jpg", base)

	// Create a chromedp context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Screenshot for Facebook dimensions (1920x1005)
	var bufFB []byte
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1005),
		chromedp.Navigate(fileURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the zoom level by scaling the CSS
			return chromedp.Evaluate(`document.body.style.zoom = "1"`, nil).Do(ctx)
		}),
		chromedp.CaptureScreenshot(&bufFB),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fbFileName, bufFB, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved screenshot as", fbFileName)

	// Screenshot for Meetup dimensions (1920x1080)
	var bufMeetup []byte
	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(fileURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the zoom level by scaling the CSS
			return chromedp.Evaluate(`document.body.style.zoom = "1"`, nil).Do(ctx)
		}),
		chromedp.CaptureScreenshot(&bufMeetup),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(meetupFileName, bufMeetup, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved screenshot as", meetupFileName)
}
