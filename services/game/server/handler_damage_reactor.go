package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type DamageReactor struct {
	gs *GameServer
}

func (DamageReactor) New(gs *GameServer) *DamageReactor {
	return &DamageReactor{
		gs: gs,
	}
}

func (h *DamageReactor) Handle(ctx *core.ClientContext, req *request.DamageReactor) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		return nil
	}

	reactor := mapInstance.GetReactor(req.OID)
	if reactor == nil {
		return nil
	}

	reactor.Hit(character, req.HitSide, int32(req.Stance))
	return nil
}
