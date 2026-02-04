package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
)

// Damaged handles damage packet requests
type Damaged struct {
	gs     *GameServer
	opcode byte
}

func (Damaged) New(gs *GameServer) *Damaged {
	return &Damaged{
		gs:     gs,
		opcode: 0x1F,
	}
}

func (h *Damaged) GetOpcode() byte {
	return h.opcode
}

func (h *Damaged) Handle(ctx *core.ClientContext, req *request.Damaged) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	stats := map[constant.Stat]int32{}
	if !character.Invincible {
		newHp := int32(character.Life.Hp) - int32(req.Damage)
		if newHp < 0 {
			newHp = 0
		}
		if newHp > int32(character.Life.GetMaxHp()) {
			newHp = int32(character.Life.GetMaxHp())
		}

		character.Life.Hp = uint16(newHp)
		stats[constant.STAT_HP] = int32(character.Life.Hp)
	}

	character.Listener.OnUpdateStats(stats, true)

	return nil
}
