package entity

import (
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) BindCarnival(match *CarnivalMatch) {
	ch.carnivalMatch = match
}

func (ch *Character) UnbindCarnival() {
	ch.carnivalMatch = nil
}

func (ch *Character) CarnivalMatch() *CarnivalMatch {
	return ch.carnivalMatch
}

func (ch *Character) CarnivalTeam() *CarnivalTeam {
	match := ch.CarnivalMatch()
	if match == nil {
		return nil
	}
	return match.FindTeam(ch)
}

func (ch *Character) PartyOnMap() ([]*Character, *Party) {
	partyID := ch.GetPartyID()
	if partyID == nil || ch.GameWorld == nil {
		return nil, nil
	}
	party := ch.GameWorld.GetPartySystem().Get(*partyID)
	leaderMap := ch.GetMap()
	if party == nil || party.GetLeaderCharacterId() != ch.GetID() || leaderMap == nil {
		return nil, nil
	}
	members := make([]*Character, 0)
	for _, mem := range party.GetMembers() {
		if mem == nil {
			continue
		}
		member := leaderMap.GetPlayer(mem.GetCharacterId())
		if member == nil {
			return nil, nil
		}
		members = append(members, member)
	}
	if len(members) == 0 {
		return nil, nil
	}
	return members, party
}

func (ch *Character) PickupCarnivalItem(wzConsume *wz.Consume) bool {
	match := ch.CarnivalMatch()
	if match == nil || wzConsume == nil {
		return false
	}
	team := match.FindTeam(ch)
	if team == nil {
		return false
	}
	if wzConsume.CP > 0 {
		team.AddCP(ch, wzConsume.CP)
	}
	if wzConsume.NuffSkillID > 0 && ch.GameWorld != nil {
		resources := ch.GameWorld.GetResources()
		if resources != nil {
			skillDef := resources.GetCarnivalSkill(wzConsume.NuffSkillID)
			enemy := match.enemyTeam(team.TeamID)
			if skillDef != nil && enemy != nil {
				match.debuffEnemies(ch.GameWorld, resources, enemy, wzConsume.NuffSkillID, skillDef)
			}
		}
	}
	return true
}
