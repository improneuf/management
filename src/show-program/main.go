package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	//SHOW_PROGRAM_SHEET_ID string = "1ejEDxQJIwQ1ougcpWIKTqauT-05PDVT1" // Test Sheet
	SHOW_PROGRAM_SHEET_ID    string = "167cJAqP9fON3ExyLnJLFaJ0MHdu5K--z" // Live Sheet
	SHOW_PROGRAM_SHEET_NAME  string = "ShowProgram"
	SERVICE_ACCOUNT_KEY_FILE string = "impro-neuf-management-99d59b5d3102.json"
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

// ReadExcelFile reads and prints the contents of the given Excel file.
func ReadShowScheduleFromFile(filePath string) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	// Get the names of all the sheets
	sheets := f.GetSheetList()

	// Iterate through each sheet and print its contents
	for _, sheet := range sheets {
		fmt.Printf("Sheet: %s\n", sheet)
		if sheet != SHOW_PROGRAM_SHEET_NAME {
			continue
		}
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Fatalf("Failed to get rows for sheet %s: %v", sheet, err)
		}

		for r, row := range rows {
			if r < 10 {
				continue
			}
			for c, cell := range row {
				if c > 9 {
					break
				}
				fmt.Print(cell, "\t")
			}
			fmt.Println()
		}
	}
}

// Retrieve a token, save the token, then return the token
func tokenFromFile(file string) (*oauth2.Token, error) {
	content, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{}
	err = json.Unmarshal(content, tok)
	if err != nil {
		return nil, err
	}

	return tok, err
}

// Request a token from the web, then returns the retrieved token
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// Generate a URL and ask the user to visit it for authorization.
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	// Wait for the authorization code to be pasted in the terminal.
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	// Exchange the authorization code for an access token.
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web %v", err)
	}

	return token, nil
}

// getServiceAccountToken retrieves an OAuth2 token using service account credentials.
func getServiceAccountToken(serviceAccountKeyFile string) (*oauth2.Token, error) {
	// Read the service account key file
	data, err := os.ReadFile(serviceAccountKeyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read service account key file: %v", err)
	}

	// Parse the service account key
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/drive.readonly")
	if err != nil {
		return nil, fmt.Errorf("unable to parse service account key file: %v", err)
	}

	// Get the HTTP client using the service account's JWT configuration
	client := conf.Client(context.Background())

	// Fetch the token
	token, err := client.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from service account: %v", err)
	}

	return token, nil
}

// Save a token to a file path
func saveToken(path string, token *oauth2.Token) {
	// Open the file for writing.
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()

	// Serialize the token to JSON and write it to the file.
	json.NewEncoder(file).Encode(token)
}

func main() {
	// Load your client_secret.json
	b, err := os.ReadFile(SERVICE_ACCOUNT_KEY_FILE)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.JWTConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(context.Background())

	srv, err := drive.New(client)
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

	ReadShowScheduleFromFile(xlsxFilePath)
}
