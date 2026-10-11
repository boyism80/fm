package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyRemoveSenior struct {
	familyHandler
}

func (FamilyRemoveSenior) New(gs *GameServer) *FamilyRemoveSenior {
	return &FamilyRemoveSenior{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyRemoveSenior) Handle(ctx *core.ClientContext, req *request.FamilyRemoveSenior) error {
	return h.handle(ctx, &req.Family)
}
