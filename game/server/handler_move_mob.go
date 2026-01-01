package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
)

// MoveMob handles mob movement packet requests
type MoveMob struct {
	gameServer *GameServer
	opcode     byte
}

func (MoveMob) New(gameServer *GameServer) *MoveMob {
	return &MoveMob{
		gameServer: gameServer,
		opcode:     0x95,
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

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	mapInstance := h.gameServer.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	mob := mapInstance.GetMob(req.OID)
	if mob == nil {
		log.Printf("Mob %d not found on map %d", req.OID, character.GetMap())
		return fmt.Errorf("mob %d not found", req.OID)
	}

	controllerTable := mapInstance.GetControllerTable()
	controller, exists := controllerTable.GetController(mob)
	if !exists {
		log.Printf("No controller found for mob %d", req.OID)
		return fmt.Errorf("no controller found for this mob")
	}

	if controller.GetID() != character.GetID() {
		return nil
	}

	startPoint := mob.Position

	for _, mnt := range req.Movements {
		if move, ok := mnt.(*dto.AbsoluteLifeMovement); ok {
			mob.Position = move.Position
		}

		mob.Stance = mnt.GetStance()
	}

	character.Listener.OnControlMoveMob(req.OID, uint8(req.MovementId), req.IsAggroed, mob.Mp, 0, 0)

	character.Listener.OnMobMoved(
		character.GetMap(),
		req.OID,
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
