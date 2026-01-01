package server

import (
	"fmt"

	"github.com/boyism80/fm/game/client"
)

type Invincible struct {
	gameServer *GameServer
}

func (*Invincible) New(gameServer *GameServer) *Invincible {
	return &Invincible{
		gameServer: gameServer,
	}
}

func (h *Invincible) GetCommandName() string {
	return "무적"
}

func (h *Invincible) GetUsage() string {
	return "- 무적 상태 토글"
}

func (h *Invincible) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Invincible = !character.Invincible
	status := "enabled"
	if !character.Invincible {
		status = "disabled"
	}

	character.Message(fmt.Sprintf("무적 상태: %s", status))

	return nil
}

