package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/boyism80/fm/game/client"
)

type SpawnMob struct {
	gs *GameServer
}

func (*SpawnMob) New(gs *GameServer) *SpawnMob {
	return &SpawnMob{
		gs: gs,
	}
}

func (h *SpawnMob) GetCommandName() string {
	return "몬스터생성"
}

func (h *SpawnMob) GetUsage() string {
	return "<몬스터ID/이름> [개수] - 몬스터 생성 (개수 기본값 1)"
}

func (h *SpawnMob) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing mobId or mob name")
	}

	var mobId uint32
	var err error

	mobId64, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		mobIdUint, ok := h.gs.resources.NameToMob(args[0])
		if !ok {
			return fmt.Errorf("invalid mobId or mob name: %s", args[0])
		}
		mobId = mobIdUint
	} else {
		mobId = uint32(mobId64)
	}

	count := 1
	if len(args) >= 2 {
		count64, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil || count64 < 1 {
			return fmt.Errorf("count must be a positive integer")
		}
		count = int(count64)
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}

	_, ok := h.gs.resources.Monsters[mobId]
	if !ok {
		return fmt.Errorf("mob specification not found for ID: %d", mobId)
	}

	for i := 0; i < count; i++ {
		mob, err := mapInstance.SpawnMob(mobId, character.Position, nil)
		if err != nil {
			log.Printf("Failed to spawn mob %d (attempt %d/%d): %v", mobId, i+1, count, err)
			return fmt.Errorf("failed to spawn mob: %v", err)
		}
		log.Printf("Mob spawned - ID: %d, Position: %v, OID: %d (%d/%d)",
			mobId, character.Position, mob.OID, i+1, count)
	}

	return nil
}
