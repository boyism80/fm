package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeLevel struct {
	gameServer *GameServer
}

func (*ChangeLevel) New(gameServer *GameServer) *ChangeLevel {
	return &ChangeLevel{
		gameServer: gameServer,
	}
}

func (h *ChangeLevel) GetCommandName() string {
	return "레벨바꾸기"
}

func (h *ChangeLevel) GetUsage() string {
	return "<레벨> - 레벨 설정"
}

func (h *ChangeLevel) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing level value")
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid level value: %s", args[0])
	}

	if level < 1 {
		level = 1
	} else if level > 255 {
		level = 255
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Level = uint8(level)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_LEVEL: int32(character.Level),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

