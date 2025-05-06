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

// Global state (for demonstration purposes) and a mutex for safe concurrent access.
var (
	pageData = PageData{
		PageTitle:           "Dancing with the Dunes: Celebrating Spontaneity",
		BackgroundHazeURL:   "bg_dunes.webp",
		LogoPath:            "logo.png",
		MainTitle:           "Emotional Whirlwinds",
		Subtitle:            "Playing with Big Feelings",
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

// main renders the full page template.
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/update", updateHandler)
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// indexHandler renders the full template.
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
	log.Printf("Update request received: field=%s, value=%s", field, value)

	// Update the global state based on the field.
	switch field {
	case "MainTitle":
		pageData.MainTitle = value
		// Render the updated h1 markup.
		fmt.Fprintf(w, `<h1 contenteditable="true" data-field="MainTitle" onblur="updateField(this)" class="text-5xl md:text-7xl font-bold tracking-widest text-transparent bg-clip-text bg-gradient-to-r from-green-200 via-pink-300 to-yellow-100" style="font-family: 'Fredoka One', sans-serif; padding-bottom: 0.1em;">%s</h1>`, pageData.MainTitle)
	case "Subtitle":
		pageData.Subtitle = value
		fmt.Fprintf(w, `<h2 contenteditable="true" data-field="Subtitle" onblur="updateField(this)" class="text-2xl md:text-3xl font-bold tracking-widest text-pink-200 mt-2" style="font-family: 'Fredoka One', sans-serif;">%s</h2>`, pageData.Subtitle)
	case "HostName1":
		pageData.HostName1 = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="HostName1" onblur="updateField(this)" class="text-3xl font-semibold text-fuchsia-200" style="font-family: 'Fredoka One', sans-serif;">%s</p>`, pageData.HostName1)
	case "HostName2":
		pageData.HostName2 = value
		// In this template HostName2 appears twice (for "and" and for the name). For simplicity, update one instance.
		fmt.Fprintf(w, `<p contenteditable="true" data-field="HostName2" onblur="updateField(this)" class="text-3xl font-semibold text-fuchsia-200" style="font-family: 'Fredoka One', sans-serif;">%s</p>`, pageData.HostName2)
	case "WorkshopDescription":
		pageData.WorkshopDescription = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="WorkshopDescription" onblur="updateField(this)" class="text-lg md:text-2xl text-pink-100" style="font-family: 'Montserrat', sans-serif; font-weight: bolder;">%s</p>`, pageData.WorkshopDescription)
	case "EventDate":
		pageData.EventDate = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="EventDate" onblur="updateField(this)" class="text-xl md:text-3xl text-indigo-200">📅 <span class="font-bold text-yellow-300">Date:</span> %s</p>`, pageData.EventDate)
	case "EventTime":
		pageData.EventTime = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="EventTime" onblur="updateField(this)" class="text-xl md:text-3xl text-purple-200">⏰ <span class="font-bold text-lime-200">Time:</span> %s</p>`, pageData.EventTime)
	case "Location":
		pageData.Location = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="Location" onblur="updateField(this)" class="text-xl md:text-3xl text-teal-200">📍 <span class="font-bold text-yellow-300">Location:</span> %s</p>`, pageData.Location)
	case "Room":
		pageData.Room = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="Room" onblur="updateField(this)" class="text-xl md:text-3xl text-emerald-200">🏢 <span class="font-bold text-pink-100">Room:</span> %s</p>`, pageData.Room)
	case "SignUpPrompt":
		pageData.SignUpPrompt = value
		fmt.Fprintf(w, `<p contenteditable="true" data-field="SignUpPrompt" onblur="updateField(this)" class="text-xl md:text-2xl font-bold text-yellow-300">%s</p>`, pageData.SignUpPrompt)
	default:
		http.Error(w, "Unknown field", http.StatusBadRequest)
		return
	}
}
