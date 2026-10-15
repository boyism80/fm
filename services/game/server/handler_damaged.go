package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type Damaged struct {
	gs *GameServer
}

func (Damaged) New(gs *GameServer) *Damaged {
	return &Damaged{
		gs: gs,
	}
}

func (h *Damaged) Handle(ctx *core.ClientContext, req *request.Damaged) error {
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

	character.TakeHit(ctx.ActorContext, game.Hit{
		Type:       req.Type,
		Damage:     req.Damage,
		MobOID:     req.OID,
		SkillID:    uint32(req.SkillID),
		SkillLevel: req.Level,
		Reflect:    req.Reflect,
	})
	return nil
}
