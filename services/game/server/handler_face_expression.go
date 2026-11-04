package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FaceExpression struct {
	gs *GameServer
}

func (FaceExpression) New(gs *GameServer) *FaceExpression { return &FaceExpression{gs: gs} }
func (h *FaceExpression) Handle(ctx *core.ClientContext, req *request.FaceExpression) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	ch.Express(req.Expression)
	return nil
}
