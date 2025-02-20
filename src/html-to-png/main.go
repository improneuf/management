package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	// Create a context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Capture screenshot of an entire webpage in JPEG format
	var buf []byte
	if err := chromedp.Run(ctx,
	    chromedp.EmulateViewport(1920, 1005),
		//chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(`
			file:///Users/pravindahal/impro-neuf-management/src/open-workshops/haley-brigi.html
		`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the zoom level by scaling the CSS
			return chromedp.Evaluate(`document.body.style.zoom = "1"`, nil).Do(ctx)
		}),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		log.Fatal(err)
	}

	// Save the screenshot to a file
	if err := ioutil.WriteFile("screenshot.jpg", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
