package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type SwitchChannel struct {
	gs *GameServer
}

func (SwitchChannel) New(gs *GameServer) *SwitchChannel {
	return &SwitchChannel{
		gs: gs,
	}
}

func (h *SwitchChannel) Handle(ctx *core.ClientContext, req *request.SwitchChannel) error {
	if h.gs == nil || req == nil {
		return fmt.Errorf("switch channel: invalid state")
	}
	log.Printf("SwitchChannel recv: channel_id=%d (0-based) from %s",
		req.Channel, ctx.Client.GetConnection().RemoteAddr())
	return nil
}
