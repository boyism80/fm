package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type NpcClick struct {
	gs *GameServer
}

func (NpcClick) New(gs *GameServer) *NpcClick {
	return &NpcClick{
		gs: gs,
	}
}

func (h *NpcClick) Handle(ctx *core.ClientContext, req *request.NpcClick) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	if character.Spectating() || character.UsingCashItem() || character.Occupied() {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Map not found")
		return fmt.Errorf("map not found")
	}

	npc, ok := mapInstance.GetNpcs()[req.OID].(*game.Npc)
	if !ok {
		log.Printf("NPC %d not found on map", req.OID)
		return fmt.Errorf("npc %d not found", req.OID)
	}

	if err := character.Dialog().Open(ctx.ActorContext, npc); err != nil {
		if errors.Is(err, game.ErrDialogBusy) {
			character.Listener().OnUpdateStats(character, nil, true)
			return nil
		}
		log.Printf("Failed to open NPC %d: %v", npc.Wz().ID, err)
		return err
	}
	return nil
}
