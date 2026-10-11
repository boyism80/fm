package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilySummon struct {
	familyHandler
}

func (FamilySummon) New(gs *GameServer) *FamilySummon {
	return &FamilySummon{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilySummon) Handle(ctx *core.ClientContext, req *request.FamilySummon) error {
	return h.handle(ctx, &req.Family)
}
