package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Carnival struct {
	gs *GameServer
}

func (Carnival) New(gs *GameServer) *Carnival {
	return &Carnival{gs: gs}
}

func (h *Carnival) Handle(ctx *core.ClientContext, req *request.Carnival) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return err
	}
	defer ch.Listener().OnUnlockAction(ch)

	if ch.Carnival().Match() == nil || ch.StateMachine() == nil {
		return nil
	}
	ch.StateMachine().CallHook("on_carnival_summon", ch, int(req.Tab), req.Num)
	return nil
}
