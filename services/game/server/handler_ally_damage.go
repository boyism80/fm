package server

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
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
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}

	from := mapInstance.GetMob(req.FromOID)
	to := mapInstance.GetMob(req.ToOID)
	if from == nil || to == nil || to.Wz == nil || !to.Wz.DamagedByMob {
		return nil
	}

	level := int(to.Wz.Level)
	if level <= 0 {
		return nil
	}
	damage := (level * rand.Intn(level)) / 3
	if damage > 0 {
		killed := to.ApplyDamage(ch, uint32(damage))
		if !killed && to.Listener != nil {
			to.Listener.OnMobAllyDamaged(to, int32(damage))
		}
	} else {
		to.MarkHit()
	}
	return nil
}
