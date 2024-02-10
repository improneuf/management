package main

func GetTeamPhoto(teamName string) string {
	switch teamName {
	case "Problemfikserne/Problemfixers":
		return "Problemfikserne.png"
	case "Open DropIn Mixer w/beginners":
		return "Open Drop In Mixer.png"
	default:
		return teamName + ".png"
	}
}
