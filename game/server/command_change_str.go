package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeStr struct {
	gs *GameServer
}

func (*ChangeStr) New(gs *GameServer) *ChangeStr {
	return &ChangeStr{
		gs: gs,
	}
}

func (h *ChangeStr) GetCommandName() string {
	return "힘바꾸기"
}

func (h *ChangeStr) GetUsage() string {
	return "<힘값> - 힘 설정"
}

func (h *ChangeStr) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing STR value")
	}

	str, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid STR value: %s", args[0])
	}

	if str > 32767 {
		str = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.BaseStats.Str = uint16(str)
	character.BonusStats.Str = 0
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_STR: int32(character.GetTotalStr()),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}
