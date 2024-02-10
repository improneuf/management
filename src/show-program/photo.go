package main

import "strings"

func GetTeamPhoto(teamName string) string {
	switch teamName {
	case "Problemfikserne/Problemfixers":
		return "Problemfikserne.png"
	case "Open DropIn Mixer w/beginners":
		return "Open Drop In Mixer.png"
	default:
		if strings.HasPrefix(teamName, "Aree and") {
			return "Aree and a Friend.png"
		}
		return teamName + ".png"
	}
}
