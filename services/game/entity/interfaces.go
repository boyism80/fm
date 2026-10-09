package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

type ItemWorld interface {
	GetResources() *wz.Resources
	NewUniqueID() uint64
}

type GameWorld interface {
	core.Server
	PacketActorPID(client core.Client) *actor.PID
	GetResources() *wz.Resources
	GetExpRate() int
	GetDropRate() int
	GetMesoRate() int
	DueyIdentityPrompt() bool
	NewUniqueID() uint64
	BroadcastNotice(messageType constant.ServerMessageType, message string, channel int, ear bool)
	GetWorldID() uint32
	GetChannelID() int32
	SaveAsync(ctx actor.Context, entries []*internal.CharacterSaveEntry) *async.Task
	SaveShopAsync(ctx actor.Context, shop *internal.Shop, entries []*internal.CharacterSaveEntry, storeBank []*internal.Shop, close bool) *async.Promise[*internal.SaveShopReply]
	GetMapSystem() MapSystem
	GetSchedulerSystem() SchedulerSystem
	GetGuildSystem() GuildSystem
	GetAllianceSystem() AllianceSystem
	GetPartySystem() PartySystem
	GetDispatchSystem() DispatchSystem
	GetStateMachineRegistry() StateMachineRegistry
	GetCarnivalRegistry() *CarnivalRegistry
	GetExpeditionRegistry() *ExpeditionRegistry
	StartStateMachineActor(sm *StateMachine) *actor.PID
	SendStateMachineMessage(pid *actor.PID, msg interface{})
	StopStateMachineActor(sm *StateMachine)
	ResumeLua(pid *actor.PID, root *lua.LState, thread *lua.LState, args []lua.LValue)
}

type StateMachineRegistry interface {
	Get(name string) *StateMachineGroup
}
