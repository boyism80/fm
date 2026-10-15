package game

import "github.com/boyism80/fm/maple/wz"

func (ch *Character) CarnivalMatch() *CarnivalMatch {
	team := ch.CarnivalTeam()
	if team == nil {
		return nil
	}
	return team.Match
}

// The team outlives the match slot so reward NPCs can read the result after Conclude.
func (ch *Character) CarnivalTeam() *CarnivalTeam {
	if ch.carnivalTeam == nil || !ch.carnivalTeam.HasMember(ch.GetID()) {
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
	party := ch.GameWorld.Systems().Get[*PartySystem]().Get(*partyID)
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
