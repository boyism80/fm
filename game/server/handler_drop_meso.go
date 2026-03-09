package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
)

// DropMeso handles meso dropping packet requests
type DropMeso struct {
	gs     *GameServer
	opcode byte
}

func (DropMeso) New(gs *GameServer) *DropMeso {
	return &DropMeso{
		gs:     gs,
		opcode: 0x4D,
	}
}

func (h *DropMeso) GetOpcode() byte {
	return h.opcode
}

func (h *DropMeso) Handle(ctx *core.ClientContext, req *request.DropMeso) error {
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

	if req.Count < 10 || req.Count > 50000 {
		log.Printf("Invalid meso count: %d", req.Count)
		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	if req.Count > character.Meso {
		log.Printf("Character doesn't have enough meso: %d < %d", character.Meso, req.Count)
		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	character.SetMeso(character.Meso - req.Count)
	character.Listener.OnUpdateStats(map[constant.Stat]int32{
		constant.STAT_MESO: character.Meso,
	}, true)

	mapInstance := h.gs.GetMap(character.Map)
	if mapInstance != nil {
		if err := mapInstance.SpawnMeso(req.Count, character.Position, character.GetID(), constant.DROP_TYPE_FFA); err != nil {
			log.Printf("Failed to spawn meso on map: %v", err)
		}
	}

	return nil
}
