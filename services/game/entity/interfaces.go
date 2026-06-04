package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/services/game/wz"
)

type GameWorld interface {
	core.Server
	GetResources() *wz.Resources
	GetExpRate() int
	GetDropRate() int
	GetMesoRate() int
	SaveAsync(ctx actor.Context, chars []*Character) *async.Promise
	GetMapSystem() MapSystem
	GetSchedulerSystem() SchedulerSystem
	GetGuildSystem() GuildSystem
	GetAllianceSystem() AllianceSystem
	GetPartySystem() PartySystem
	GetDispatchSystem() DispatchSystem
}
