package entity

import "github.com/boyism80/fm/services/game/constant"

func (ch *Character) BindCarnival(match *CarnivalMatch) {
	ch.carnivalMatch = match
}

func (ch *Character) UnbindCarnival() {
	ch.carnivalMatch = nil
}

func (ch *Character) CarnivalMatch() *CarnivalMatch {
	return ch.carnivalMatch
}

func (ch *Character) CarnivalTeamID() constant.CarnivalTeam {
	match := ch.CarnivalMatch()
	if match == nil {
		return constant.CarnivalTeamNone
	}
	return match.TeamIDOf(ch)
}

func (ch *Character) CarnivalTeam() *CarnivalTeam {
	match := ch.CarnivalMatch()
	if match == nil {
		return nil
	}
	return match.TeamOf(ch)
}

func (ch *Character) IsInCarnivalBattleMap() bool {
	match := ch.CarnivalMatch()
	m := ch.GetMap()
	if match == nil || m == nil {
		return false
	}
	mapID := m.GetMapID()
	return mapID == match.FieldMapID || mapID == match.ReviveMapID
}

func (ch *Character) IsOnCarnivalWaitingMap() bool {
	match := ch.CarnivalMatch()
	m := ch.GetMap()
	if match == nil || m == nil {
		return false
	}
	return m.GetMapID() == match.WaitingMapID
}
