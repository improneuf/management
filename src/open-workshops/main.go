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
		{
			PageTitle:      "Celebrating Mistakes",
			BackgroundHaze: "banner-bg-6.png",
			MainTitle:      "Celebrating Mistakes",
			Subtitle:       "Celebrating Mistakes",
			HostName1:      "Mike Altorjay",
			HostName2:      "",
			HostImage1:     "host_miklos.png",
			HostImage2:     "",
			EventDate:      "Wednesday, February 18, 2026",
			Room:           "Betong",
		},
		{
			PageTitle:      "Make things that are not real, real",
			BackgroundHaze: "banner-bg-1.png",
			MainTitle:      "Make things that are not real, real",
			Subtitle:       "Make things that are not real, real",
			HostName1:      "Liv Grøthe",
			HostName2:      "Nikki Michelle Soo",
			HostImage1:     "host_liv.png",
			HostImage2:     "host_nikki.png",
			EventDate:      "Wednesday, February 25, 2026",
			Room:           "Teaterscenen",
		},
		{
			PageTitle:      "Let's tell a story",
			BackgroundHaze: "banner-bg-2.png",
			MainTitle:      "Let's tell a story",
			Subtitle:       "Let's tell a story",
			HostName1:      "Anne Ko",
			HostName2:      "Julie Outterside",
			HostImage1:     "host_anne.png",
			HostImage2:     "host_julie.png",
			EventDate:      "Wednesday, March 4, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "The power of listening",
			BackgroundHaze: "banner-bg-3.png",
			MainTitle:      "The power of listening",
			Subtitle:       "The power of listening",
			HostName1:      "Naya PK",
			HostName2:      "Remi Rossi",
			HostImage1:     "host_naya.png",
			HostImage2:     "host_remi.png",
			EventDate:      "Wednesday, March 11, 2026",
			Room:           "Sysesalen",
		},
		{
			PageTitle:      "The Loser Takes It All",
			BackgroundHaze: "banner-bg-4.png",
			MainTitle:      "The Loser Takes It All",
			Subtitle:       "The Loser Takes It All",
			HostName1:      "Peter Müller",
			HostName2:      "",
			HostImage1:     "host_peter.png",
			HostImage2:     "",
			EventDate:      "Wednesday, March 18, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Just Go For It!",
			BackgroundHaze: "banner-bg-5.png",
			MainTitle:      "Just Go For It!",
			Subtitle:       "Just Go For It!",
			HostName1:      "Cole Grabinsky",
			HostName2:      "Mari Svenkerud",
			HostImage1:     "host_cole.png",
			HostImage2:     "host_mari_sven.png",
			EventDate:      "Wednesday, March 25, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Using chores to empower your improv",
			BackgroundHaze: "banner-bg-6.png",
			MainTitle:      "Using chores to empower your improv",
			Subtitle:       "Using chores to empower your improv",
			HostName1:      "Paul Omar",
			HostName2:      "",
			HostImage1:     "host_paul.png",
			HostImage2:     "",
			EventDate:      "Wednesday, April 8, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Be present and react",
			BackgroundHaze: "banner-bg-1.png",
			MainTitle:      "Be present and react",
			Subtitle:       "Be present and react",
			HostName1:      "Siddanth Nayak",
			HostName2:      "Alejandro Saksida",
			HostImage1:     "host_sid.png",
			HostImage2:     "host_alejandro.png",
			EventDate:      "Wednesday, April 15, 2026",
			Room:           "Betong",
		},
		{
			PageTitle:      "Be Great to Play with",
			BackgroundHaze: "banner-bg-2.png",
			MainTitle:      "Be Great to Play with",
			Subtitle:       "Be Great to Play with",
			HostName1:      "Kevin Gow",
			HostName2:      "",
			HostImage1:     "host_kevin.png",
			HostImage2:     "",
			EventDate:      "Wednesday, April 22, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Who, what, where - the basics",
			BackgroundHaze: "banner-bg-3.png",
			MainTitle:      "Who, what, where - the basics",
			Subtitle:       "Who, what, where - the basics",
			HostName1:      "Mari Olimstad",
			HostName2:      "Øyvind Dragsten",
			HostImage1:     "host_mari.png",
			HostImage2:     "host_øyvind.png",
			EventDate:      "Wednesday, April 29, 2026",
			Room:           "Sysesalen",
		},
		{
			PageTitle:      "Flow",
			BackgroundHaze: "banner-bg-4.png",
			MainTitle:      "Flow",
			Subtitle:       "Flow",
			HostName1:      "Natalia Wydra",
			HostName2:      "Brigitta Peto",
			HostImage1:     "host_natalia.png",
			HostImage2:     "host_brigitta.png",
			EventDate:      "Wednesday, May 6, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Between the lines",
			BackgroundHaze: "banner-bg-5.png",
			MainTitle:      "Between the lines",
			Subtitle:       "Between the lines",
			HostName1:      "Rifat Naim",
			HostName2:      "Maike Helder",
			HostImage1:     "host_rifat.png",
			HostImage2:     "host_maike.png",
			EventDate:      "Wednesday, May 13, 2026",
			Room:           "Betong",
		},
		{
			PageTitle:      "Building scenes together",
			BackgroundHaze: "banner-bg-6.png",
			MainTitle:      "Building scenes together",
			Subtitle:       "Building scenes together",
			HostName1:      "Magnus Seines",
			HostName2:      "Ruidi Collins",
			HostImage1:     "host_magnus.png",
			HostImage2:     "host_ruidi.png",
			EventDate:      "Wednesday, May 20, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Improv in the fast lane",
			BackgroundHaze: "banner-bg-1.png",
			MainTitle:      "Improv in the fast lane",
			Subtitle:       "Improv in the fast lane",
			HostName1:      "Peter Müller",
			HostName2:      "Erlend Lunde",
			HostImage1:     "host_peter.png",
			HostImage2:     "host_erlend.png",
			EventDate:      "Wednesday, June 3, 2026",
			Room:           "Betong",
		},
		{
			PageTitle:      "Emotional Playground",
			BackgroundHaze: "banner-bg-2.png",
			MainTitle:      "Emotional Playground",
			Subtitle:       "Emotional Playground",
			HostName1:      "Santiago Beltran",
			HostName2:      "Hanna Saastamoinen",
			HostImage1:     "host_santiago.png",
			HostImage2:     "host_hanna.png",
			EventDate:      "Wednesday, June 10, 2026",
			Room:           "Teaterscenen",
		},
		{
			PageTitle:      "The Little Things",
			BackgroundHaze: "banner-bg-3.png",
			MainTitle:      "The Little Things",
			Subtitle:       "The Little Things",
			HostName1:      "Julie Outterside",
			HostName2:      "",
			HostImage1:     "host_julie.png",
			HostImage2:     "",
			EventDate:      "Wednesday, June 17, 2026",
			Room:           "Teaterscenen",
		},
		{
			PageTitle:      "Mind Meld",
			BackgroundHaze: "banner-bg-4.png",
			MainTitle:      "Mind Meld",
			Subtitle:       "Mind Meld",
			HostName1:      "Kevin Gow",
			HostName2:      "",
			HostImage1:     "host_kevin.png",
			HostImage2:     "",
			EventDate:      "Wednesday, June 24, 2026",
			Room:           "Teaterscenen",
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
