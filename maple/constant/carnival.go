package constant

type CarnivalTeam int8

const (
	CarnivalTeamNone CarnivalTeam = -1
	CarnivalTeamRed  CarnivalTeam = 0
	CarnivalTeamBlue CarnivalTeam = 1
)

type CarnivalTab uint8

const (
	CarnivalTabMob      CarnivalTab = 0
	CarnivalTabSkill    CarnivalTab = 1
	CarnivalTabGuardian CarnivalTab = 2
)
