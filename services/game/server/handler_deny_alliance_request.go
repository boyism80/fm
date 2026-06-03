package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type DenyAllianceRequest struct {
	gs *GameServer
}

func (DenyAllianceRequest) New(gs *GameServer) *DenyAllianceRequest {
	return &DenyAllianceRequest{
		gs: gs,
	}
}

func (h *DenyAllianceRequest) Handle(ctx *core.ClientContext, req *request.DenyAllianceRequest) error {
	return nil
}
