package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
)

type HealOverTime struct {
	gameServer *GameServer
	opcode     byte
}

func (HealOverTime) New(gameServer *GameServer) *HealOverTime {
	return &HealOverTime{
		gameServer: gameServer,
		opcode:     0x48,
	}
}

func (h *HealOverTime) GetOpcode() byte {
	return h.opcode
}

func (h *HealOverTime) Handle(ctx *core.ClientContext, req *request.HealOverTime) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	if character.Hp <= 0 {
		return nil
	}

	healHP := req.HealHP
	healMP := req.HealMP

	if healHP > 0 {
		newHP := character.Hp + healHP
		if newHP > character.MaxHp {
			newHP = character.MaxHp
		}
		character.Hp = newHP
	}

	if healMP > 0 {
		newMP := character.Mp + healMP
		if newMP > character.MaxMp {
			newMP = character.MaxMp
		}
		character.Mp = newMP
	}

	if character.Listener != nil {
		stats := map[constant.Stat]int32{
			constant.STAT_HP: int32(character.Hp),
			constant.STAT_MP: int32(character.Mp),
		}
		character.Listener.OnUpdateStats(stats, false)
	}

	return nil
}

