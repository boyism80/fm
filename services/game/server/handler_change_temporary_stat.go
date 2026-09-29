package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type ChangeTemporaryStat struct {
	gs *GameServer
}

func (ChangeTemporaryStat) New(gs *GameServer) *ChangeTemporaryStat {
	return &ChangeTemporaryStat{
		gs: gs,
	}
}

// Empty body. The client sends it after a temporary stat set/reset that touches specific stat flags, and when Berserk (1320006) toggles on its HP threshold.
func (h *ChangeTemporaryStat) Handle(ctx *core.ClientContext, req *request.ChangeTemporaryStat) error {
	// TODO
	return nil
}
