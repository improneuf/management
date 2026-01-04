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
		{
			PageTitle:      "Connecting the Dots",
			BackgroundHaze: "bg_wind.webp",
			MainTitle:      "Connecting the dots",
			Subtitle:       "Connecting the dots",
			HostName1:      "Hanna Saastamoinen",
			HostName2:      "Magnus Seines",
			HostImage1:     "host_hanna.png",
			HostImage2:     "host_magnus.png",
			EventDate:      "Wednesday, September 17, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Who Do You Think You Are?!",
			BackgroundHaze: "bg_bokeh.webp",
			MainTitle:      "Who Do You Think You Are?!",
			Subtitle:       "Who Do You Think You Are?!",
			HostName1:      "India Anderson",
			HostName2:      "Peter Müller",
			HostImage1:     "host_india.png",
			HostImage2:     "host_peter.png",
			EventDate:      "Wednesday, September 24, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "The core of FUN",
			BackgroundHaze: "bg6.webp",
			MainTitle:      "The core of FUN",
			Subtitle:       "The core of FUN",
			HostName1:      "Carlos Moreno",
			HostName2:      "",
			HostImage1:     "host_carlos.png",
			HostImage2:     "",
			EventDate:      "Wednesday, October 1, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Playing with Emotions",
			BackgroundHaze: "bg8.webp",
			MainTitle:      "Playing with Emotions",
			Subtitle:       "Playing with Emotions",
			HostName1:      "Natalia Wydra",
			HostName2:      "",
			HostImage1:     "host_natalia.png",
			HostImage2:     "",
			EventDate:      "Wednesday, October 8, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Improvising things out of thin air",
			BackgroundHaze: "bg9.webp",
			MainTitle:      "Improvising things out of thin air",
			Subtitle:       "Improvising things out of thin air",
			HostName1:      "Mari Olimstad",
			HostName2:      "Liv Grøthe",
			HostImage1:     "host_mari.png",
			HostImage2:     "host_liv3.png",
			EventDate:      "Wednesday, November 5, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Let your body speak",
			BackgroundHaze: "bg10.webp",
			MainTitle:      "Let your body speak",
			Subtitle:       "Let your body speak",
			HostName1:      "Mike Altorjay",
			HostName2:      "Naya Kouzilou",
			HostImage1:     "host_miklos.png",
			HostImage2:     "host_naya.png",
			EventDate:      "Wednesday, November 12, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Playing with status",
			BackgroundHaze: "bg11.webp",
			MainTitle:      "Playing with status",
			Subtitle:       "Playing with status",
			HostName1:      "Anjitha S.G.",
			HostName2:      "Siddanth Nayak",
			HostImage1:     "host_anjitha.png",
			HostImage2:     "host_sid.png",
			EventDate:      "Wednesday, November 19, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Fun and games",
			BackgroundHaze: "bg12.webp",
			MainTitle:      "Fun and games",
			Subtitle:       "Fun and games",
			HostName1:      "Julie Outterside",
			HostName2:      "Eva Coder",
			HostImage1:     "host_julie.png",
			HostImage2:     "host_eva.png",
			EventDate:      "Wednesday, November 26, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Game of the scene",
			BackgroundHaze: "bg13.webp",
			MainTitle:      "Game of the scene",
			Subtitle:       "Game of the scene",
			HostName1:      "Santiago Beltran",
			HostName2:      "Remi Rossi",
			HostImage1:     "host_santiago.png",
			HostImage2:     "host_remi.png",
			EventDate:      "Wednesday, December 3, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Play in the chaos",
			BackgroundHaze: "bg14.webp",
			MainTitle:      "Play in the chaos",
			Subtitle:       "Play in the chaos",
			HostName1:      "Kristīne",
			HostName2:      "Naima",
			HostImage1:     "host_kristine2.png",
			HostImage2:     "host_naima.png",
			EventDate:      "Wednesday, December 10, 2025",
			Room:           "Betong",
		},
		{
			PageTitle:      "Playful! Present! Perfect!",
			BackgroundHaze: "bg15.webp",
			MainTitle:      "Playful! Present! Perfect!",
			Subtitle:       "Playful! Present! Perfect!",
			HostName1:      "Kevin Gow",
			HostName2:      "",
			HostImage1:     "host_kevin.png",
			HostImage2:     "",
			EventDate:      "Wednesday, December 17, 2025",
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
			EventTime:           "6:00 PM - 7:30 PM",
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
