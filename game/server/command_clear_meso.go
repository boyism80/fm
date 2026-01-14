package server

import (
	"fmt"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ClearMeso struct {
	gs *GameServer
}

func (*ClearMeso) New(gs *GameServer) *ClearMeso {
	return &ClearMeso{
		gs: gs,
	}
}

func (h *ClearMeso) GetCommandName() string {
	return "메소초기화"
}

func (h *ClearMeso) GetUsage() string {
	return "- 메소 초기화"
}

func (h *ClearMeso) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Meso = 0
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}
