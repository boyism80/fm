package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

type DamageSummon struct {
	gs     *GameServer
	opcode byte
}

func (DamageSummon) New(gs *GameServer) *DamageSummon {
	return &DamageSummon{
		gs:     gs,
		opcode: 0x8E,
	}
}

func (h *DamageSummon) GetOpcode() byte {
	return h.opcode
}

func (h *DamageSummon) Handle(ctx *core.ClientContext, req *request.DamageSummon) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return nil
	}

	summon := mapInstance.GetSummon(req.OID)
	if summon == nil {
		return nil
	}
	if summon.OwnerID != character.GetID() {
		return nil
	}

	damage := req.Damage
	if damage > 0 {
		if summon.Hp <= damage {
			summon.Hp = 0
		} else {
			summon.Hp -= damage
		}
		character.Listener.OnSummonDamaged(character, summon, req.Unknown, damage, req.MonsterIdFrom)
		if summon.Hp == 0 {
			character.Buffs.RemoveSkillBuff(uint32(summon.SkillID))
		}
	}

	return nil
}
