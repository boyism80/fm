package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeLuk struct {
	gs *GameServer
}

func (*ChangeLuk) New(gs *GameServer) *ChangeLuk {
	return &ChangeLuk{
		gs: gs,
	}
}

func (h *ChangeLuk) GetCommandName() string {
	return "행운바꾸기"
}

func (h *ChangeLuk) GetUsage() string {
	return "<행운값> - 행운 설정"
}

func (h *ChangeLuk) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing LUK value")
	}

	luk, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid LUK value: %s", args[0])
	}

	if luk > 32767 {
		luk = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.BaseStats.Luk = uint16(luk)
	character.BonusStats.Luk = 0
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_LUK: int32(character.GetTotalLuk()),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}
