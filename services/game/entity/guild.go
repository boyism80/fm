package entity

import (
	"github.com/boyism80/fm/core/clock"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

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
	GameWorld         GameWorld
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
	AllianceID        *uint32
	AllianceInvites   map[uint32]time.Time
	Members           []*GuildMember
}

func CloneAllianceInvites(src map[uint32]time.Time) map[uint32]time.Time {
	if len(src) == 0 {
		return nil
	}
	out := make(map[uint32]time.Time, len(src))
	for allianceID, expiresAt := range src {
		out[allianceID] = expiresAt
	}
	return out
}

func (g *Guild) PruneAllianceInvites() {
	if g == nil || len(g.AllianceInvites) == 0 {
		return
	}
	now := clock.Now()
	for allianceID, expiresAt := range g.AllianceInvites {
		if !now.Before(expiresAt) {
			delete(g.AllianceInvites, allianceID)
		}
	}
}

func (g *Guild) HasAllianceInvite() bool {
	if g == nil {
		return false
	}
	g.PruneAllianceInvites()
	return len(g.AllianceInvites) > 0
}

func (g *Guild) SetAllianceInvite(allianceID uint32, expiresAt time.Time) {
	if g == nil || allianceID == 0 {
		return
	}
	if g.AllianceInvites == nil {
		g.AllianceInvites = make(map[uint32]time.Time)
	}
	g.AllianceInvites[allianceID] = expiresAt
}

func (g *Guild) PendingAllianceInvite() (allianceID uint32, expiresAt time.Time, ok bool) {
	if g == nil {
		return 0, time.Time{}, false
	}
	g.PruneAllianceInvites()
	for id, exp := range g.AllianceInvites {
		return id, exp, true
	}
	return 0, time.Time{}, false
}

func (g *Guild) ClearAllianceInvite(allianceID uint32) {
	if g == nil || allianceID == 0 {
		return
	}
	delete(g.AllianceInvites, allianceID)
}

func (g *Guild) ClearAllianceInvites() {
	if g == nil {
		return
	}
	for allianceID := range g.AllianceInvites {
		delete(g.AllianceInvites, allianceID)
	}
}

func (g *Guild) GetAllianceID() (uint32, bool) {
	if g == nil || g.AllianceID == nil {
		return 0, false
	}
	return *g.AllianceID, true
}

func (g *Guild) Alliance() *Alliance {
	if g == nil {
		return nil
	}
	allianceID, ok := g.GetAllianceID()
	if !ok {
		return nil
	}
	return g.GameWorld.GetAllianceSystem().Get(allianceID)
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

func (g *Guild) FindMember(characterID uint32) *GuildMember {
	if g == nil {
		return nil
	}
	for _, m := range g.Members {
		if m == nil || m.CharacterID != characterID {
			continue
		}
		return m
	}
	return nil
}

func (g *Guild) IsGuildMaster(characterID uint32) bool {
	if g == nil || characterID == 0 {
		return false
	}
	if g.LeaderCharacterID == characterID {
		return true
	}
	m := g.FindMember(characterID)
	if m == nil {
		return false
	}
	return m.Rank == internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER
}

func (g *Guild) Clone() *Guild {
	if g == nil {
		return nil
	}
	out := &Guild{
		GameWorld:         g.GameWorld,
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
	if g.AllianceID != nil {
		id := *g.AllianceID
		out.AllianceID = &id
	}
	out.AllianceInvites = CloneAllianceInvites(g.AllianceInvites)
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

func GuildFromProto(gw GameWorld, pb *internal.Guild) *Guild {
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
	guild := &Guild{
		GameWorld:         gw,
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
		AllianceInvites:   make(map[uint32]time.Time),
	}
	if pb.AllianceId != nil {
		id := pb.GetAllianceId()
		guild.AllianceID = &id
	}
	return guild
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
	pb := &internal.Guild{
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
	if g.AllianceID != nil {
		id := *g.AllianceID
		pb.AllianceId = &id
	}
	return pb
}
