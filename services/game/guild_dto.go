package game

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

func (m *GuildMember) RankWire() uint8 {
	return uint8(guildRankToUint32(m.rank))
}

func (m *GuildMember) BuildS2C() dto.GuildMemberStatus {
	return dto.GuildMemberStatus{
		CharacterID:  m.characterID,
		Name:         m.characterName,
		ClassID:      m.classID,
		Level:        m.level,
		GuildRank:    guildRankToUint32(m.rank),
		Online:       guildMemberOnline(m.channelIndex),
		AllianceRank: cloneU32Ptr(m.allianceRank),
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

func (g *Guild) BuildS2C() *dto.GuildInfo {
	logoBG := uint16(0)
	logoBGColor := uint8(0)
	logo := uint16(0)
	logoColor := uint8(0)
	if g.logo != nil {
		logoBG = uint16(g.logo.logoBG)
		logoBGColor = uint8(g.logo.logoBGColor)
		logo = uint16(g.logo.logo)
		logoColor = uint8(g.logo.logoColor)
	}
	members := make([]dto.GuildMemberStatus, 0, len(g.members))
	for _, m := range g.members {
		if m == nil {
			continue
		}
		members = append(members, m.BuildS2C())
	}
	return &dto.GuildInfo{
		GuildID:     g.guildID,
		Name:        g.name,
		RankTitles:  g.rankTitles,
		Members:     members,
		Capacity:    g.capacity,
		LogoBG:      logoBG,
		LogoBGColor: logoBGColor,
		Logo:        logo,
		LogoColor:   logoColor,
		Notice:      g.notice,
		GP:          g.gp,
		AllianceID:  cloneU32Ptr(g.allianceID),
	}
}
