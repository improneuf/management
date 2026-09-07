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

func createBannersManifest(workshops []WorkshopConfig) []Banner {
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
		return banners
	}
	file, err := os.Create("output/banners.json")
	if err != nil {
		log.Printf("Failed to create banners.json: %v", err)
		return banners
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(banners); err != nil {
		log.Printf("Failed to encode banners.json: %v", err)
	}

	return banners
}

func createWorkshopGallery(banners []Banner) {
	const galleryTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Impro Neuf Open English Workshops</title>
  <style>
    :root { color-scheme: dark; }
    * { box-sizing: border-box; }
    body { margin: 0; background: #111; color: #fff; font-family: Arial, sans-serif; }
    header { padding: 2rem 1rem 1rem; text-align: center; }
    header h1 { margin: 0 0 .5rem; }
    header a { color: #9fd8ff; }
    main { width: min(1200px, calc(100% - 2rem)); margin: 0 auto 3rem; display: grid; grid-template-columns: repeat(auto-fit, minmax(290px, 1fr)); gap: 1.25rem; }
    article { overflow: hidden; border: 1px solid #333; border-radius: 14px; background: #1d1d1d; box-shadow: 0 8px 28px rgba(0,0,0,.28); }
    article img { display: block; width: 100%; height: auto; }
    .details { padding: 1rem; }
    h2 { margin: 0 0 .6rem; font-size: 1.2rem; }
    p { margin: .35rem 0; color: #ddd; }
    .date { color: #fff; font-weight: 700; }
  </style>
</head>
<body>
  <header>
    <h1>Impro Neuf Open English Workshops</h1>
    <p>Automatically generated workshop banners.</p>
    <a href="/management/workshop-banners.json">Workshop banner JSON</a>
  </header>
  <main>
    {{range .}}{{if not .IsPast}}
    <article>
      <a href="{{index .Images "fb"}}"><img src="{{index .Images "fb"}}" alt="{{.Title}}"></a>
      <div class="details">
        <h2>{{.Title}}</h2>
        <p class="date">{{.Date}} at {{.ShowStart}}</p>
        <p>{{.Room}}</p>
        <p>{{range $index, $name := .Facilitators}}{{if $index}}, {{end}}{{$name}}{{end}}</p>
      </div>
    </article>
    {{end}}{{end}}
  </main>
</body>
</html>`

	tmpl, err := template.New("workshop-gallery").Parse(galleryTemplate)
	if err != nil {
		log.Printf("Failed to parse workshop gallery template: %v", err)
		return
	}

	file, err := os.Create("output/index.html")
	if err != nil {
		log.Printf("Failed to create workshop gallery: %v", err)
		return
	}
	defer file.Close()

	if err := tmpl.Execute(file, banners); err != nil {
		log.Printf("Failed to render workshop gallery: %v", err)
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
		{
			PageTitle:      "Improv for Everyone",
			BackgroundHaze: "banner-bg-4.png",
			MainTitle:      "Improv for Everyone",
			Subtitle:       "Improv for Everyone",
			HostName1:      "Mike Altorjay",
			HostName2:      "",
			HostImage1:     "host_miklos.png",
			HostImage2:     "",
			EventDate:      "Wednesday, September 09, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-5.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, September 16, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-6.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, September 23, 2026",
			Room:           "Galleriet",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-1.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Magnus",
			HostImage1:     "host_magnus.png",
			EventDate:      "Wednesday, September 30, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-2.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, October 07, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-3.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, October 14, 2026",
			Room:           "Betong",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-4.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, October 21, 2026",
			Room:           "Room TBA",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-5.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, October 28, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-6.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, November 04, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-1.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, November 11, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-2.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, November 18, 2026",
			Room:           "Klubbscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-3.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, November 25, 2026",
			Room:           "Teaterscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-4.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, December 02, 2026",
			Room:           "Teaterscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-5.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, December 09, 2026",
			Room:           "Teaterscenen",
		},
		{
			PageTitle:      "Open Improv Workshop",
			BackgroundHaze: "banner-bg-6.png",
			MainTitle:      "Open Improv Workshop",
			Subtitle:       "Free improv workshop for everyone",
			HostName1:      "Impro Neuf",
			HostImage1:     "host_improv.png",
			EventDate:      "Wednesday, December 16, 2026",
			Room:           "Klubbscenen",
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

	banners := createBannersManifest(workshops)
	createWorkshopGallery(banners)

	fmt.Printf("\nGenerated %d HTML files successfully!\n", len(workshops))
	fmt.Println("\nTo convert to PNG, run this command from the html-to-png directory:")
	fmt.Println("cd ../html-to-png")
	fmt.Println("go run convert-workshops.go")
}

func main() {
	GenerateAllWorkshops()
}
