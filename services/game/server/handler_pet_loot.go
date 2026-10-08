package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PetLoot struct {
	gs *GameServer
}

func (PetLoot) New(gs *GameServer) *PetLoot {
	return &PetLoot{
		gs: gs,
	}
}

func (h *PetLoot) Handle(ctx *core.ClientContext, req *request.PetLoot) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil || character.Pets.Active == nil {
		return nil
	}

	err := character.Pets.Active.Loot(req.OID, req.Position)
	if err != nil {
		character.Listener.OnUnlockAction(character)
	}

	return nil
}
