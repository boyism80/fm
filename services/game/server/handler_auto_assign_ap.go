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

	if req.Amount == 0 && req.Amount2 == 0 {
		return nil
	}

	if int32(req.Amount) < 0 || int32(req.Amount2) < 0 {
		return nil
	}

	stats := map[constant.Stat]int32{}
	character.Listener.OnUpdateStats(character, stats, true)

	totalAmount := req.Amount + req.Amount2
	if character.AbilityPoint != uint16(totalAmount) {
		return nil
	}

	statUpdate := map[constant.Stat]int32{}
	success := true

	if !h.processStat(character, req.PrimaryStat, req.Amount, &statUpdate) {

		character.Listener.OnUpdateStats(character, stats, true)
		return nil
	}

	if !h.processStat(character, req.SecondaryStat, req.Amount2, &statUpdate) {

		character.Listener.OnUpdateStats(character, stats, true)
		return nil
	}

	if success {

		character.AbilityPoint = character.AbilityPoint - uint16(totalAmount)
		statUpdate[constant.STAT_AVAILABLE_AP] = int32(character.AbilityPoint)

		character.Listener.OnUpdateStats(character, statUpdate, true)
	}

	return nil
}

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
		return false
	}
}
