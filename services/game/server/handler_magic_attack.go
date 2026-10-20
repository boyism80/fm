package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type MagicAttack struct {
	gs *GameServer
}

func (MagicAttack) New(gs *GameServer) *MagicAttack {
	return &MagicAttack{
		gs: gs,
	}
}

func (h *MagicAttack) Handle(ctx *core.ClientContext, req *request.MagicAttack) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := gameClient.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	return character.Attack(ctx.ActorContext, game.AttackMagic, req)
}
