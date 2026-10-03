package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type DamageSummon struct {
	gs *GameServer
}

func (DamageSummon) New(gs *GameServer) *DamageSummon {
	return &DamageSummon{
		gs: gs,
	}
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

	summon.TakeDamage(req.Unknown, req.Damage, req.MonsterIdFrom)
	return nil
}
