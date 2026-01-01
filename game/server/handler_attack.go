package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

// Attack handles attack packet requests
type Attack struct {
	gameServer *GameServer
	opcode     byte
}

func (Attack) New(gameServer *GameServer) *Attack {
	return &Attack{
		gameServer: gameServer,
		opcode:     0x1B,
	}
}

func (h *Attack) GetOpcode() byte {
	return h.opcode
}

func (h *Attack) Handle(ctx *core.ClientContext, req *request.Attack) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	mapID := character.GetMap()
	mapInstance := h.gameServer.GetMap(mapID)
	if mapInstance == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	for _, damage := range req.AttackInfo.Damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil {
			log.Printf("Mob not found for OID: %d", damage.OID)
			continue
		}

		for _, damagePair := range damage.DamagePairs {
			mob.Damage(uint16(damagePair.Damage), character)
		}
	}

	character.Listener.OnAttack(mapID, character.GetID(), req.AttackInfo, 0)

	return nil
}
