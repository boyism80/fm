package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type LieDetectorSkill struct {
	gs *GameServer
}

func (LieDetectorSkill) New(gs *GameServer) *LieDetectorSkill {
	return &LieDetectorSkill{
		gs: gs,
	}
}

func (h *LieDetectorSkill) Handle(ctx *core.ClientContext, req *request.LieDetectorSkill) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	err := character.UseLieDetectorSkill(req.Target)
	if err != nil {
		character.Listener().OnUnlockAction(character)
	}

	return nil
}
