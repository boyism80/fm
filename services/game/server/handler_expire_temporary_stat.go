package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type ExpireTemporaryStat struct {
	gs *GameServer
}

func (ExpireTemporaryStat) New(gs *GameServer) *ExpireTemporaryStat {
	return &ExpireTemporaryStat{
		gs: gs,
	}
}

// Empty body. The client sends it from its field update when a local temporary stat is expiring (throttled to 200ms).
func (h *ExpireTemporaryStat) Handle(ctx *core.ClientContext, req *request.ExpireTemporaryStat) error {
	// TODO
	return nil
}
