package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

// NpcClick handles NPC click packet requests
type NpcClick struct {
	gameServer *GameServer
	opcode     byte
}

func (NpcClick) New(gameServer *GameServer) *NpcClick {
	return &NpcClick{
		gameServer: gameServer,
		opcode:     0x29,
	}
}

func (h *NpcClick) GetOpcode() byte {
	return h.opcode
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

	mapInstance := h.gameServer.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	npcs := mapInstance.GetNpcs()
	npc, exists := npcs[req.OID]
	if !exists {
		log.Printf("NPC %d not found on map %d", req.OID, character.GetMap())
		return fmt.Errorf("npc %d not found", req.OID)
	}

	if err := h.gameServer.ExecuteNpcScript(character, npc); err != nil {
		log.Printf("Failed to execute NPC script: %v", err)
		return err
	}

	return nil
}
