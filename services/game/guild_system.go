package game

import (
	"strings"
	"sync"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type GuildEventEnvelope struct {
	EventType string `json:"event_type"`
	WorldID   uint32 `json:"world_id"`
	GuildID   uint32 `json:"guild_id"`
	GuildPB   string `json:"guild_pb"`
}

type GuildSystem struct {
	gw      GameWorld
	worldID uint32
	mu      sync.Mutex
	guilds  map[uint32]*Guild
}

func NewGuildSystem(gw GameWorld) *GuildSystem {
	return &GuildSystem{
		gw:      gw,
		worldID: gw.GetWorldID(),
		guilds:  make(map[uint32]*Guild),
	}
}

func (s *GuildSystem) Start() error { return nil }
func (s *GuildSystem) Stop()        {}

func (s *GuildSystem) Update(guildPb *internal.Guild) bool {
	if guildPb == nil {
		return false
	}
	if s.Get(guildPb.GetGuildId()) == nil {
		found := false
		for _, m := range guildPb.GetMembers() {
			if s.gw.Systems().Get[*DispatchSystem]().Exists(m.GetCharacterId()) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return s.store(guildPb)
}

func (s *GuildSystem) store(guildPb *internal.Guild) bool {
	stored := GuildFromProto(s.gw, guildPb)
	if stored == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if prev := s.guilds[stored.guildID]; prev != nil {
		if stored.revision < prev.revision {
			return false
		}
		if _, inAlliance := stored.GetAllianceID(); !inAlliance {
			stored.allianceInvites = prev.CloneAllianceInvites()
		}
	}
	s.guilds[stored.guildID] = stored
	if _, inAlliance := stored.GetAllianceID(); inAlliance {
		stored.ClearAllianceInvites()
	}
	return true
}

func (s *GuildSystem) Remove(guildID uint32) {
	s.mu.Lock()
	delete(s.guilds, guildID)
	s.mu.Unlock()
}

func (s *GuildSystem) SetMemberChannel(guildID, characterID uint32, channelIndex int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	g := s.guilds[guildID]
	if g == nil {
		return
	}
	if m := g.FindMember(characterID); m != nil {
		m.channelIndex = &channelIndex
	}
}

func (s *GuildSystem) NameToGuildID(guildName string) (uint32, bool) {
	if s == nil || guildName == "" {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, g := range s.guilds {
		if g != nil && strings.EqualFold(g.name, guildName) {
			return g.guildID, true
		}
	}
	return 0, false
}

func (s *GuildSystem) TrySetAllianceInvite(targetGuildID, allianceID uint32, expiresAt time.Time) bool {
	if s == nil || targetGuildID == 0 || allianceID == 0 {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	g := s.guilds[targetGuildID]
	if g == nil {
		return false
	}
	if _, inAlliance := g.GetAllianceID(); inAlliance {
		return false
	}
	if g.HasAllianceInvite() {
		return false
	}
	g.SetAllianceInvite(allianceID, expiresAt)
	return true
}

func (s *GuildSystem) PendingAllianceInvite(guildID uint32) (allianceID uint32, expiresAt time.Time, ok bool) {
	if s == nil || guildID == 0 {
		return 0, time.Time{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	g := s.guilds[guildID]
	if g == nil {
		return 0, time.Time{}, false
	}
	return g.PendingAllianceInvite()
}

func (s *GuildSystem) ClearAllianceInvite(guildID, allianceID uint32) {
	if s == nil || guildID == 0 || allianceID == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	g := s.guilds[guildID]
	if g == nil {
		return
	}
	g.ClearAllianceInvite(allianceID)
}

func (s *GuildSystem) Get(guildID uint32) *Guild {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	guild := s.guilds[guildID]
	s.mu.Unlock()
	if guild == nil {
		return nil
	}
	return guild.Clone()
}

func (s *GuildSystem) GuildIDForCharacter(characterID uint32) (uint32, bool) {
	if s == nil {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for guildID, g := range s.guilds {
		if g == nil {
			continue
		}
		for _, m := range g.members {
			if m != nil && m.characterID == characterID {
				return guildID, true
			}
		}
	}
	return 0, false
}
