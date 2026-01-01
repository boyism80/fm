package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type ChangeHp struct {
	gameServer *GameServer
}

func (*ChangeHp) New(gameServer *GameServer) *ChangeHp {
	return &ChangeHp{
		gameServer: gameServer,
	}
}

func (h *ChangeHp) GetCommandName() string {
	return "체력바꾸기"
}

func (h *ChangeHp) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing HP value")
	}

	hp, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid HP value: %s", args[0])
	}

	if hp > 32767 {
		hp = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Hp = uint16(hp)
	character.MaxHp = uint16(hp)
	gameClient.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_HP:     int32(character.Hp),
			constant.STAT_MAX_HP: int32(character.MaxHp),
		},
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}

