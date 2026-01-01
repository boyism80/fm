package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type GainMeso struct {
	gameServer *GameServer
}

func (*GainMeso) New(gameServer *GameServer) *GainMeso {
	return &GainMeso{
		gameServer: gameServer,
	}
}

func (h *GainMeso) GetCommandName() string {
	return "메소얻기"
}

func (h *GainMeso) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing meso amount")
	}

	amount, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid meso amount: %s", args[0])
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Meso += int32(amount)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

