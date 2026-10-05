package entity

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type MapSystem interface {
	Get(mapID uint32) *Map
	GetInstance(instanceKey uint32) *Map
	CreateInstanceMap(templateID uint32, opts MapInitOpts) (*Map, error)
	CreateStateMachineMap(templateID uint32, sm *StateMachine, opts MapInitOpts) (*Map, error)
	SlotInstance(templateID uint32, slot uint32) (*Map, error)
	RemoveInstanceMap(instanceKey uint32) error
	CloseInstance(m *Map)
	Call(m *Map, run func(ctx actor.Context))
	Broadcast(m *Map, message types.Packet)
	Warp(ctx actor.Context, character *Character, targetMap *Map, spawnPoint uint8, onEnter func(actor.Context)) error
	CreateReturnDoor(ch *Character, key DoorKey, skillID constant.SkillID) *Map
	DespawnDoor(m *Map, key DoorKey, animated bool, notifyCounterpart bool)
	SetDoorPartyID(m *Map, key DoorKey, partyID *uint32)
	Reset(L *lua.LState, mapInstance *Map, actorCtx actor.Context) int
	Respawn(L *lua.LState, mapInstance *Map, actorCtx actor.Context, includeNegativeMobTime bool) int
	RunScript(L *lua.LState, actorCtx actor.Context, mapID uint32, scriptPath string, funcName string, args []interface{}) int
	Visit(fn func(*Map))
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
	GainGPAsync(ctx actor.Context, guildID uint32, amount int32) *async.Promise
	SendMessageAsync(ctx actor.Context, guildID uint32, messageType constant.ServerMessageType, message string) *async.Promise
	ShowRankingAsync(ctx actor.Context, ch *Character, npcID uint32) *async.Promise
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
