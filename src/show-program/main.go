package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	//SHOW_PROGRAM_SHEET_ID string = "1ejEDxQJIwQ1ougcpWIKTqauT-05PDVT1" // Test Sheet
	SHOW_PROGRAM_SHEET_ID    string = "167cJAqP9fON3ExyLnJLFaJ0MHdu5K--z" // Live Sheet
	SHOW_PROGRAM_SHEET_NAME  string = "ShowProgram"
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
	ctx := context.Background()

	b, err := os.ReadFile(SERVICE_ACCOUNT_KEY_FILE)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(ctx)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	xlsxFilePath := SHOW_PROGRAM_SHEET_ID + ".xlsx"

	// if the file already exists
	_, err = os.Stat(xlsxFilePath)
	fileExistsLocally := !os.IsNotExist(err)
	localFileIsUpToDate := false

	// if the file exists locally, check if it's up to date
	if fileExistsLocally {
		localFileModifiedTime, err := GetLocalFileModifiedTime(xlsxFilePath)
		if err != nil {
			log.Fatalf("Unable to get local file modified date: %v", err)
		}
		googleDriveFileModifiedTime, err := GetGoogleDriveFileModifiedTime(srv, SHOW_PROGRAM_SHEET_ID)
		if err != nil {
			log.Fatalf("Unable to get google drive file modified date: %v", err)
		}

		if localFileModifiedTime.After(googleDriveFileModifiedTime) {
			localFileIsUpToDate = true
		}
	}

	if fileExistsLocally && localFileIsUpToDate {
		fmt.Println("File already exists and is up to date, not downloading again.")
	} else {
		fmt.Println("File does not exist or is out of date, downloading...")

		downloadedFileTemp, err := DownloadFileFromGoogleDrive(srv, SHOW_PROGRAM_SHEET_ID)

		if err != nil {
			log.Fatalf("Unable to download file: %v", err)
		}

		fmt.Println("Downloaded file to: " + downloadedFileTemp)

		// move tempfile to the correct location
		os.Rename(downloadedFileTemp, xlsxFilePath)
	}

	fmt.Println(ReadShowScheduleFromFile(xlsxFilePath))

	funcMap := template.FuncMap{
		"GetTeamPhoto": GetTeamPhoto,
	}

	for _, show := range ReadShowScheduleFromFile(xlsxFilePath) {
		if show.Types[0] != ShowTypeRegular {
			continue
		}

		// Parse the template file
		tmpl, err := template.New("regular.tmpl").Funcs(funcMap).ParseFiles("regular.tmpl")
		if err != nil {
			panic(err)
		}

		// Create output file
		outputFile, err := os.Create("output/" + show.Title + " - " + show.Date.Format("2006-01-02") + ".html")
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
