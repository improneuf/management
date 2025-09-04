package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
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

// GenerateAllWorkshops creates HTML files for all workshop combinations
func GenerateAllWorkshops() {
	// Define all the workshop configurations
	workshops := []WorkshopConfig{
		{
			PageTitle:      "Building the Basics of a Scene",
			BackgroundHaze: "bg_dunes.webp",
			MainTitle:      "Building the Basics of a Scene",
			Subtitle:       "Improv Theatre",
			HostName1:      "Barathy Pirabahar",
			HostName2:      "Cole Grabinsky",
			HostImage1:     "host_barathy.png",
			HostImage2:     "host_cole.png",
			EventDate:      "Wednesday, August 20, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "The FUN in Fundamentals: Start with Agreement",
			BackgroundHaze: "bg_fire.webp",
			MainTitle:      "The FUN in Fundamentals: Start with Agreement",
			Subtitle:       "Start with Agreement",
			HostName1:      "Kevin Gow",
			HostName2:      "",
			HostImage1:     "host_kevin.png",
			HostImage2:     "",
			EventDate:      "Wednesday, August 27, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Celebrating Our Mistakes",
			BackgroundHaze: "bg_snow.webp",
			MainTitle:      "Celebrating Our Mistakes",
			Subtitle:       "Celebrating Our Mistakes",
			HostName1:      "Liv Grøthe",
			HostName2:      "Nikki Michelle Soo",
			HostImage1:     "host_liv.png",
			HostImage2:     "host_nikki.png",
			EventDate:      "Wednesday, September 3, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "The ABC of Improv - Yes And",
			BackgroundHaze: "bg_tides.webp",
			MainTitle:      "The ABC of Improv - Yes And",
			Subtitle:       "The ABC of Improv - Yes And",
			HostName1:      "Auritro Paldas",
			HostName2:      "Naya Kouzilou",
			HostImage1:     "host_auritro.png",
			HostImage2:     "host_naya2.png",
			EventDate:      "Wednesday, September 10, 2025",
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
		safeTitle := strings.ReplaceAll(workshop.MainTitle, " ", "-")
		safeTitle = strings.ReplaceAll(safeTitle, ":", "")
		safeTitle = strings.ReplaceAll(safeTitle, "***", "star")

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
			PageTitle:           workshop.PageTitle,
			BackgroundHazeURL:   workshop.BackgroundHaze,
			LogoPath:            "logo.png",
			MainTitle:           workshop.MainTitle,
			Subtitle:            workshop.Subtitle,
			HostName1:           workshop.HostName1,
			HostName2:           workshop.HostName2,
			HostImage1:          workshop.HostImage1,
			HostImage2:          workshop.HostImage2,
			WorkshopDescription: "A free English open workshop by Impro Neuf.",
			EventDate:           workshop.EventDate,
			EventTime:           "6:00 PM - 8:00 PM",
			Location:            "Chateau Neuf",
			Room:                workshop.Room,
			SignUpPrompt:        "Sign up now!",
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
