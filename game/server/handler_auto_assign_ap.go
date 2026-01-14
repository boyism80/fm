package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
)

// AutoAssignAP handles AUTO_ASSIGN_AP packet requests (0x47)
type AutoAssignAP struct {
	gs     *GameServer
	opcode byte
}

func (AutoAssignAP) New(gs *GameServer) *AutoAssignAP {
	return &AutoAssignAP{
		gs:     gs,
		opcode: 0x47,
	}
}

func (h *AutoAssignAP) GetOpcode() byte {
	return h.opcode
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

	// Check if we have enough data
	if req.Amount == 0 && req.Amount2 == 0 {
		return nil
	}

	// Validate amounts
	if int32(req.Amount) < 0 || int32(req.Amount2) < 0 {
		return nil
	}

	// Send empty stat update packet first (for client synchronization)
	stats := map[constant.Stat]int32{}
	character.Listener.OnUpdateStats(stats, true)

	// Check if remaining AP matches the total amount
	totalAmount := req.Amount + req.Amount2
	if character.AbilityPoint != uint16(totalAmount) {
		return nil // AP doesn't match, ignore
	}

	statUpdate := map[constant.Stat]int32{}
	success := true

	// Process primary stat
	if !h.processStat(character, req.PrimaryStat, req.Amount, &statUpdate) {
		// Invalid primary stat - send empty stat update
		character.Listener.OnUpdateStats(stats, true)
		return nil
	}

	// Process secondary stat
	if !h.processStat(character, req.SecondaryStat, req.Amount2, &statUpdate) {
		// Invalid secondary stat - send empty stat update
		character.Listener.OnUpdateStats(stats, true)
		return nil
	}

	if success {
		// Decrease AP
		character.AbilityPoint -= uint16(totalAmount)
		statUpdate[constant.STAT_AVAILABLE_AP] = int32(character.AbilityPoint)

		// Send stat update packet
		character.Listener.OnUpdateStats(statUpdate, true)
	}

	return nil
}

// processStat processes stat distribution for a single stat
// Returns true if successful, false if invalid stat type
func (h *AutoAssignAP) processStat(character *entity.Character, statType uint32, amount uint32, statUpdate *map[constant.Stat]int32) bool {
	switch constant.StatType(statType) {
	case constant.STAT_TYPE_STR:
		if character.Str+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		character.Str += uint16(amount)
		(*statUpdate)[constant.STAT_STR] = int32(character.Str)
		return true

	case constant.STAT_TYPE_DEX:
		if character.Dex+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		character.Dex += uint16(amount)
		(*statUpdate)[constant.STAT_DEX] = int32(character.Dex)
		return true

	case constant.STAT_TYPE_INT:
		if character.Int+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		character.Int += uint16(amount)
		(*statUpdate)[constant.STAT_INT] = int32(character.Int)
		return true

	case constant.STAT_TYPE_LUK:
		if character.Luk+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		character.Luk += uint16(amount)
		(*statUpdate)[constant.STAT_LUK] = int32(character.Luk)
		return true

	default:
		return false // Invalid stat type (HP/MP not supported in AUTO_ASSIGN_AP)
	}
}
