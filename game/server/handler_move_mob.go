package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
)

type MoveMob struct {
	gs     *GameServer
	opcode byte
}

func (MoveMob) New(gs *GameServer) *MoveMob {
	return &MoveMob{
		gs:     gs,
		opcode: 0x95,
	}
}

func (h *MoveMob) GetOpcode() byte {
	return h.opcode
}

func (h *MoveMob) Handle(ctx *core.ClientContext, req *request.MoveMob) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		log.Printf("Character not on a map")
		return fmt.Errorf("map not found")
	}

	mob := mapInstance.GetMob(req.OID)
	if mob == nil {
		return nil
	}

	controllerTable := mapInstance.GetControllerTable()
	controller, exists := controllerTable.GetController(mob)
	if !exists {
		log.Printf("No controller found for mob %d", req.OID)
		return fmt.Errorf("no controller found for this mob")
	}

	if controller.GetID() != ch.GetID() {
		return nil
	}

	startPoint := mob.Position

	for _, mnt := range req.Movements {
		if move, ok := mnt.(*dto.AbsoluteLifeMovement); ok {
			mob.Position = move.Position
		}

		mob.Stance = mnt.GetStance()
	}

	ch.Listener.OnControlMoveMob(ch, mob, req.MovementId, req.IsAggroed, uint16(min(mob.Mp, 65535)), 0, 0)

	ch.Listener.OnMobMoved(
		ch,
		mob,
		req.IsAggroed,
		req.CenterSplit,
		req.Skill1,
		req.Skill2,
		req.Skill3,
		req.Skill4,
		startPoint,
		req.Movements,
	)

	return nil
}
