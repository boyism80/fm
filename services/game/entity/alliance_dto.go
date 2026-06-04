package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func AllianceInfoFromProto(a *internal.Alliance) *dto.AllianceInfo {
	if a == nil {
		return nil
	}
	var rankTitles [5]string
	titles := a.GetRankTitles()
	for i := 0; i < 5; i++ {
		if i < len(titles) {
			rankTitles[i] = titles[i]
		}
	}
	guildIDs := append([]uint32(nil), a.GetGuildIds()...)
	return &dto.AllianceInfo{
		AllianceID: a.GetAllianceId(),
		Name:       a.GetName(),
		RankTitles: rankTitles,
		GuildIDs:   guildIDs,
		Capacity:   a.GetCapacity(),
		Notice:     a.GetNotice(),
	}
}

func AllianceCreateGuildsFromProto(a *internal.Alliance) []*dto.GuildInfo {
	if a == nil {
		return nil
	}
	out := make([]*dto.GuildInfo, 0, len(a.GetGuilds()))
	for _, g := range a.GetGuilds() {
		if info := GuildInfoFromProto(g); info != nil {
			out = append(out, info)
		}
	}
	return out
}

func AllianceGuildMemberRanksFromProto(g *internal.Guild) []dto.AllianceGuildMemberRank {
	if g == nil {
		return nil
	}
	out := make([]dto.AllianceGuildMemberRank, 0, len(g.GetMembers()))
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		rank := uint8(0)
		if m.AllianceRank != nil {
			rank = uint8(m.GetAllianceRank())
		}
		out = append(out, dto.AllianceGuildMemberRank{
			CharacterID:  m.GetCharacterId(),
			AllianceRank: rank,
		})
	}
	return out
}

func allianceMembershipChangeGuildFromProto(g *internal.Guild) (dto.AllianceMembershipChangeGuild, bool) {
	if g == nil {
		return dto.AllianceMembershipChangeGuild{}, false
	}
	block := dto.AllianceMembershipChangeGuild{
		GuildID: g.GetGuildId(),
	}
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		if m.AllianceRank == nil {
			continue
		}
		block.Members = append(block.Members, dto.AllianceMembershipChangeMember{
			CharacterID:  m.GetCharacterId(),
			AllianceRank: uint8(m.GetAllianceRank()),
		})
	}
	return block, true
}

func AllianceMembershipChangeGuildFromProto(g *internal.Guild) (dto.AllianceMembershipChangeGuild, bool) {
	return allianceMembershipChangeGuildFromProto(g)
}

func AllianceMembershipChangeGuildsFromProto(a *internal.Alliance) []dto.AllianceMembershipChangeGuild {
	if a == nil {
		return nil
	}
	out := make([]dto.AllianceMembershipChangeGuild, 0, len(a.GetGuilds()))
	for _, g := range a.GetGuilds() {
		block, ok := allianceMembershipChangeGuildFromProto(g)
		if ok {
			out = append(out, block)
		}
	}
	return out
}
