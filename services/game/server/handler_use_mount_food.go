package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseMountFood struct {
	gs *GameServer
}

func (UseMountFood) New(gs *GameServer) *UseMountFood {
	return &UseMountFood{
		gs: gs,
	}
}

func (h *UseMountFood) Handle(ctx *core.ClientContext, req *request.UseMountFood) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	// The client locks its actions on use and waits for an unlock even on success.
	character.Mount.Feed(req.Slot, req.ItemID)
	character.Listener.OnUnlockAction(character)
	return nil
}
