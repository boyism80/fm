package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeAllStats struct {
	gs *GameServer
}

func (*ChangeAllStats) New(gs *GameServer) *ChangeAllStats {
	return &ChangeAllStats{
		gs: gs,
	}
}

func (h *ChangeAllStats) GetCommandName() string {
	return "모든스탯바꾸기"
}

func (h *ChangeAllStats) GetUsage() string {
	return "<스탯값> - 모든 스탯 설정"
}

func (h *ChangeAllStats) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing stat value")
	}

	statValue, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid stat value: %s", args[0])
	}

	if statValue > 32767 {
		statValue = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Str = uint16(statValue)
	character.Dex = uint16(statValue)
	character.Int = uint16(statValue)
	character.Luk = uint16(statValue)

	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_STR: int32(character.Str),
			constant.STAT_DEX: int32(character.Dex),
			constant.STAT_INT: int32(character.Int),
			constant.STAT_LUK: int32(character.Luk),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}
