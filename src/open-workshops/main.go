package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

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

type WorkshopConfig struct {
	PageTitle      string
	BackgroundHaze string
	MainTitle      string
	Subtitle       string
	HostName1      string
	HostName2      string
	HostImage1     string
	HostImage2     string
	EventDate      string
	Room           string
}

type Banner struct {
	Date         string            `json:"date"`
	Filename     string            `json:"filename"`
	URL          string            `json:"url"`
	Teams        []string          `json:"teams"`
	Facilitators []string          `json:"facilitators"`
	Title        string            `json:"title"`
	ShowStart    string            `json:"showStart"`
	Room         string            `json:"room"`
	IsPast       bool              `json:"isPast"`
	Types        []string          `json:"types"`
	Images       map[string]string `json:"images"`
}

func createBannersManifest(workshops []WorkshopConfig) {
	const base = "https://improneuf.github.io/management/open-workshops/output/"
	today := time.Now().Truncate(24 * time.Hour)
	banners := make([]Banner, 0, len(workshops))

	for i, workshop := range workshops {
		safeTitle := sanitizeFilename(workshop.MainTitle)
		filename := fmt.Sprintf("workshop-%d-%s-fb.jpg", i+1, safeTitle)
		date, err := time.Parse("Monday, January 02, 2006", workshop.EventDate)
		if err != nil {
			log.Printf("Failed to parse workshop date %q: %v", workshop.EventDate, err)
			continue
		}

		images := map[string]string{}
		for _, postType := range []string{"fb", "meetup"} {
			imageFilename := fmt.Sprintf("workshop-%d-%s-%s.jpg", i+1, safeTitle, postType)
			images[postType] = base + url.PathEscape(imageFilename)
		}

		hosts := []string{workshop.HostName1}
		if workshop.HostName2 != "" {
			hosts = append(hosts, workshop.HostName2)
		}

		banners = append(banners, Banner{
			Date:         date.Format("2006-01-02"),
			Filename:     filename,
			URL:          images["fb"],
			Teams:        hosts,
			Facilitators: hosts,
			Title:        workshop.MainTitle,
			ShowStart:    "18:00",
			Room:         workshop.Room,
			IsPast:       date.Before(today),
			Types:        []string{"Workshop"},
			Images:       images,
		})
	}

	if err := os.MkdirAll("output", 0755); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		return
	}

	file, err := os.Create("output/banners.json")
	if err != nil {
		log.Printf("Failed to create banners.json: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(banners); err != nil {
		log.Printf("Failed to encode banners.json: %v", err)
	}
}

func sanitizeFilename(title string) string {
	safeTitle := strings.ReplaceAll(title, " ", "-")

	invalidChars := regexp.MustCompile(`[<>:"|?*\\/!@#$%^&()+={}[\]~` + "`" + `;,]`)
	safeTitle = invalidChars.ReplaceAllString(safeTitle, "")

	safeTitle = regexp.MustCompile(`-+`).ReplaceAllString(safeTitle, "-")

	safeTitle = strings.Trim(safeTitle, "-")

	if safeTitle == "" {
		safeTitle = "workshop"
	}

	return safeTitle
}

func GenerateAllWorkshops() {
	workshops := []WorkshopConfig{
		{
			PageTitle:      "Improv for everyone",
			BackgroundHaze: "banner-bg-1.png",
			MainTitle:      "Improv for everyone",
			Subtitle:       "Improv for everyone",
			HostName1:      "Kevin Gow",
			HostName2:      "",
			HostImage1:     "host_kevin.png",
			HostImage2:     "",
			EventDate:      "Wednesday, August 19, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Short Form Games",
			BackgroundHaze: "banner-bg-2.png",
			MainTitle:      "Short Form Games",
			Subtitle:       "Short Form Games",
			HostName1:      "Liv Grøthe",
			HostName2:      "Nikki Michelle Soo",
			HostImage1:     "host_liv.png",
			HostImage2:     "host_nikki.png",
			EventDate:      "Wednesday, August 26, 2026",
			Room:           "Galleriet",
		},
		{
			PageTitle:      "Turning Mistakes into Magic",
			BackgroundHaze: "banner-bg-3.png",
			MainTitle:      "Turning Mistakes into Magic",
			Subtitle:       "Turning Mistakes into Magic",
			HostName1:      "Remi Rossi",
			HostName2:      "Santiago Beltran",
			HostImage1:     "host_remi.png",
			HostImage2:     "host_santiago.png",
			EventDate:      "Wednesday, September 02, 2026",
			Room:           "Betong",
		},
	}

	if err := os.MkdirAll(".", 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	for i, workshop := range workshops {
		safeTitle := sanitizeFilename(workshop.MainTitle)

		filename := fmt.Sprintf("workshop-%d-%s.html", i+1, safeTitle)

		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Failed to create %s: %v", filename, err)
			continue
		}
		defer file.Close()

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

		if err := tmpl.Execute(file, data); err != nil {
			log.Printf("Failed to execute template for %s: %v", filename, err)
			continue
		}

		fmt.Printf("Generated: %s\n", filename)
	}

	createBannersManifest(workshops)

	fmt.Printf("\nGenerated %d HTML files successfully!\n", len(workshops))
	fmt.Println("\nTo convert to PNG, run this command from the html-to-png directory:")
	fmt.Println("cd ../html-to-png")
	fmt.Println("go run convert-workshops.go")
}

func main() {
	GenerateAllWorkshops()
}
