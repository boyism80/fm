package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type LieDetectorAnswer struct {
	gs *GameServer
}

func (LieDetectorAnswer) New(gs *GameServer) *LieDetectorAnswer {
	return &LieDetectorAnswer{
		gs: gs,
	}
}

func (h *LieDetectorAnswer) Handle(ctx *core.ClientContext, req *request.LieDetectorAnswer) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	err = character.LieDetector().Answer(ctx.ActorContext, req.Answer)
	if err != nil {
		character.Listener().OnUnlockAction(character)
	}

	return nil
}
