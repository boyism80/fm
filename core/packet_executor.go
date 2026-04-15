package core

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

func ExecutePacketHandler(actorCtx actor.Context, context ServerContext, client Client, opcode int, data []byte, logicActorPID *actor.PID) error {
	handler := context.GetPacketHandler().GetHandler(opcode)
	if handler == nil {
		return fmt.Errorf("no handler for opcode 0x%02X", opcode)
	}

	ctx := &ClientContext{
		Client:        client,
		Server:        nil,
		LogicActorPID: logicActorPID,
		ActorContext:  actorCtx,
	}

	return handler(ctx, data)
}
