package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/game/client"
)

type GetPosition struct {
	gs *GameServer
}

func (*GetPosition) New(gs *GameServer) *GetPosition {
	return &GetPosition{
		gs: gs,
	}
}

func (h *GetPosition) GetCommandName() string {
	return "좌표"
}

func (h *GetPosition) GetUsage() string {
	return "- 현재 좌표 확인"
}

func (h *GetPosition) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	mapID := uint32(0)
	if m := character.GetMap(); m != nil {
		mapID = m.ID
	}
	log.Printf("Command: Character %d position - Map: %d, Position: %v",
		character.GetID(), mapID, character.Position)
	return nil
}
