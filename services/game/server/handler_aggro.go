package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Aggro struct {
	gs *GameServer
}

func (Aggro) New(gs *GameServer) *Aggro {
	return &Aggro{
		gs: gs,
	}
}

func (h *Aggro) Handle(ctx *core.ClientContext, req *request.Aggro) error {
	return nil
}
