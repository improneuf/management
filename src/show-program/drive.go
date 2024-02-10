package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// retrieve the last modified date of a file on Google Drive.
func GetGoogleDriveFileModifiedTime(service *drive.Service, fileID string) (time.Time, error) {
	// Retrieve the file's metadata
	file, err := service.Files.Get(fileID).Fields("modifiedTime").Do()
	if err != nil {
		return time.Time{}, err
	}
	// Parse the ModifiedTime
	modifiedTime, err := time.Parse(time.RFC3339, file.ModifiedTime)
	if err != nil {
		fmt.Println("Error parsing ModifiedTime:", err)
		return time.Time{}, err
	}

	fmt.Println("Parsed Google Drive File Modified Time:", modifiedTime)
	return modifiedTime, nil
}

// Download a file from google drive to a temporary file and return the path to the file
func DownloadFileFromGoogleDrive(service *drive.Service, file_id string) (string, error) {
	fmt.Println("Downloading file with id " + file_id + " from Google Drive...")

	resp, err := service.Files.Get(file_id).Download()
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded content: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}

	fmt.Println("Created a temporary file: " + tmpFile.Name())

	err = os.WriteFile(tmpFile.Name(), data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Println("Wrote downloaded content to the temporary file: " + tmpFile.Name() + ".")

	return tmpFile.Name(), nil
}

func GetGoogleSheetsPath(sheetId string) string {
	xlsxFilePath := sheetId + ".xlsx"
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
	return xlsxFilePath
}
