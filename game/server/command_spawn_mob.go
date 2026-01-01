package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/boyism80/fm/game/client"
)

type SpawnMob struct {
	gameServer *GameServer
}

func (*SpawnMob) New(gameServer *GameServer) *SpawnMob {
	return &SpawnMob{
		gameServer: gameServer,
	}
}

func (h *SpawnMob) GetCommandName() string {
	return "몬스터생성"
}

func (h *SpawnMob) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing mobId or mob name")
	}

	var mobId uint32
	var err error

	mobId64, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		mobIdUint, ok := h.gameServer.resources.NameToMob(args[0])
		if !ok {
			return fmt.Errorf("invalid mobId or mob name: %s", args[0])
		}
		mobId = mobIdUint
	} else {
		mobId = uint32(mobId64)
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	mapInstance := h.gameServer.GetMap(character.Map)
	if mapInstance == nil {
		return fmt.Errorf("map %d not found", character.Map)
	}

	_, ok := h.gameServer.resources.Monsters[mobId]
	if !ok {
		return fmt.Errorf("mob specification not found for ID: %d", mobId)
	}

	mob, err := mapInstance.SpawnMob(mobId, character.Position)
	if err != nil {
		log.Printf("Failed to spawn mob %d: %v", mobId, err)
		return fmt.Errorf("failed to spawn mob: %v", err)
	}

	log.Printf("Mob spawned successfully - ID: %d, Position: %v, OID: %d",
		mobId, character.Position, mob.OID)

	return nil
}

