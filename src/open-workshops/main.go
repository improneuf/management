package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
)

// PageData holds the dynamic values for the page.
type PageData struct {
	PageTitle           string
	BackgroundHazeURL   string
	LogoPath            string
	MainTitle           string
	Subtitle            string
	HostName1           string
	HostName2           string
	HostImage1          string
	HostImage2          string
	WorkshopDescription string
	EventDate           string
	EventTime           string
	Location            string
	Room                string
	SignUpPrompt        string
}

// WorkshopConfig holds all the parameters for a workshop
type WorkshopConfig struct {
	PageTitle      string
	BackgroundHaze string
	MainTitle      string
	Subtitle       string
	HostName1      string
	HostName2      string // Empty string means no second host
	HostImage1     string
	HostImage2     string // Empty string means no second host image
	EventDate      string
	Room           string
}

// sanitizeFilename creates a Windows-safe filename by removing or replacing invalid characters
func sanitizeFilename(title string) string {
	// Replace spaces with hyphens
	safeTitle := strings.ReplaceAll(title, " ", "-")

	// Remove or replace invalid characters for Windows filenames
	// Invalid characters: < > : " | ? * \ /
	// Also remove other potentially problematic characters
	invalidChars := regexp.MustCompile(`[<>:"|?*\\/!@#$%^&()+={}[\]~` + "`" + `;,]`)
	safeTitle = invalidChars.ReplaceAllString(safeTitle, "")

	// Remove multiple consecutive hyphens
	safeTitle = regexp.MustCompile(`-+`).ReplaceAllString(safeTitle, "-")

	// Remove leading/trailing hyphens
	safeTitle = strings.Trim(safeTitle, "-")

	// Ensure the filename is not empty
	if safeTitle == "" {
		safeTitle = "workshop"
	}

	return safeTitle
}

// GenerateAllWorkshops creates HTML files for all workshop combinations
func GenerateAllWorkshops() {
	// Define all the workshop configurations
	workshops := []WorkshopConfig{
		{
			PageTitle:      "An improduction to freedom and fun!",
			BackgroundHaze: "banner-bg-1.png",
			MainTitle:      "An improduction to freedom and fun!",
			Subtitle:       "An improduction to freedom and fun!",
			HostName1:      "Peter Müller",
			HostName2:      "",
			HostImage1:     "host_peter.png",
			HostImage2:     "",
			EventDate:      "Wednesday, January 14, 2026",
			Room:           "Betong",
		},
		{
			PageTitle:      "When characters meet",
			BackgroundHaze: "banner-bg-2.png",
			MainTitle:      "When characters meet",
			Subtitle:       "When characters meet",
			HostName1:      "Kevin Gow",
			HostName2:      "",
			HostImage1:     "host_kevin.png",
			HostImage2:     "",
			EventDate:      "Wednesday, January 21, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Let the Body Lead",
			BackgroundHaze: "banner-bg-3.png",
			MainTitle:      "Let the Body Lead",
			Subtitle:       "Let the Body Lead",
			HostName1:      "Cole Grabinsky",
			HostName2:      "Anjitha S.g.",
			HostImage1:     "host_cole.png",
			HostImage2:     "host_anjitha.png",
			EventDate:      "Wednesday, January 28, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Freedom to fail",
			BackgroundHaze: "banner-bg-4.png",
			MainTitle:      "Freedom to fail",
			Subtitle:       "Freedom to fail",
			HostName1:      "Natalia Wydra",
			HostName2:      "",
			HostImage1:     "host_natalia.png",
			HostImage2:     "",
			EventDate:      "Wednesday, February 4, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Making Bold Choices",
			BackgroundHaze: "banner-bg-5.png",
			MainTitle:      "Making Bold Choices",
			Subtitle:       "Making Bold Choices",
			HostName1:      "Carlos Moreno",
			HostName2:      "Ruidi Collins",
			HostImage1:     "host_carlos.png",
			HostImage2:     "host_ruidi.png",
			EventDate:      "Wednesday, February 11, 2026",
			Room:           "Betong",
		},
	}

	// Create output directory
	if err := os.MkdirAll(".", 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Parse the template
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Generate HTML for each workshop
	for i, workshop := range workshops {
		// Create filename-safe version of the title
		safeTitle := sanitizeFilename(workshop.MainTitle)

		filename := fmt.Sprintf("workshop-%d-%s.html", i+1, safeTitle)

		// Create the HTML file
		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Failed to create %s: %v", filename, err)
			continue
		}
		defer file.Close()

		// Prepare the data for the template
		data := PageData{
			PageTitle:         workshop.PageTitle,
			BackgroundHazeURL: workshop.BackgroundHaze,
			LogoPath:          "logo.png",
			MainTitle:         workshop.MainTitle,
			Subtitle:          workshop.Subtitle,
			HostName1:         workshop.HostName1,
			HostName2:         workshop.HostName2,
			HostImage1:        workshop.HostImage1,
			HostImage2:        workshop.HostImage2,
			EventDate:         workshop.EventDate,
			EventTime:         "6:00 PM - 8:00 PM",
			Location:          "Chateau Neuf",
			Room:              workshop.Room,
			SignUpPrompt:      "Sign up now!",
		}

		// Execute the template
		if err := tmpl.Execute(file, data); err != nil {
			log.Printf("Failed to execute template for %s: %v", filename, err)
			continue
		}

		fmt.Printf("Generated: %s\n", filename)
	}

	fmt.Printf("\nGenerated %d HTML files successfully!\n", len(workshops))
	fmt.Println("\nTo convert to PNG, run this command from the html-to-png directory:")
	fmt.Println("cd ../html-to-png")
	fmt.Println("go run convert-workshops.go")
}

// main generates all workshop HTML files and exits
func main() {
	GenerateAllWorkshops()
}
