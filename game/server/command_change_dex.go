package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeDex struct {
	gameServer *GameServer
}

func (*ChangeDex) New(gameServer *GameServer) *ChangeDex {
	return &ChangeDex{
		gameServer: gameServer,
	}
}

func (h *ChangeDex) GetCommandName() string {
	return "민첩바꾸기"
}

func (h *ChangeDex) GetUsage() string {
	return "<민첩값> - 민첩 설정"
}

func (h *ChangeDex) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing DEX value")
	}

	dex, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid DEX value: %s", args[0])
	}

	if dex > 32767 {
		dex = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Dex = uint16(dex)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_DEX: int32(character.Dex),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

