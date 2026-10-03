package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
)

type HandlePacketHandler struct{}

func (HandlePacketHandler) New() *HandlePacketHandler {
	return &HandlePacketHandler{}
}

func (h *HandlePacketHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *c_actor.HandlePacket) {
	// The character may have changed maps after the packet was queued here; only its current actor may touch its map.
	pid := a.GameWorld.PacketActorPID(msg.Client)
	if pid != nil && pid.Equal(ctx.Self()) == false {
		msg.LogicActorPID = pid
		ctx.Send(pid, msg)
		return
	}

	err := core.ExecutePacketHandler(ctx, a.GameWorld, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}
