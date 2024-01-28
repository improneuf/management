package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
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

// ReadExcelFile reads and prints the contents of the given Excel file.
func ReadShowScheduleFromFile(filePath string) []Show {
	events := make([]Show, 0)

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

			var teams []string
			var date time.Time
			var day string
			var crewSjefTeam string
			var showLanguages []ShowLanguage
			var showTypes []ShowType

			for c, cell := range row {
				if c > 9 {
					break
				}
				cellTrimmed := strings.TrimSpace(cell)
				fmt.Printf("|%s|", cellTrimmed)
				switch c {
				case 0:
					date, err = time.Parse("2 Jan 2006", cellTrimmed)
					if err != nil {
						date2, err2 := time.Parse("2-Jan-06", cellTrimmed)
						if err2 != nil {
							log.Fatalf("Failed to parse date: %v", cellTrimmed)
						}
						date = date2
					}
				case 1:
					day = cellTrimmed
				case 2:
					crewSjefTeam = cellTrimmed
				case 4, 5, 6, 7, 8:
					if cellTrimmed != "" {
						teams = append(teams, cellTrimmed)
						showTypes = append(showTypes, getShowType(cellTrimmed))
					}
				case 9:
					showLanguagesStr := strings.Split(cellTrimmed, "/")
					for _, languageStr := range showLanguagesStr {
						languageStr = strings.TrimSpace(languageStr)
						showLanguages = append(showLanguages, getShowLanguage(languageStr))
					}
				}
			}
			showTypes = removeDuplicateShowTypes(showTypes)
			event := Show{
				Date:          date,
				Day:           day,
				CrewSjefTeam:  crewSjefTeam,
				Teams:         teams,
				ShowLanguages: showLanguages,
				ShowTypes:     showTypes,
			}
			events = append(events, event)
			fmt.Println()
		}
	}
	return events
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
}
