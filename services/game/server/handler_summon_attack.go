package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type SummonAttack struct {
	gs *GameServer
}

func (SummonAttack) New(gs *GameServer) *SummonAttack {
	return &SummonAttack{
		gs: gs,
	}
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

	character.DamageTo(req.Damages)
	CallSummonOnAttackHooks(character, mapInstance, req.Damages, uint32(summon.SkillID))

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
