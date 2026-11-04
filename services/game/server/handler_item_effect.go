package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type ItemEffect struct {
	gs *GameServer
}

func (ItemEffect) New(gs *GameServer) *ItemEffect { return &ItemEffect{gs: gs} }
func (h *ItemEffect) Handle(ctx *core.ClientContext, req *request.ItemEffect) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	ch.SetItemEffect(req.ItemID)
	return nil
}
