package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyPedigree struct {
	familyHandler
}

func (FamilyPedigree) New(gs *GameServer) *FamilyPedigree {
	return &FamilyPedigree{
		familyHandler: familyHandler{gs: gs},
	}
}

func (h *FamilyPedigree) Handle(ctx *core.ClientContext, req *request.FamilyPedigree) error {
	return h.handle(ctx, &req.Family)
}
