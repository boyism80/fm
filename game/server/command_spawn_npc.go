package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/boyism80/fm/game/client"
)

type SpawnNpc struct {
	gs *GameServer
}

func (*SpawnNpc) New(gs *GameServer) *SpawnNpc {
	return &SpawnNpc{
		gs: gs,
	}
}

func (h *SpawnNpc) GetCommandName() string {
	return "엔피씨생성"
}

func (h *SpawnNpc) GetUsage() string {
	return "<NPCID/이름> - 현재 위치에 NPC 생성"
}

func (h *SpawnNpc) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing npcId or npc name")
	}

	var npcId uint32
	var err error

	npcId64, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		npcIdUint, ok := h.gs.resources.NameToNpc(args[0])
		if !ok {
			return fmt.Errorf("invalid npcId or npc name: %s", args[0])
		}
		npcId = npcIdUint
	} else {
		npcId = uint32(npcId64)
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}

	npc, err := mapInstance.SpawnNpc(npcId, character.Position)
	if err != nil {
		log.Printf("Failed to spawn npc %d: %v", npcId, err)
		return fmt.Errorf("failed to spawn npc: %v", err)
	}

	log.Printf("NPC spawned successfully - ID: %d, Position: %v, OID: %d",
		npcId, character.Position, npc.OID)

	return nil
}
