package game

import (
	"encoding/json"
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type GuildLogo struct {
	logo        uint32
	logoColor   uint32
	logoBG      uint32
	logoBGColor uint32
}

func (l *GuildLogo) Clone() *GuildLogo {
	if l == nil {
		return nil
	}
	out := *l
	return &out
}

type Guild struct {
	gameWorld         GameWorld
	worldID           uint32
	guildID           uint32
	name              string
	leaderCharacterID uint32
	revision          uint64
	gp                uint32
	capacity          uint32
	notice            string
	logo              *GuildLogo
	rankTitles        [5]string
	allianceID        *uint32
	allianceInvites   map[uint32]time.Time
	members           []*GuildMember
}

func (g *Guild) CloneAllianceInvites() map[uint32]time.Time {
	if len(g.allianceInvites) == 0 {
		return nil
	}
	out := make(map[uint32]time.Time, len(g.allianceInvites))
	for allianceID, expiresAt := range g.allianceInvites {
		out[allianceID] = expiresAt
	}
	return out
}

func (g *Guild) PruneAllianceInvites() {
	if g == nil || len(g.allianceInvites) == 0 {
		return
	}
	now := clock.Now()
	for allianceID, expiresAt := range g.allianceInvites {
		if !now.Before(expiresAt) {
			delete(g.allianceInvites, allianceID)
		}
	}
}

func (g *Guild) HasAllianceInvite() bool {
	if g == nil {
		return false
	}
	g.PruneAllianceInvites()
	return len(g.allianceInvites) > 0
}

func (g *Guild) SetAllianceInvite(allianceID uint32, expiresAt time.Time) {
	if g == nil || allianceID == 0 {
		return
	}
	if g.allianceInvites == nil {
		g.allianceInvites = make(map[uint32]time.Time)
	}
	g.allianceInvites[allianceID] = expiresAt
}

func (g *Guild) PendingAllianceInvite() (allianceID uint32, expiresAt time.Time, ok bool) {
	if g == nil {
		return 0, time.Time{}, false
	}
	g.PruneAllianceInvites()
	for id, exp := range g.allianceInvites {
		return id, exp, true
	}
	return 0, time.Time{}, false
}

func (g *Guild) ClearAllianceInvite(allianceID uint32) {
	if g == nil || allianceID == 0 {
		return
	}
	delete(g.allianceInvites, allianceID)
}

func (g *Guild) ClearAllianceInvites() {
	if g == nil {
		return
	}
	for allianceID := range g.allianceInvites {
		delete(g.allianceInvites, allianceID)
	}
}

func (g *Guild) GetAllianceID() (uint32, bool) {
	if g == nil || g.allianceID == nil {
		return 0, false
	}
	return *g.allianceID, true
}

func (g *Guild) Alliance() *Alliance {
	if g == nil {
		return nil
	}
	allianceID, ok := g.GetAllianceID()
	if !ok {
		return nil
	}
	return g.gameWorld.Systems().Get[*AllianceSystem]().Get(allianceID)
}

func (g *Guild) GetWorldId() uint32 {
	return g.worldID
}

func (g *Guild) GetGuildId() uint32 {
	return g.guildID
}

func (g *Guild) GetLeaderCharacterId() uint32 {
	return g.leaderCharacterID
}

func (g *Guild) Level() uint32 {
	for level := 1; level < len(constant.GuildLevelGP); level++ {
		if g.gp < constant.GuildLevelGP[level] {
			return uint32(level)
		}
	}
	return uint32(len(constant.GuildLevelGP))
}

func (g *Guild) GetMembers() []*GuildMember {
	return g.members
}

func (g *Guild) FindMember(characterID uint32) *GuildMember {
	if g == nil {
		return nil
	}
	for _, m := range g.members {
		if m == nil || m.characterID != characterID {
			continue
		}
		return m
	}
	return nil
}

func (g *Guild) FindMaster(characterID uint32) *GuildMember {
	if g == nil || g.leaderCharacterID != characterID {
		return nil
	}
	m := g.FindMember(characterID)
	if m == nil || m.rank != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	return m
}

func (g *Guild) IsAllianceLeader(characterID uint32) bool {
	rank, ok := g.FindMaster(characterID).GetAllianceRank()
	return ok && rank == constant.AllianceRankLeader
}

func (g *Guild) IsAllianceOfficer(characterID uint32) bool {
	rank, ok := g.FindMaster(characterID).GetAllianceRank()
	return ok && rank <= constant.AllianceRankOfficer
}

func (g *Guild) IsGuildMaster(characterID uint32) bool {
	if g == nil {
		return false
	}
	if g.leaderCharacterID == characterID {
		return true
	}
	m := g.FindMember(characterID)
	if m == nil {
		return false
	}
	return m.rank == internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER
}

func (g *Guild) CanInvite(characterID uint32) bool {
	if g == nil {
		return false
	}
	m := g.FindMember(characterID)
	if m == nil {
		return false
	}
	switch m.rank {
	case internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER,
		internal.GuildMemberRank_GUILD_MEMBER_RANK_JUNIOR:
		return true
	default:
		return false
	}
}

func (g *Guild) Clone() *Guild {
	if g == nil {
		return nil
	}
	out := &Guild{
		gameWorld:         g.gameWorld,
		worldID:           g.worldID,
		guildID:           g.guildID,
		name:              g.name,
		leaderCharacterID: g.leaderCharacterID,
		revision:          g.revision,
		gp:                g.gp,
		capacity:          g.capacity,
		notice:            g.notice,
		logo:              g.logo.Clone(),
		rankTitles:        g.rankTitles,
		members:           make([]*GuildMember, 0, len(g.members)),
	}
	if g.allianceID != nil {
		id := *g.allianceID
		out.allianceID = &id
	}
	out.allianceInvites = g.CloneAllianceInvites()
	for _, m := range g.members {
		out.members = append(out.members, m.Clone())
	}
	return out
}

func (g *Guild) MemberCharacterIDs() []uint32 {
	out := make([]uint32, 0, len(g.members))
	for _, m := range g.members {
		if m == nil {
			continue
		}
		out = append(out, m.characterID)
	}
	return out
}

func NewGuildLogo(logo, logoColor, logoBG, logoBGColor uint32) *GuildLogo {
	return &GuildLogo{logo: logo, logoColor: logoColor, logoBG: logoBG, logoBGColor: logoBGColor}
}

func (l *GuildLogo) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("null"), nil
	}
	return json.Marshal(struct {
		Logo        uint32
		LogoColor   uint32
		LogoBg      uint32
		LogoBgColor uint32
	}{l.logo, l.logoColor, l.logoBG, l.logoBGColor})
}

