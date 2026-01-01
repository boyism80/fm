package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
)

// LoginGame handles game login packet requests
type LoginGame struct {
	gameServer *GameServer
	opcode     byte
}

func (LoginGame) New(gameServer *GameServer) *LoginGame {
	return &LoginGame{
		gameServer: gameServer,
		opcode:     0x06,
	}
}

func (h *LoginGame) GetOpcode() byte {
	return h.opcode
}

func (h *LoginGame) Handle(ctx *core.ClientContext, req *request.LoginGame) error {
	name := "채승현"
	if req.PlayerId != 1 {
		name = "채진영"
	}

	character := entity.NewDummyCharacter(ctx.Client, nil, req.PlayerId, name, h.gameServer)
	character.Listener = NewGameCharacterListener(h.gameServer, &character)

	// Set GM mode (for testing: player ID 1 is GM)
	if req.PlayerId == 1 {
		character.Admin = true
		character.Invincible = true
	}

	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	client.SetCharacter(&character)

	mapInstance := h.gameServer.GetMap(character.Map)
	if mapInstance == nil {
		log.Printf("Map %d not found (should have been pre-created)", character.Map)
		return fmt.Errorf("map %d not found", character.Map)
	}

	mapSpec := mapInstance.GetSpec()
	if mapSpec == nil {
		log.Printf("MapSpec not found for map %d", character.Map)
		return fmt.Errorf("mapSpec not found for map %d", character.Map)
	}

	portal, ok := mapSpec.Portals[character.SpawnPoint]
	if !ok {
		log.Printf("Portal %d not found in map %d", character.SpawnPoint, character.Map)
		return fmt.Errorf("portal %d not found in map %d", character.SpawnPoint, character.Map)
	}

	character.Position = portal.Position
	character.Stance = 0

	if err := mapInstance.AddPlayer(character.ID, &character, true); err != nil {
		log.Printf("Failed to add player to map: %v", err)
		return err
	}

	return nil
}
