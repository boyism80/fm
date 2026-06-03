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
	guildIDs := make([]uint32, 0, len(a.GetGuildIds()))
	for _, id := range a.GetGuildIds() {
		if id > 0 {
			guildIDs = append(guildIDs, id)
		}
	}
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
		ent := GuildFromProto(g)
		if info := GuildToDTO(ent); info != nil {
			out = append(out, info)
		}
	}
	return out
}

func AllianceMembershipChangeGuildsFromProto(a *internal.Alliance) []dto.AllianceMembershipChangeGuild {
	if a == nil {
		return nil
	}
	out := make([]dto.AllianceMembershipChangeGuild, 0, len(a.GetGuilds()))
	for _, g := range a.GetGuilds() {
		if g == nil {
			continue
		}
		ent := GuildFromProto(g)
		if ent == nil {
			continue
		}
		block := dto.AllianceMembershipChangeGuild{
			GuildID: ent.GuildID,
		}
		for _, m := range ent.Members {
			if m == nil || m.CharacterID == 0 {
				continue
			}
			block.Members = append(block.Members, dto.AllianceMembershipChangeMember{
				CharacterID:  m.CharacterID,
				AllianceRank: uint8(m.AllianceRank),
			})
		}
		out = append(out, block)
	}
	return out
}
