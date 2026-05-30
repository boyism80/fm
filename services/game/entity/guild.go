package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

type GuildLogo struct {
	Logo        uint32
	LogoColor   uint32
	LogoBG      uint32
	LogoBGColor uint32
}

func (l *GuildLogo) Clone() *GuildLogo {
	if l == nil {
		return nil
	}
	out := *l
	return &out
}

type Guild struct {
	WorldID           uint32
	GuildID           uint32
	Name              string
	LeaderCharacterID uint32
	Revision          uint64
	GP                uint32
	Capacity          uint32
	Notice            string
	Logo              *GuildLogo
	RankTitles        [5]string
	Members           []*GuildMember
}

func (g *Guild) GetWorldId() uint32 {
	return g.WorldID
}

func (g *Guild) GetGuildId() uint32 {
	return g.GuildID
}

func (g *Guild) GetLeaderCharacterId() uint32 {
	return g.LeaderCharacterID
}

func (g *Guild) GetRevision() uint64 {
	return g.Revision
}

func (g *Guild) GetMembers() []*GuildMember {
	return g.Members
}

func (g *Guild) Clone() *Guild {
	if g == nil {
		return nil
	}
	out := &Guild{
		WorldID:           g.WorldID,
		GuildID:           g.GuildID,
		Name:              g.Name,
		LeaderCharacterID: g.LeaderCharacterID,
		Revision:          g.Revision,
		GP:                g.GP,
		Capacity:          g.Capacity,
		Notice:            g.Notice,
		Logo:              g.Logo.Clone(),
		RankTitles:        g.RankTitles,
		Members:           make([]*GuildMember, 0, len(g.Members)),
	}
	for _, m := range g.Members {
		out.Members = append(out.Members, m.Clone())
	}
	return out
}

func GuildMemberCharacterIDs(members []*GuildMember) []uint32 {
	out := make([]uint32, 0, len(members))
	for _, m := range members {
		if m == nil || m.CharacterID == 0 {
			continue
		}
		out = append(out, m.CharacterID)
	}
	return out
}

func GuildLogoFromProto(pb *internal.GuildLogo) *GuildLogo {
	if pb == nil {
		return &GuildLogo{}
	}
	return &GuildLogo{
		Logo:        pb.GetLogo(),
		LogoColor:   pb.GetLogoColor(),
		LogoBG:      pb.GetLogoBg(),
		LogoBGColor: pb.GetLogoBgColor(),
	}
}

func (l *GuildLogo) ToProto() *internal.GuildLogo {
	if l == nil {
		return &internal.GuildLogo{}
	}
	return &internal.GuildLogo{
		Logo:        l.Logo,
		LogoColor:   l.LogoColor,
		LogoBg:      l.LogoBG,
		LogoBgColor: l.LogoBGColor,
	}
}

func GuildFromProto(pb *internal.Guild) *Guild {
	if pb == nil {
		return nil
	}
	members := make([]*GuildMember, 0, len(pb.GetMembers()))
	for _, mm := range pb.GetMembers() {
		if m := GuildMemberFromProto(mm); m != nil {
			members = append(members, m)
		}
	}
	var rankTitles [5]string
	titles := pb.GetRankTitles()
	for i := 0; i < 5; i++ {
		if i < len(titles) {
			rankTitles[i] = titles[i]
		}
	}
	return &Guild{
		WorldID:           pb.GetWorldId(),
		GuildID:           pb.GetGuildId(),
		Name:              pb.GetName(),
		LeaderCharacterID: pb.GetLeaderCharacterId(),
		Revision:          pb.GetRevision(),
		GP:                pb.GetGp(),
		Capacity:          pb.GetCapacity(),
		Notice:            pb.GetNotice(),
		Logo:              GuildLogoFromProto(pb.GetLogo()),
		RankTitles:        rankTitles,
		Members:           members,
	}
}

func (g *Guild) ToProto() *internal.Guild {
	if g == nil {
		return nil
	}
	members := make([]*internal.GuildMember, 0, len(g.Members))
	for _, m := range g.Members {
		if pm := m.ToProto(); pm != nil {
			members = append(members, pm)
		}
	}
	rankTitles := make([]string, 5)
	for i := 0; i < 5; i++ {
		rankTitles[i] = g.RankTitles[i]
	}
	return &internal.Guild{
		WorldId:           g.WorldID,
		GuildId:           g.GuildID,
		Name:              g.Name,
		LeaderCharacterId: g.LeaderCharacterID,
		Revision:          g.Revision,
		Gp:                g.GP,
		Capacity:          g.Capacity,
		Notice:            g.Notice,
		Logo:              g.Logo.ToProto(),
		RankTitles:        rankTitles,
		Members:           members,
	}
}
