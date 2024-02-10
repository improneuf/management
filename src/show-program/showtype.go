package main

import (
	"strings"
	"time"
)

type ShowType string

const (
	Normal        ShowType = "Normal"
	Jam           ShowType = "Jam"
	ClashOfTitans ShowType = "ClashOfTitans"
	StoryNight    ShowType = "StoryNight"
	DuoLab        ShowType = "DuoLab"
	CProject      ShowType = "CProject"
	Project       ShowType = "Project"
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
	case "PROJ":
		return CProject
	case "CPROJ":
		return Project
	case "JAM":
		return Jam
	default:
		return Normal
	}
}

func getShowNameFromType(showType ShowType, dt time.Time) string {
	switch showType {
	case ClashOfTitans:
		return "Clash of Titans"
	case DuoLab:
		return "DuoLab"
	case StoryNight:
		return "Story Night"
	case Project:
		return "Project"
	case CProject:
		return "CProject"
	case Jam:
		return "Impro Neuf Jam"
	case Normal:
		switch dt.Weekday() {
		case time.Wednesday:
			return "Impro Neuf Wednesday Show"
		case time.Tuesday:
			return "Impro Tirsdagshow"
		}
	}
	return ""
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
