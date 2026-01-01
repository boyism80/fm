package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Unknown0C struct {
	gameServer *GameServer
	opcode     byte
}

func (Unknown0C) New(gameServer *GameServer) *Unknown0C {
	return &Unknown0C{
		gameServer: gameServer,
		opcode:     0x0C,
	}
}

func (h *Unknown0C) GetOpcode() byte {
	return h.opcode
}

func (h *Unknown0C) Handle(ctx *core.ClientContext, req *request.Unknown0C) error {
	return nil
}

