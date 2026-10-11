package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyRemoveJunior struct {
	familyHandler
}

func (FamilyRemoveJunior) New(gs *GameServer) *FamilyRemoveJunior {
	return &FamilyRemoveJunior{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyRemoveJunior) Handle(ctx *core.ClientContext, req *request.FamilyRemoveJunior) error {
	return h.handle(ctx, &req.Family)
}
