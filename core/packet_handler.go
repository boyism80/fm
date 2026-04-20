package core

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type PacketHandler struct {
	handlers map[int]func(ctx *ClientContext, data []byte) error
}

type ClientContext struct {
	Client        Client
	Server        *ServerCore
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

func (h *PacketHandler) Handle(ctx *ClientContext, opcode int, data []byte) error {
	handler, exists := h.handlers[opcode]
	if !exists {
		return fmt.Errorf("no handler registered for opcode 0x%02X", opcode)
	}

	return handler(ctx, data)
}

func (h *PacketHandler) GetHandlerCount() int {
	return len(h.handlers)
}

func (h *PacketHandler) GetRegisteredOpcodes() []int {
	opcodes := make([]int, 0, len(h.handlers))
	for opcode := range h.handlers {
		opcodes = append(opcodes, opcode)
	}
	return opcodes
}
