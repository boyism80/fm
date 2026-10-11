package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyOpen struct {
	familyHandler
}

func (FamilyOpen) New(gs *GameServer) *FamilyOpen {
	return &FamilyOpen{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyOpen) Handle(ctx *core.ClientContext, req *request.FamilyOpen) error {
	return h.handle(ctx, &req.Family)
}
