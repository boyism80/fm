package core

import (
	"fmt"
)

func ExecutePacketHandler(context ServerContext, client Client, opcode int, data []byte, logicActorID string) error {
	handler := context.GetPacketHandler().GetHandler(opcode)
	if handler == nil {
		return fmt.Errorf("no handler for opcode 0x%02X", opcode)
	}

	ctx := &ClientContext{
		Client:       client,
		Server:       nil,
		LogicActorID: logicActorID,
	}

	return handler(ctx, data)
}
