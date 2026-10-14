package core

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type PacketHandler struct {
	handlers map[int]func(ctx *ClientContext, data []byte) error
}

type ClientContext struct {
	Client        Client
	LogicActorPID *actor.PID
	ActorContext  actor.Context
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		handlers: make(map[int]func(ctx *ClientContext, data []byte) error),
	}
}

func (h *PacketHandler) RegisterHandler(opcode int, handler func(ctx *ClientContext, data []byte) error) {
	h.handlers[opcode] = handler
	log.Printf("Registered packet handler for opcode 0x%02X", opcode)
}

func (h *PacketHandler) GetHandler(opcode int) func(ctx *ClientContext, data []byte) error {
	return h.handlers[opcode]
}
