package constant

type CarnivalTeam int8

const (
	CarnivalTeamNone CarnivalTeam = -1
	CarnivalTeamRed  CarnivalTeam = 0
	CarnivalTeamBlue CarnivalTeam = 1
)

func AllCarnivalTeams() map[string]CarnivalTeam {
	return map[string]CarnivalTeam{
		"None": CarnivalTeamNone,
		"Blue": CarnivalTeamBlue,
		"Red":  CarnivalTeamRed,
	}
}