func (l *GuildLogo) UnmarshalJSON(data []byte) error {
	var logo struct {
		Logo        uint32
		LogoColor   uint32
		LogoBg      uint32
		LogoBgColor uint32
	}
	if err := json.Unmarshal(data, &logo); err != nil {
		return err
	}
	l.logo = logo.Logo
	l.logoColor = logo.LogoColor
	l.logoBG = logo.LogoBg
	l.logoBGColor = logo.LogoBgColor
	return nil
}

func GuildLogoFromProto(pb *internal.GuildLogo) *GuildLogo {
	if pb == nil {
		return &GuildLogo{}
	}
	return &GuildLogo{
		logo:        pb.GetLogo(),
		logoColor:   pb.GetLogoColor(),
		logoBG:      pb.GetLogoBg(),
		logoBGColor: pb.GetLogoBgColor(),
	}
}

func (l *GuildLogo) BuildRPC() *internal.GuildLogo {
	if l == nil {
		return &internal.GuildLogo{}
	}
	return &internal.GuildLogo{
		Logo:        l.logo,
		LogoColor:   l.logoColor,
		LogoBg:      l.logoBG,
		LogoBgColor: l.logoBGColor,
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
		gameWorld:         gw,
		worldID:           pb.GetWorldId(),
		guildID:           pb.GetGuildId(),
		name:              pb.GetName(),
		leaderCharacterID: pb.GetLeaderCharacterId(),
		revision:          pb.GetRevision(),
		gp:                pb.GetGp(),
		capacity:          pb.GetCapacity(),
		notice:            pb.GetNotice(),
		logo:              GuildLogoFromProto(pb.GetLogo()),
		rankTitles:        rankTitles,
		members:           members,
		allianceInvites:   make(map[uint32]time.Time),
	}
	if pb.AllianceId != nil {
		id := pb.GetAllianceId()
		guild.allianceID = &id
	}
	return guild
}

func (g *Guild) BuildRPC() *internal.Guild {
	if g == nil {
		return nil
	}
	members := make([]*internal.GuildMember, 0, len(g.members))
	for _, m := range g.members {
		if pm := m.BuildRPC(); pm != nil {
			members = append(members, pm)
		}
	}
	rankTitles := make([]string, 5)
	for i := 0; i < 5; i++ {
		rankTitles[i] = g.rankTitles[i]
	}
	pb := &internal.Guild{
		WorldId:           g.worldID,
		GuildId:           g.guildID,
		Name:              g.name,
		LeaderCharacterId: g.leaderCharacterID,
		Revision:          g.revision,
		Gp:                g.gp,
		Capacity:          g.capacity,
		Notice:            g.notice,
		Logo:              g.logo.BuildRPC(),
		RankTitles:        rankTitles,
		Members:           members,
	}
	if g.allianceID != nil {
		id := *g.allianceID
		pb.AllianceId = &id
	}
	return pb
}
