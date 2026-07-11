package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func cloneU32Ptr(p *uint32) *uint32 {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

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

func GuildMemberRankWire(rank internal.GuildMemberRank) uint8 {
	return uint8(guildRankToUint32(rank))
}

func GuildMemberToDTO(m *GuildMember) dto.GuildMemberStatus {
	if m == nil {
		return dto.GuildMemberStatus{}
	}
	return dto.GuildMemberStatus{
		CharacterID:  m.CharacterID,
		Name:         m.CharacterName,
		ClassID:      m.ClassID,
		Level:        m.Level,
		GuildRank:    guildRankToUint32(m.Rank),
		Online:       guildMemberOnline(m.ChannelIndex),
		AllianceRank: cloneU32Ptr(m.AllianceRank),
	}
}

func GuildMemberStatusFromProto(m *internal.GuildMember) dto.GuildMemberStatus {
	if m == nil {
		return dto.GuildMemberStatus{}
	}
	online := false
	if m.ChannelIndex != nil && *m.ChannelIndex >= 0 {
		online = true
	}
	var allianceRank *uint32
	if m.AllianceRank != nil {
		ar := m.GetAllianceRank()
		allianceRank = &ar
	}
	return dto.GuildMemberStatus{
		CharacterID:  m.GetCharacterId(),
		Name:         m.GetCharacterName(),
		ClassID:      m.GetClassId(),
		Level:        m.GetLevel(),
		GuildRank:    guildRankToUint32(m.GetRank()),
		Online:       online,
		AllianceRank: allianceRank,
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

func GuildInfoFromProto(pb *internal.Guild) *dto.GuildInfo {
	if pb == nil {
		return nil
	}
	logoBG := uint16(0)
	logoBGColor := uint8(0)
	logo := uint16(0)
	logoColor := uint8(0)
	if logoPb := pb.GetLogo(); logoPb != nil {
		logoBG = uint16(logoPb.GetLogoBg())
		logoBGColor = uint8(logoPb.GetLogoBgColor())
		logo = uint16(logoPb.GetLogo())
		logoColor = uint8(logoPb.GetLogoColor())
	}
	var rankTitles [5]string
	titles := pb.GetRankTitles()
	for i := 0; i < 5; i++ {
		if i < len(titles) {
			rankTitles[i] = titles[i]
		}
	}
	members := make([]dto.GuildMemberStatus, 0, len(pb.GetMembers()))
	for _, m := range pb.GetMembers() {
		members = append(members, GuildMemberStatusFromProto(m))
	}
	var allianceID *uint32
	if pb.AllianceId != nil {
		id := pb.GetAllianceId()
		allianceID = &id
	}
	return &dto.GuildInfo{
		GuildID:     pb.GetGuildId(),
		Name:        pb.GetName(),
		RankTitles:  rankTitles,
		Members:     members,
		Capacity:    pb.GetCapacity(),
		LogoBG:      logoBG,
		LogoBGColor: logoBGColor,
		Logo:        logo,
		LogoColor:   logoColor,
		Notice:      pb.GetNotice(),
		GP:          pb.GetGp(),
		AllianceID:  allianceID,
	}
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
		AllianceID:  cloneU32Ptr(g.AllianceID),
	}
}
