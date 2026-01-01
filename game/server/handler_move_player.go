package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
)

// MovePlayer handles player movement packet requests
type MovePlayer struct {
	gameServer *GameServer
	opcode     byte
}

func (MovePlayer) New(gameServer *GameServer) *MovePlayer {
	return &MovePlayer{
		gameServer: gameServer,
		opcode:     0x18,
	}
}

func (h *MovePlayer) GetOpcode() byte {
	return h.opcode
}

func (h *MovePlayer) Handle(ctx *core.ClientContext, req *request.MovePlayer) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	beforePosition := character.Position

	for _, frag := range req.Fragments {
		if move, ok := frag.(*dto.AbsoluteLifeMovement); ok {
			character.Position = move.Position
		}
		character.Stance = frag.GetStance()
	}

	mapInstance := h.gameServer.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found for character movement", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	character.Listener.OnPlayerMove(character.GetMap(), character.GetID(), character, beforePosition, req.Fragments)

	return nil
}
