package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyUse struct {
	familyHandler
}

func (FamilyUse) New(gs *GameServer) *FamilyUse {
	return &FamilyUse{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyUse) Handle(ctx *core.ClientContext, req *request.FamilyUse) error {
	return h.handle(ctx, &req.Family)
}
