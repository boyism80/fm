package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Secure struct {
	gs     *GameServer
	opcode byte
}

func (Secure) New(gs *GameServer) *Secure {
	return &Secure{
		gs:     gs,
		opcode: 0x0C,
	}
}

func (h *Secure) GetOpcode() byte {
	return h.opcode
}

func (h *Secure) Handle(ctx *core.ClientContext, req *request.Secure) error {
	return nil
}
