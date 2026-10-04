package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type Warp struct {
	gs *GameServer
}

func (Warp) New(gs *GameServer) *Warp {
	return &Warp{
		gs: gs,
	}
}

func (h *Warp) Handle(ctx *core.ClientContext, req *request.Warp) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	currentMap := character.GetMap()
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	wz := currentMap.Wz
	if wz == nil {
		return fmt.Errorf("map model not found")
	}

	if req.Target == 0xFFFFFFFF {
		portal := currentMap.FindPortalByName(req.PortalName)
		if portal == nil || portal.Wz == nil {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
		return character.EnterPortal(ctx.ActorContext, portal)
	}

	if character.GetHp() > 0 {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	character.SetHp(50, false)
	character.Stance = constant.StanceDefaultValue
	character.Listener.OnUpdateStats(character, map[constant.Stat]int32{
		constant.StatHP: int32(character.GetHp()),
	}, true)
	if sm := character.StateMachine(); sm != nil {
		sm.CallHook("on_player_revive", character)
	}

	targetMap := h.gs.GetMapSystem().Get(uint32(wz.ReturnMapId))
	if targetMap == nil {
		return fmt.Errorf("target map not found")
	}
	return character.Warp(ctx.ActorContext, targetMap, 0)
}
