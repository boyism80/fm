package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetReviveInquiry struct {
	gs *GameServer
}

func (PetReviveInquiry) New(gs *GameServer) *PetReviveInquiry {
	return &PetReviveInquiry{
		gs: gs,
	}
}

func (h *PetReviveInquiry) Handle(ctx *core.ClientContext, req *request.PetReviveInquiry) error {
	// TODO: open the cash shop revive flow once the cash shop exists
	return nil
}
