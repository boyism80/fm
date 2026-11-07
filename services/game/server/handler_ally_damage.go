package server

import (
	"fmt"
	"math/rand"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type AllyDamage struct {
	gs *GameServer
}

func (AllyDamage) New(gs *GameServer) *AllyDamage {
	return &AllyDamage{
		gs: gs,
	}
}

func (h *AllyDamage) Handle(ctx *core.ClientContext, req *request.AllyDamage) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}

	from := mapInstance.GetMob(req.FromOID)
	to := mapInstance.GetMob(req.ToOID)
	if from == nil || to == nil || to.Wz() == nil || !to.Wz().DamagedByMob {
		return nil
	}

	level := int(to.Wz().Level)
	if level <= 0 {
		return nil
	}
	damage := (level * rand.Intn(level)) / 3
	if damage > 0 {
		if to.TakeDamage(ch, uint32(damage)) {
			ch.OnKill([]*game.Mob{to})
		} else if to.Listener() != nil {
			to.Listener().OnMobAllyDamaged(to, int32(damage))
		}
	} else {
		to.MarkHit()
	}
	return nil
}
