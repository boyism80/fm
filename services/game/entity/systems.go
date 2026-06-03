package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/services/game/constant"
)

type MapSystem interface {
	Get(mapID uint32) *Map
	Warp(character *Character, targetMap *Map, spawnPoint uint8) error
	CreateReturnDoor(ch *Character, skillID constant.SkillID)
	RemoveReturnDoor(ownerID uint32, skillID uint32, counterpartMapWZID uint32)
}

type SchedulerSystem interface {
	RunObjectTimer(pid *actor.PID, obj Object, key string)
}

type PartySystem interface {
	Get(partyID uint32) *Party
}

type GuildSystem interface {
	Get(guildID uint32) *Guild
	DisbandAsync(ctx actor.Context, ch *Character, result *int) *async.Promise
	IncCapacityAsync(ctx actor.Context, ch *Character, extendedCap bool, result *int) *async.Promise
	CreateAllianceAsync(ctx actor.Context, ch *Character, allianceName string, result *int) *async.Promise
}

type DispatchSystem interface {
	SendTo(characterID uint32, msg interface{})
}
