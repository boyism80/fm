package entity

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

type MapSystem interface {
	Get(mapID uint32) *Map
	GetInstance(instanceKey uint32) *Map
	CreateInstanceMap(templateID uint32, opts MapInitOpts) (*Map, error)
	RemoveInstanceMap(instanceKey uint32) error
	Warp(ctx actor.Context, character *Character, targetMap *Map, spawnPoint uint8) error
	CreateReturnDoor(ch *Character, skillID constant.SkillID)
	RemoveReturnDoor(ownerID uint32, skillID uint32, counterpartMapWZID uint32)
	ResetFromLua(L *lua.LState, mapInstance *Map, actorCtx actor.Context) int
	RespawnFromLua(L *lua.LState, mapInstance *Map, actorCtx actor.Context, includeNegativeMobTime bool) int
	RunOnMapFromLua(L *lua.LState, actorCtx actor.Context, mapID uint32, scriptPath string, funcName string, args []interface{}) int
}

type SchedulerSystem interface {
	RunObjectTimer(pid *actor.PID, obj Object, key string)
	RunReactorRespawn(pid *actor.PID, mapID uint32, spawnID uint32)
}

type PartySystem interface {
	Get(partyID uint32) *Party
	UpdateMemberAsync(ctx actor.Context, ch *Character) *async.Promise
}

type GuildSystem interface {
	Get(guildID uint32) *Guild
	TrySetAllianceInvite(guildID, allianceID uint32, expiresAt time.Time) bool
	DisbandAsync(ctx actor.Context, ch *Character, result *int) *async.Promise
	IncCapacityAsync(ctx actor.Context, ch *Character, extendedCap bool, result *int) *async.Promise
	CreateAllianceAsync(ctx actor.Context, ch *Character, allianceName string, result *int) *async.Promise
	DisbandAllianceAsync(ctx actor.Context, ch *Character, result *int) *async.Promise
}

type AllianceSystem interface {
	Get(allianceID uint32) *Alliance
	IncCapacityAsync(ctx actor.Context, ch *Character, result *int) *async.Promise
}

type DispatchSystem interface {
	SendTo(characterID uint32, msg interface{})
}
