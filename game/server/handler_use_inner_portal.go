package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseInnerPortal struct {
	gs     *GameServer
	opcode byte
}

func (UseInnerPortal) New(gs *GameServer) *UseInnerPortal {
	return &UseInnerPortal{
		gs:     gs,
		opcode: 0x54,
	}
}

func (h *UseInnerPortal) GetOpcode() byte {
	return h.opcode
}

func (h *UseInnerPortal) Handle(ctx *core.ClientContext, req *request.UseInnerPortal) error {
	_ = ctx
	_ = req
	return nil
}
