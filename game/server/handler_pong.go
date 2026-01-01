package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

// Pong handles pong packet requests
type Pong struct {
	gameServer *GameServer
	opcode     byte
}

func (Pong) New(gameServer *GameServer) *Pong {
	return &Pong{
		gameServer: gameServer,
		opcode:     0x0A,
	}
}

func (h *Pong) GetOpcode() byte {
	return h.opcode
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	return nil
}
