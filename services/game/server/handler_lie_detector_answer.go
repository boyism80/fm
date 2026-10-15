package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
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
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	err := character.LieDetector.Answer(ctx.ActorContext, req.Answer)
	if err != nil {
		character.Listener.OnUnlockAction(character)
	}

	return nil
}
