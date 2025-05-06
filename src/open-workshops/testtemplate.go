package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
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

// Global state (for demonstration) and a mutex for safe concurrent access.
var (
	pageData = PageData{
		PageTitle:           "Dancing with the Dunes: Celebrating Spontaneity",
		BackgroundHazeURL:   "bg_dunes.webp",
		LogoPath:            "logo.png",
		MainTitle:           "Emotional Whirlwinds",
		Subtitle:            "", // Initially empty
		HostName1:           "Barathy Pirabahar",
		HostName2:           "Anjitha S.g.",
		HostImage1:          "host_barathy.png",
		HostImage2:          "host_anjitha.png",
		WorkshopDescription: "A free English open workshop by Impro Neuf.",
		EventDate:           "Wednesday, April 9, 2025",
		EventTime:           "7:00 PM - 9:00 PM",
		Location:            "Chateau Neuf",
		Room:                "Betong",
		SignUpPrompt:        "Sign up now!",
	}
	mu sync.Mutex
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/subtitle", subtitleHandler)
	http.HandleFunc("/subtitle/edit", subtitleEditHandler)
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// indexHandler renders the full page template.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, pageData); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// updateHandler processes htmx updates from editable fields.
func updateHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	field := r.FormValue("field")
	value := r.FormValue("value")
	log.Printf("Update: field=%s, value=%s", field, value)
	switch field {
	case "MainTitle":
		pageData.MainTitle = value
		// Re-render the main title snippet.
		fmt.Fprintf(w, `<h1 contenteditable="true" data-field="MainTitle" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-5xl md:text-7xl font-bold tracking-widest text-transparent bg-clip-text bg-gradient-to-r from-green-200 via-pink-300 to-yellow-100" style="font-family: 'Fredoka One', sans-serif; padding-bottom: 0.1em;">%s</h1>`, pageData.MainTitle)
	case "Subtitle":
		pageData.Subtitle = value
		// If empty, render the Add subtitle button; otherwise, the editable field.
		fmt.Fprint(w, renderSubtitleSnippet())
	case "HostName1":
		pageData.HostName1 = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="HostName1" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-3xl font-semibold text-fuchsia-200" style="font-family: 'Fredoka One', sans-serif;">%s</p>`, pageData.HostName1)
	case "HostName2":
		pageData.HostName2 = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="HostName2" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-3xl font-semibold text-fuchsia-200" style="font-family: 'Fredoka One', sans-serif;">%s</p>`, pageData.HostName2)
	case "WorkshopDescription":
		pageData.WorkshopDescription = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="WorkshopDescription" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-lg md:text-2xl text-pink-100" style="font-family: 'Montserrat', sans-serif; font-weight: bolder;">%s</p>`, pageData.WorkshopDescription)
	case "EventDate":
		pageData.EventDate = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="EventDate" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-xl md:text-3xl text-indigo-200">📅 <span class="font-bold text-yellow-300">Date:</span> %s</p>`, pageData.EventDate)
	case "EventTime":
		pageData.EventTime = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="EventTime" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-xl md:text-3xl text-purple-200">⏰ <span class="font-bold text-lime-200">Time:</span> %s</p>`, pageData.EventTime)
	case "Location":
		pageData.Location = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="Location" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-xl md:text-3xl text-teal-200">📍 <span class="font-bold text-yellow-300">Location:</span> %s</p>`, pageData.Location)
	case "Room":
		pageData.Room = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="Room" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-xl md:text-3xl text-emerald-200">🏢 <span class="font-bold text-pink-100">Room:</span> %s</p>`, pageData.Room)
	case "SignUpPrompt":
		pageData.SignUpPrompt = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="SignUpPrompt" hx-post="/update" hx-trigger="blur" hx-target="#title-container" class="text-xl md:text-2xl font-bold text-yellow-300">%s</p>`, pageData.SignUpPrompt)
	default:
		http.Error(w, "Unknown field", http.StatusBadRequest)
	}
}

// subtitleHandler renders the subtitle container snippet.
func subtitleHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprint(w, renderSubtitleSnippet())
}

// subtitleEditHandler returns an editable subtitle field with a placeholder.
func subtitleEditHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	// Return the editable field with a placeholder value.
	// Note: We do not update the global state here yet.
	s := `<h2 contenteditable="true" data-field="Subtitle" hx-post="/update" hx-trigger="blur" hx-target="#subtitle-container" hx-include="[data-field=Subtitle]" class="text-2xl md:text-3xl font-bold tracking-widest text-pink-200 mt-2" style="font-family: 'Fredoka One', sans-serif;">subtitle here</h2>`
	fmt.Fprint(w, s)
}

// renderSubtitleSnippet returns the HTML snippet for the subtitle container.
func renderSubtitleSnippet() string {
	if pageData.Subtitle == "" {
		// The button is hidden by default ("hidden") and shown on hover (group-hover:inline-block).
		return `<button id="addSubtitleButton" hx-get="/subtitle/edit" hx-target="#subtitle-container" class="mt-2 text-blue-500 underline hidden group-hover:inline-block">Add subtitle</button>`
	}
	return fmt.Sprintf(`<h2 contenteditable="true" data-field="Subtitle" hx-post="/update" hx-trigger="blur" hx-target="#subtitle-container" class="text-2xl md:text-3xl font-bold tracking-widest text-pink-200 mt-2" style="font-family: 'Fredoka One', sans-serif;">%s</h2>`, pageData.Subtitle)
}


