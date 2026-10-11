package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyAccept struct {
	familyHandler
}

func (FamilyAccept) New(gs *GameServer) *FamilyAccept {
	return &FamilyAccept{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyAccept) Handle(ctx *core.ClientContext, req *request.FamilyAccept) error {
	return h.handle(ctx, &req.Family)
}
