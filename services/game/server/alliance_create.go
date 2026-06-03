package server

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
)

func (gs *GameServer) isGuildMaster(characterID uint32, g *entity.Guild) bool {
	if g == nil || characterID == 0 {
		return false
	}
	if g.LeaderCharacterID == characterID {
		return true
	}
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != characterID {
			continue
		}
		return m.GetRank() == internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER
	}
	return false
}

func (gs *GameServer) PartyPartnerCharacterID(ch *entity.Character) (uint32, bool) {
	if ch == nil || gs == nil || gs.party == nil {
		return 0, false
	}
	partyID := ch.GetPartyID()
	if partyID == nil {
		return 0, false
	}
	party := gs.party.Get(*partyID)
	if party == nil {
		return 0, false
	}
	members := party.GetMembers()
	if len(members) != 2 {
		return 0, false
	}
	selfID := ch.GetID()
	for _, pm := range members {
		if pm == nil || pm.CharacterID == 0 || pm.CharacterID == selfID {
			continue
		}
		return pm.CharacterID, true
	}
	return 0, false
}

func (gs *GameServer) ValidateCreateAlliance(ch *entity.Character) (uint32, bool) {
	if ch == nil || gs == nil || gs.guild == nil || gs.party == nil {
		return 0, false
	}
	leaderID := ch.GetID()
	guildID, inGuild := ch.GetGuildID()
	if !inGuild || guildID == 0 {
		return 0, false
	}
	leaderGuild := gs.guild.Get(guildID)
	if leaderGuild == nil {
		return 0, false
	}
	if !gs.isGuildMaster(leaderID, leaderGuild) {
		return 0, false
	}
	if leaderGuild.AllianceID > 0 {
		return 0, false
	}
	partnerID, ok := gs.PartyPartnerCharacterID(ch)
	if !ok || partnerID == 0 {
		return 0, false
	}
	partnerGuildID, partnerInGuild := gs.guild.GuildIDForCharacter(partnerID)
	if !partnerInGuild || partnerGuildID == 0 || partnerGuildID == guildID {
		return 0, false
	}
	partnerGuild := gs.guild.Get(partnerGuildID)
	if partnerGuild == nil {
		return 0, false
	}
	if !gs.isGuildMaster(partnerID, partnerGuild) {
		return 0, false
	}
	if partnerGuild.AllianceID > 0 {
		return 0, false
	}
	return partnerID, true
}
