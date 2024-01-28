package main

import (
	"strings"
)

type ShowType string

const (
	Normal        ShowType = "Normal"
	Jam           ShowType = "Jam"
	ClashOfTitans ShowType = "ClashOfTitans"
	StoryNight    ShowType = "StoryNight"
	DuoLab        ShowType = "DuoLab"
	CProject      ShowType = "CProject"
	Other         ShowType = "Other"
)

func getShowType(teamName string) ShowType {
	teamName = strings.ToUpper(teamName)
	switch teamName {
	case "CLASH OF TITANS":
		return ClashOfTitans
	case "DUOLAB":
		return DuoLab
	case "STORY NIGHT":
		return StoryNight
	case "CPROJ":
		return CProject
	case "JAM":
		return Jam
	default:
		return Normal
	}
}

func removeDuplicateShowTypes(slice []ShowType) []ShowType {
	set := make(map[ShowType]struct{})
	var uniqueSlice []ShowType

	for _, element := range slice {
		if _, exists := set[element]; !exists {
			set[element] = struct{}{}
			uniqueSlice = append(uniqueSlice, element)
		}
	}

	return uniqueSlice
}
