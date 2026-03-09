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
	gs *GameServer
}

func (*GainMeso) New(gs *GameServer) *GainMeso {
	return &GainMeso{
		gs: gs,
	}
}

func (h *GainMeso) GetCommandName() string {
	return "메소얻기"
}

func (h *GainMeso) GetUsage() string {
	return "<금액> - 메소 획득"
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

	character.SetMeso(character.Meso + int32(amount))
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}
