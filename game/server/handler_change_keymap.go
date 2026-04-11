package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type ChangeKeymap struct {
	gs     *GameServer
	opcode byte
}

func (ChangeKeymap) New(gs *GameServer) *ChangeKeymap {
	return &ChangeKeymap{
		gs:     gs,
		opcode: 0x71,
	}
}

func (h *ChangeKeymap) GetOpcode() byte {
	return h.opcode
}

func (h *ChangeKeymap) Handle(ctx *core.ClientContext, req *request.ChangeKeymap) error {
	_ = ctx
	_ = req
	return nil
}
