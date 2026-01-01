package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/game/client"
)

type GetPosition struct {
	gameServer *GameServer
}

func (*GetPosition) New(gameServer *GameServer) *GetPosition {
	return &GetPosition{
		gameServer: gameServer,
	}
}

func (h *GetPosition) GetCommandName() string {
	return "좌표"
}

func (h *GetPosition) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	log.Printf("Command: Character %d position - Map: %d, Position: %v",
		character.GetID(), character.GetMap(), character.Position)
	return nil
}

