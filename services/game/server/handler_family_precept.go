package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyPrecept struct {
	familyHandler
}

func (FamilyPrecept) New(gs *GameServer) *FamilyPrecept {
	return &FamilyPrecept{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyPrecept) Handle(ctx *core.ClientContext, req *request.FamilyPrecept) error {
	return h.handle(ctx, &req.Family)
}
