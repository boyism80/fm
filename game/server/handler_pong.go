package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Pong struct {
	gs     *GameServer
	opcode byte
}

func (Pong) New(gs *GameServer) *Pong {
	return &Pong{
		gs:     gs,
		opcode: 0x0A,
	}
}

func (h *Pong) GetOpcode() byte {
	return h.opcode
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	return nil
}
