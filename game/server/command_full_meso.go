package server

import (
	"fmt"
	"math"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type FullMeso struct {
	gameServer *GameServer
}

func (*FullMeso) New(gameServer *GameServer) *FullMeso {
	return &FullMeso{
		gameServer: gameServer,
	}
}

func (h *FullMeso) GetCommandName() string {
	return "풀메소"
}

func (h *FullMeso) GetUsage() string {
	return "- 메소 최대치로 설정"
}

func (h *FullMeso) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Meso = math.MaxInt32
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

