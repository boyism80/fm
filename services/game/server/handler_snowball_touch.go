package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type SnowballTouch struct {
	gs *GameServer
}

func (SnowballTouch) New(gs *GameServer) *SnowballTouch {
	return &SnowballTouch{gs: gs}
}

// Running into a snowball knocks the character back.
func (h *SnowballTouch) Handle(ctx *core.ClientContext, req *request.SnowballTouch) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	if ch.StateMachine() == nil {
		return nil
	}
	ch.StateMachine().CallHook("on_snowball_touch", ch)
	return nil
}
