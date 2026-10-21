package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
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

	summon, err := character.Summons().Attack(ctx.ActorContext, req.SummonOID, req.Damages)
	if err != nil {
		return err
	}
	if summon == nil {
		return nil
	}

	var targets []game.SummonAttackTarget
	for _, ap := range req.Damages {
		for _, dp := range ap.DamagePairs {
			if dp.Damage == 0 {
				continue
			}
			targets = append(targets, game.SummonAttackTarget{
				OID:    ap.OID,
				Damage: dp.Damage,
			})
		}
	}
	if len(targets) > 0 {
		character.Listener().OnSummonAttack(character, summon, req.Animation, targets)
	}
	summon.CallAttackHook(req.Damages)
	return nil
}
