package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) BindCarnival(team *CarnivalTeam) {
	ch.carnivalTeam = team
}

func (ch *Character) UnbindCarnival() {
	ch.carnivalTeam = nil
}

func (ch *Character) CarnivalMatch() *CarnivalMatch {
	team := ch.CarnivalTeam()
	if team == nil {
		return nil
	}
	return team.Match
}

// The team outlives the match slot so reward NPCs can read the result after Conclude.
func (ch *Character) CarnivalTeam() *CarnivalTeam {
	if ch.carnivalTeam == nil || ch.carnivalTeam.HasMember(ch.GetID()) == false {
		return nil
	}
	return ch.carnivalTeam
}

func (ch *Character) CarnivalMember() *CarnivalMember {
	return &CarnivalMember{
		ID:    ch.GetID(),
		Name:  ch.GetName(),
		Level: ch.GetLevel(),
		Class: ch.Class,
	}
}

func (ch *Character) PartyOnMap() ([]*Character, *Party) {
	partyID := ch.Party.ID()
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

func (ch *Character) useCarnivalItem(wzConsume *wz.Consume) {
	team := ch.CarnivalTeam()
	if team == nil || team.Match == nil {
		return
	}
	if wzConsume.CP > 0 {
		team.AddCP(ch, wzConsume.CP)
	}
	if wzConsume.NuffSkillID == 0 {
		return
	}
	enemy := team.Match.enemyTeam(team.TeamID)
	if enemy == nil {
		return
	}
	enemy.Debuff(ch.GetMap(), wzConsume.NuffSkillID)
}

// Returns false when the mob skill is not a disease, so the caller can fall back to a dispel.
func (ch *Character) GiveMobSkillDebuff(skillID uint32, skillLevel uint8) bool {
	mobSkill := ch.GameWorld.GetResources().GetMobSkill(skillID, skillLevel)
	if mobSkill == nil {
		return false
	}
	for _, flag := range constant.AllDebuffFlags() {
		if flag.DiseaseSkillID != uint16(skillID) {
			continue
		}
		duration := time.Duration(mobSkill.DurationMs) * time.Millisecond
		if duration <= 0 {
			duration = 5 * time.Second
		}
		ch.Debuffs.Give(flag, duration, int16(mobSkill.X), uint16(skillID), uint16(skillLevel))
		return true
	}
	return false
}
