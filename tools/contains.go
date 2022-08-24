package tools

func Contains(wantedTeams []string, teamName string) bool {
	for _, name := range wantedTeams {
		if name == teamName {
			return true
		}
	}
	return false

}
