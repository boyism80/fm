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
		if character.GetTotalStr()+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		newStr := character.BaseStats.Str + uint16(amount)
		if newStr > constant.STAT_MAX_STR_DEX_INT_LUK {
			newStr = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Str = newStr
		(*statUpdate)[constant.STAT_STR] = int32(character.GetTotalStr())
		return true

	case constant.STAT_TYPE_DEX:
		if character.GetTotalDex()+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		newDex := character.BaseStats.Dex + uint16(amount)
		if newDex > constant.STAT_MAX_STR_DEX_INT_LUK {
			newDex = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Dex = newDex
		(*statUpdate)[constant.STAT_DEX] = int32(character.GetTotalDex())
		return true

	case constant.STAT_TYPE_INT:
		if character.GetTotalInt()+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		newInt := character.BaseStats.Int + uint16(amount)
		if newInt > constant.STAT_MAX_STR_DEX_INT_LUK {
			newInt = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Int = newInt
		(*statUpdate)[constant.STAT_INT] = int32(character.GetTotalInt())
		return true

	case constant.STAT_TYPE_LUK:
		if character.GetTotalLuk()+uint16(amount) > constant.STAT_MAX_STR_DEX_INT_LUK {
			return false
		}
		newLuk := character.BaseStats.Luk + uint16(amount)
		if newLuk > constant.STAT_MAX_STR_DEX_INT_LUK {
			newLuk = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Luk = newLuk
		(*statUpdate)[constant.STAT_LUK] = int32(character.GetTotalLuk())
		return true

	default:
		return false // Invalid stat type (HP/MP not supported in AUTO_ASSIGN_AP)
	}
}
