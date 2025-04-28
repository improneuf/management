package main

import (
	"os"
	"strings"
)

func GetTeamPhoto(teamName string) string {
	// Normalize team name for comparison
	normalized := strings.ToLower(strings.TrimSpace(teamName))

	// Special cases for teams with known variations
	switch normalized {
	case "2legit2quit", "2legit 2 quit":
		return "2 Legit 2 Quit.png"
	case "far and the rest of the ...", "far and the rest of...", "far and the rest of the family":
		return "Far and the Rest of the Family.png"
	case "livikki":
		return "LIVIKKI.png"
	case "aree and a friend", "aree & a friend":
		return "Aree and a Friend.png"
	case "pompel og tilt", "pompel & tilt":
		return "Pompel & Tilt.png"
	case "one man movie", "1manmovie":
		return "One Man Movie.png"
	case "mixer", "mixer team", "mixerb":
		return "Mixer.png"
	case "showcase", "showcase medlemsworkshop", "showcase norsk medlemsworkshop":
		return "Showcase medlemsworkshop (Eirik).png"
	case "impro neuf ensemble", "impro neuf ensemblet":
		return "The Impro Neuf Ensemble.png"
	case "loose  connections":
		return "Loose Connections.png"
	case "arti kombo":
		return "Arti' Kombo.png"
	}

	files, err := os.ReadDir("team-photos")
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
			// Compare normalized names (without .png extension)
			photoName := strings.ToLower(strings.TrimSuffix(file.Name(), ".png"))
			if photoName == normalized {
				return file.Name()
			}
		}
	}

	return ""
}
