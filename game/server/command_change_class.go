package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeClass struct {
	gameServer *GameServer
}

func (*ChangeClass) New(gameServer *GameServer) *ChangeClass {
	return &ChangeClass{
		gameServer: gameServer,
	}
}

func (h *ChangeClass) GetCommandName() string {
	return "직업바꾸기"
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

	character.Class = uint16(class)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_JOB: int32(character.Class),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

