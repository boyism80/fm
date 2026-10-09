package server

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
	"github.com/boyism80/fm/services/game/entity"
)

type GuildEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	GuildID    uint32 `json:"guild_id"`
	OccurredAt string `json:"occurred_at"`
}

type GuildCache struct {
	gs             *GameServer
	worldID        uint32
	internalClient internal.InternalClient
	mu             sync.Mutex
	guilds         map[uint32]*entity.Guild
}

func NewGuildCache(gs *GameServer, worldID uint32, ic internal.InternalClient) *GuildCache {
	return &GuildCache{
		gs:             gs,
		worldID:        worldID,
		internalClient: ic,
		guilds:         make(map[uint32]*entity.Guild),
	}
}

func (gc *GuildCache) UpdateAsync(ctx actor.Context, evt GuildEventEnvelope) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if gc == nil {
		return p
	}
	guildID := evt.GuildID
	p.OnError(func(err error) {
		log.Printf("guild consumer: apply type=%s guild_id=%d: %v", evt.EventType, guildID, err)
	})
	if gc.internalClient == nil {
		return p
	}
	p.ThenRPC(
		func(c context.Context) (*internal.GetGuildReply, error) {
			return gc.internalClient.GetGuild(c, &internal.GetGuildRequest{
				WorldId: gc.worldID,
				GuildId: guildID,
			})
		},
		func(reply *internal.GetGuildReply) error {
			if reply == nil || !reply.GetFound() || reply.GetGuild() == nil {
				gc.mu.Lock()
				delete(gc.guilds, guildID)
				gc.mu.Unlock()
				return nil
			}
			gc.Update(reply.GetGuild())
			return nil
		},
	)
	return p
}

func (gc *GuildCache) Update(guildPb *internal.Guild) {
	if gc == nil || guildPb == nil {
		return
	}
	ent := entity.GuildFromProto(gc.gs, guildPb)
	if ent == nil {
		return
	}
	stored := ent.Clone()
	if stored == nil {
		return
	}
	guildID := stored.GuildID
	gc.mu.Lock()
	defer gc.mu.Unlock()
	if prev := gc.guilds[guildID]; prev != nil {
		if stored.UpdatedAt.Before(prev.UpdatedAt) {
			return
		}
		if _, inAlliance := stored.GetAllianceID(); !inAlliance {
			stored.AllianceInvites = entity.CloneAllianceInvites(prev.AllianceInvites)
		}
	}
	gc.guilds[guildID] = stored
	if _, inAlliance := stored.GetAllianceID(); inAlliance {
		stored.ClearAllianceInvites()
	}
}

func (gc *GuildCache) NameToGuildID(guildName string) (uint32, bool) {
	if gc == nil || guildName == "" {
		return 0, false
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	for _, g := range gc.guilds {
		if g != nil && strings.EqualFold(g.Name, guildName) {
			return g.GuildID, true
		}
	}
	return 0, false
}

func (gc *GuildCache) TrySetAllianceInvite(targetGuildID, allianceID uint32, expiresAt time.Time) bool {
	if gc == nil || targetGuildID == 0 || allianceID == 0 {
		return false
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	g := gc.guilds[targetGuildID]
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

func (gc *GuildCache) PendingAllianceInvite(guildID uint32) (allianceID uint32, expiresAt time.Time, ok bool) {
	if gc == nil || guildID == 0 {
		return 0, time.Time{}, false
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	g := gc.guilds[guildID]
	if g == nil {
		return 0, time.Time{}, false
	}
	return g.PendingAllianceInvite()
}

func (gc *GuildCache) ClearAllianceInvite(guildID, allianceID uint32) {
	if gc == nil || guildID == 0 || allianceID == 0 {
		return
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	g := gc.guilds[guildID]
	if g == nil {
		return
	}
	g.ClearAllianceInvite(allianceID)
}

func (gc *GuildCache) Get(guildID uint32) *entity.Guild {
	if gc == nil {
		return nil
	}
	gc.mu.Lock()
	s := gc.guilds[guildID]
	gc.mu.Unlock()
	if s == nil {
		return nil
	}
	return s.Clone()
}

func (gc *GuildCache) GuildIDForCharacter(characterID uint32) (uint32, bool) {
	if gc == nil || characterID == 0 {
		return 0, false
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	for guildID, g := range gc.guilds {
		if g == nil {
			continue
		}
		for _, m := range g.Members {
			if m != nil && m.CharacterID == characterID {
				return guildID, true
			}
		}
	}
	return 0, false
}

func (gc *GuildCache) RefreshAsync(ctx actor.Context, guildIDs []uint32) *async.Task {
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if gc == nil || gc.internalClient == nil {
		return promise
	}
	worldID := gc.worldID
	for _, guildID := range guildIDs {
		if guildID == 0 {
			continue
		}
		gid := guildID
		promise.ThenRPC(func(c context.Context) (*internal.GetGuildReply, error) {
			return gc.internalClient.GetGuild(c, &internal.GetGuildRequest{
				WorldId: worldID,
				GuildId: gid,
			})
		}, func(reply *internal.GetGuildReply) error {
			if reply != nil && reply.GetFound() && reply.GetGuild() != nil {
				gc.Update(reply.GetGuild())
			}
			return nil
		})
	}
	return promise
}
