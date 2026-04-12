package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
)

type SummonAttack struct {
	gs     *GameServer
	opcode byte
}

func (SummonAttack) New(gs *GameServer) *SummonAttack {
	return &SummonAttack{
		gs:     gs,
		opcode: 0x8D,
	}
}

func (h *SummonAttack) GetOpcode() byte {
	return h.opcode
}

func (h *SummonAttack) Handle(ctx *core.ClientContext, req *request.SummonAttack) error {
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

	summon := mapInstance.GetSummon(req.SummonOID)
	if summon == nil {
		return nil
	}
	if summon.OwnerID != character.GetID() {
		return nil
	}

	ApplyDamageToMobs(character, mapInstance, req.Damages)
	CallSummonOnAttackHooks(ctx, character, mapInstance, req.Damages, uint32(summon.SkillID))

	var targets []entity.SummonAttackTarget
	for _, ap := range req.Damages {
		for _, dp := range ap.DamagePairs {
			if dp.Damage == 0 {
				continue
			}
			targets = append(targets, entity.SummonAttackTarget{
				OID:    ap.OID,
				Damage: dp.Damage,
			})
		}
	}
	if len(targets) > 0 {
		character.Listener.OnSummonAttack(character, summon, req.Animation, targets)
	}
	return nil
}
