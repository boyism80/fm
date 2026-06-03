package server

import (
	"context"
	"log"
	"sync"

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
	Revision   uint64 `json:"revision"`
	OccurredAt string `json:"occurred_at"`
}

type GuildContainer struct {
	worldID        uint32
	internalClient internal.InternalClient
	mu             sync.Mutex
	revisions      map[uint32]uint64
	guilds         map[uint32]*entity.Guild
}

func NewGuildContainer(worldID uint32, ic internal.InternalClient) *GuildContainer {
	return &GuildContainer{
		worldID:        worldID,
		internalClient: ic,
		revisions:      make(map[uint32]uint64),
		guilds:         make(map[uint32]*entity.Guild),
	}
}

func (gc *GuildContainer) UpdateAsync(ctx actor.Context, evt GuildEventEnvelope) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if gc == nil {
		return p
	}
	guildID := evt.GuildID
	p.OnError(func(err error) {
		log.Printf("guild consumer: apply type=%s guild_id=%d: %v", evt.EventType, guildID, err)
	})
	if gc.internalClient == nil {
		p.Then(func() (interface{}, error) {
			return nil, nil
		}, func(interface{}) error {
			return nil
		})
		return p
	}
	async.ThenRPC(p,
		func(c context.Context) (*internal.GetGuildReply, error) {
			return gc.internalClient.GetGuild(c, &internal.GetGuildRequest{
				WorldId: gc.worldID,
				GuildId: guildID,
			})
		},
		func(reply *internal.GetGuildReply) error {
			if reply == nil || !reply.GetFound() || reply.GetGuild() == nil {
				gc.mu.Lock()
				delete(gc.revisions, guildID)
				delete(gc.guilds, guildID)
				gc.mu.Unlock()
				return nil
			}
			gc.Update(reply.GetGuild())
			log.Printf("guild consumer: applied type=%s guild_id=%d revision=%d", evt.EventType, guildID, reply.GetGuild().GetRevision())
			return nil
		},
	)
	return p
}

func (gc *GuildContainer) Update(guildPb *internal.Guild) {
	if gc == nil || guildPb == nil {
		return
	}
	ent := entity.GuildFromProto(guildPb)
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
	gc.revisions[guildID] = stored.Revision
	gc.guilds[guildID] = stored
	log.Printf("guild: hydrated guild_id=%d revision=%d", guildID, stored.Revision)
}

func (gc *GuildContainer) Get(guildID uint32) *entity.Guild {
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

func (gc *GuildContainer) GuildIDForCharacter(characterID uint32) (uint32, bool) {
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
