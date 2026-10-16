package game

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type GuildEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	GuildID    uint32 `json:"guild_id"`
	OccurredAt string `json:"occurred_at"`
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

func (s *GuildSystem) UpdateAsync(ctx actor.Context, evt GuildEventEnvelope) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if s == nil {
		return p
	}
	guildID := evt.GuildID
	p.OnError(func(err error) {
		log.Printf("guild consumer: apply type=%s guild_id=%d: %v", evt.EventType, guildID, err)
	})
	p.ThenRPC(
		func(c context.Context) (*internal.GetGuildReply, error) {
			return s.gw.Systems().Get[*InternalSystem]().client.GetGuild(c, &internal.GetGuildRequest{
				WorldId: s.worldID,
				GuildId: guildID,
			})
		},
		func(reply *internal.GetGuildReply) error {
			if reply == nil || !reply.GetFound() || reply.GetGuild() == nil {
				s.mu.Lock()
				delete(s.guilds, guildID)
				s.mu.Unlock()
				return nil
			}
			s.Update(reply.GetGuild())
			return nil
		},
	)
	return p
}

func (s *GuildSystem) Update(guildPb *internal.Guild) {
	if s == nil || guildPb == nil {
		return
	}
	ent := GuildFromProto(s.gw, guildPb)
	if ent == nil {
		return
	}
	stored := ent.Clone()
	if stored == nil {
		return
	}
	guildID := stored.guildID
	s.mu.Lock()
	defer s.mu.Unlock()
	if prev := s.guilds[guildID]; prev != nil {
		if stored.updatedAt.Before(prev.updatedAt) {
			return
		}
		if _, inAlliance := stored.GetAllianceID(); !inAlliance {
			stored.allianceInvites = prev.CloneAllianceInvites()
		}
	}
	s.guilds[guildID] = stored
	if _, inAlliance := stored.GetAllianceID(); inAlliance {
		stored.ClearAllianceInvites()
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
	if s == nil || characterID == 0 {
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

func (s *GuildSystem) RefreshAsync(ctx actor.Context, guildIDs []uint32) *async.Task {
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if s == nil {
		return promise
	}
	worldID := s.worldID
	for _, guildID := range guildIDs {
		if guildID == 0 {
			continue
		}
		gid := guildID
		promise.ThenRPC(func(c context.Context) (*internal.GetGuildReply, error) {
			return s.gw.Systems().Get[*InternalSystem]().client.GetGuild(c, &internal.GetGuildRequest{
				WorldId: worldID,
				GuildId: gid,
			})
		}, func(reply *internal.GetGuildReply) error {
			if reply != nil && reply.GetFound() && reply.GetGuild() != nil {
				s.Update(reply.GetGuild())
			}
			return nil
		})
	}
	return promise
}
