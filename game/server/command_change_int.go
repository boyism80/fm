package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeInt struct {
	gameServer *GameServer
}

func (*ChangeInt) New(gameServer *GameServer) *ChangeInt {
	return &ChangeInt{
		gameServer: gameServer,
	}
}

func (h *ChangeInt) GetCommandName() string {
	return "지능바꾸기"
}

func (h *ChangeInt) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing INT value")
	}

	intVal, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid INT value: %s", args[0])
	}

	if intVal > 32767 {
		intVal = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Int = uint16(intVal)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_INT: int32(character.Int),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

