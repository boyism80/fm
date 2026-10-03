package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type AutoAssignAP struct {
	gs *GameServer
}

func (AutoAssignAP) New(gs *GameServer) *AutoAssignAP {
	return &AutoAssignAP{
		gs: gs,
	}
}

func (h *AutoAssignAP) Handle(ctx *core.ClientContext, req *request.AutoAssignAP) error {
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

	totalAmount := uint32(0)
	for _, entry := range req.Entries {
		if int32(entry.Amount) < 0 {
			return nil
		}
		totalAmount += entry.Amount
	}
	if totalAmount == 0 {
		return nil
	}

	emptyStats := map[constant.Stat]int32{}
	character.Listener.OnUpdateStats(character, emptyStats, true)

	if character.AbilityPoint != uint16(totalAmount) {
		return nil
	}

	entries := make([]entity.APEntry, len(req.Entries))
	for i, e := range req.Entries {
		entries[i] = entity.APEntry{Stat: constant.StatType(e.Stat), Amount: uint16(e.Amount)}
	}
	if !character.AssignAP(entries) {
		character.Listener.OnUpdateStats(character, emptyStats, true)
	}
	return nil
}
