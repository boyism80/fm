package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyInvite struct {
	familyHandler
}

func (FamilyInvite) New(gs *GameServer) *FamilyInvite {
	return &FamilyInvite{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyInvite) Handle(ctx *core.ClientContext, req *request.FamilyInvite) error {
	return h.handle(ctx, &req.Family)
}
