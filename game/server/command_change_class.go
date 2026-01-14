package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
)

type ChangeClass struct {
	gs *GameServer
}

func (*ChangeClass) New(gs *GameServer) *ChangeClass {
	return &ChangeClass{
		gs: gs,
	}
}

func (h *ChangeClass) GetCommandName() string {
	return "직업바꾸기"
}

func (h *ChangeClass) GetUsage() string {
	return "<직업ID> - 직업 변경"
}

func (h *ChangeClass) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing class value")
	}

	class, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid class value: %s", args[0])
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.ChangeClass(uint16(class))

	return nil
}
