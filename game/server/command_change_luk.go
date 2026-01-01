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
	gameServer *GameServer
}

func (*ChangeLuk) New(gameServer *GameServer) *ChangeLuk {
	return &ChangeLuk{
		gameServer: gameServer,
	}
}

func (h *ChangeLuk) GetCommandName() string {
	return "행운바꾸기"
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

	character.Luk = uint16(luk)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_LUK: int32(character.Luk),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

