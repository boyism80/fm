package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type DenyGuildRequest struct {
	gs *GameServer
}

func (DenyGuildRequest) New(gs *GameServer) *DenyGuildRequest {
	return &DenyGuildRequest{
		gs: gs,
	}
}

func (h *DenyGuildRequest) Handle(ctx *core.ClientContext, req *request.DenyGuildRequest) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	ch.Guild().ClearInvites()
	return nil
}
