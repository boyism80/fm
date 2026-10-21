package game

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/system"
	"github.com/boyism80/fm/maple/wz"
)

type ItemWorld interface {
	GetResources() *wz.Resources
	NewUniqueID() uint64
}

type GameWorld interface {
	core.Server
	Systems() *system.Systems
	PacketActorPID(client core.Client) *actor.PID
	GetResources() *wz.Resources
	GetExpRate() int
	GetDropRate() int
	GetMesoRate() int
	DueyIdentityPrompt() bool
	AntiCheat() config.GameAntiCheat
	NewUniqueID() uint64
	GetWorldID() uint32
	GetChannelID() int32
}
