package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type NpcClick struct {
	gs *GameServer
}

func (NpcClick) New(gs *GameServer) *NpcClick {
	return &NpcClick{
		gs: gs,
	}
}

func (h *NpcClick) Handle(ctx *core.ClientContext, req *request.NpcClick) error {
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
		log.Printf("Map not found")
		return fmt.Errorf("map not found")
	}

	npc, ok := mapInstance.GetNpcs()[req.OID].(*entity.Npc)
	if !ok {
		log.Printf("NPC %d not found on map", req.OID)
		return fmt.Errorf("npc %d not found", req.OID)
	}

	if err := character.OpenNpc(ctx.ActorContext, npc); err != nil {
		log.Printf("Failed to open NPC %d: %v", npc.Wz.ID, err)
		return err
	}
	return nil
}
