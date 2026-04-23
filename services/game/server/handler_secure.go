package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Secure struct {
	gs *GameServer
}

func (Secure) New(gs *GameServer) *Secure {
	return &Secure{
		gs: gs,
	}
}

func (h *Secure) Handle(ctx *core.ClientContext, req *request.Secure) error {
	return nil
}
