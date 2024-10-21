package main

import (
	"strings"
	"time"
)

type ShowType string

const (
	ShowTypeRegular       ShowType = "Regular"
	ShowTypeJam           ShowType = "Jam"
	ShowTypeClashOfTitans ShowType = "ClashOfTitans"
	ShowTypeStoryNight    ShowType = "StoryNight"
	ShowTypeDuoLab        ShowType = "DuoLab"
	ShowTypeCProject      ShowType = "CProject"
	Project               ShowType = "Project"
	ShowTypeOther         ShowType = "Other"
	ShowTypeTheme         ShowType = "Theme"
)

func getShowType(teamName string) ShowType {
	teamName = strings.ToUpper(teamName)
	switch teamName {
	case "CLASH OF TITANS":
		return ShowTypeClashOfTitans
	case "DUOLAB":
		return ShowTypeDuoLab
	case "STORYNIGHT":
		return ShowTypeStoryNight
	case "PROJ":
		return ShowTypeCProject
	case "CPROJ":
		return Project
	case "JAM":
		return ShowTypeJam
	case "THEME":
		return ShowTypeTheme
	default:
		return ShowTypeRegular
	}
}

func getShowTitleAndSubtitle(showType ShowType, dt time.Time) (string, string) {
	switch showType {
	case ShowTypeClashOfTitans:
		return "Clash of Titans", ""
	case ShowTypeDuoLab:
		return "DuoLab", ""
	case ShowTypeStoryNight:
		return "Story Night", ""
	case Project:
		return "Project", ""
	case ShowTypeCProject:
		return "CProject", ""
	case ShowTypeJam:
		return "Impro Neuf Jam", ""
	case ShowTypeRegular:
		switch dt.Weekday() {
		case time.Wednesday:
			return "Impro Neuf Wednesday Show", "Laugh, cry, and everything in between"
		case time.Tuesday:
			return "Impro Neuf Tirsdagsshow", ""
		}
	}
	return "", ""
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
