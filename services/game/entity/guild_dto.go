package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func guildMemberOnline(channelIndex *int32) bool {
	if channelIndex == nil {
		return false
	}
	return *channelIndex >= 0
}

func guildRankToUint32(rank internal.GuildMemberRank) uint32 {
	switch rank {
	case internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER:
		return 1
	case internal.GuildMemberRank_GUILD_MEMBER_RANK_JUNIOR:
		return 2
	case internal.GuildMemberRank_GUILD_MEMBER_RANK_SENIOR:
		return 3
	case internal.GuildMemberRank_GUILD_MEMBER_RANK_MEMBER:
		return 4
	case internal.GuildMemberRank_GUILD_MEMBER_RANK_NEW:
		return 5
	default:
		return 4
	}
}

func GuildMemberToDTO(m *GuildMember) dto.GuildMemberStatus {
	if m == nil {
		return dto.GuildMemberStatus{}
	}
	return dto.GuildMemberStatus{
		CharacterID:  m.CharacterID,
		Name:         m.CharacterName,
		JobID:        m.ClassID,
		Level:        m.Level,
		GuildRank:    guildRankToUint32(m.Rank),
		Online:       guildMemberOnline(m.ChannelIndex),
		AllianceRank: 0,
	}
}

func GuildMembersToDTO(members []*GuildMember) []dto.GuildMemberStatus {
	out := make([]dto.GuildMemberStatus, 0, len(members))
	for _, m := range members {
		if m == nil {
			continue
		}
		out = append(out, GuildMemberToDTO(m))
	}
	return out
}

func GuildToDTO(g *Guild) *dto.GuildInfo {
	if g == nil {
		return nil
	}
	logoBG := uint16(0)
	logoBGColor := uint8(0)
	logo := uint16(0)
	logoColor := uint8(0)
	if g.Logo != nil {
		logoBG = uint16(g.Logo.LogoBG)
		logoBGColor = uint8(g.Logo.LogoBGColor)
		logo = uint16(g.Logo.Logo)
		logoColor = uint8(g.Logo.LogoColor)
	}
	return &dto.GuildInfo{
		GuildID:     g.GuildID,
		Name:        g.Name,
		RankTitles:  g.RankTitles,
		Members:     GuildMembersToDTO(g.Members),
		Capacity:    g.Capacity,
		LogoBG:      logoBG,
		LogoBGColor: logoBGColor,
		Logo:        logo,
		LogoColor:   logoColor,
		Notice:      g.Notice,
		GP:          g.GP,
		AllianceID:  0,
	}
}
