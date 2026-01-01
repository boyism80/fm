package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeMp struct {
	gameServer *GameServer
}

func (*ChangeMp) New(gameServer *GameServer) *ChangeMp {
	return &ChangeMp{
		gameServer: gameServer,
	}
}

func (h *ChangeMp) GetCommandName() string {
	return "마력바꾸기"
}

func (h *ChangeMp) GetUsage() string {
	return "<마력값> - 마력 설정"
}

func (h *ChangeMp) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing MP value")
	}

	mp, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid MP value: %s", args[0])
	}

	if mp > 32767 {
		mp = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Mp = uint16(mp)
	character.MaxMp = uint16(mp)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MP:     int32(character.Mp),
			constant.STAT_MAX_MP: int32(character.MaxMp),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

