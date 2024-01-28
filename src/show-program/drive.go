package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"google.golang.org/api/drive/v3"
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
