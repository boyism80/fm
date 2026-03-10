package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Aggro struct {
	gs     *GameServer
	opcode byte
}

func (Aggro) New(gs *GameServer) *Aggro {
	return &Aggro{
		gs:     gs,
		opcode: 0x96,
	}
}

func (h *Aggro) GetOpcode() byte {
	return h.opcode
}

func (h *Aggro) Handle(ctx *core.ClientContext, req *request.Aggro) error {
	return nil
}
