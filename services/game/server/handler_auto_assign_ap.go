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

	stats := map[constant.Stat]int32{}
	character.Listener.OnUpdateStats(character, stats, true)

	if character.AbilityPoint != uint16(totalAmount) {
		return nil
	}

	statUpdate := map[constant.Stat]int32{}
	for _, entry := range req.Entries {
		if !h.addStat(character, entry.Stat, entry.Amount, &statUpdate) {
			character.Listener.OnUpdateStats(character, stats, true)
			return nil
		}
	}

	character.AbilityPoint = character.AbilityPoint - uint16(totalAmount)
	statUpdate[constant.StatAvailableAP] = int32(character.AbilityPoint)

	character.Listener.OnUpdateStats(character, statUpdate, true)

	return nil
}

func (h *AutoAssignAP) addStat(character *entity.Character, statType uint32, amount uint32, statUpdate *map[constant.Stat]int32) bool {
	switch constant.StatType(statType) {
	case constant.StatTypeStr:
		if character.GetTotalStr()+uint16(amount) > constant.StatMaxStrDexIntLuk {
			return false
		}
		newStr := character.BaseStats.Str + uint16(amount)
		if newStr > constant.StatMaxStrDexIntLuk {
			newStr = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Str = newStr
		(*statUpdate)[constant.StatStr] = int32(character.GetTotalStr())
		return true

	case constant.StatTypeDex:
		if character.GetTotalDex()+uint16(amount) > constant.StatMaxStrDexIntLuk {
			return false
		}
		newDex := character.BaseStats.Dex + uint16(amount)
		if newDex > constant.StatMaxStrDexIntLuk {
			newDex = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Dex = newDex
		(*statUpdate)[constant.StatDex] = int32(character.GetTotalDex())
		return true

	case constant.StatTypeInt:
		if character.GetTotalInt()+uint16(amount) > constant.StatMaxStrDexIntLuk {
			return false
		}
		newInt := character.BaseStats.Int + uint16(amount)
		if newInt > constant.StatMaxStrDexIntLuk {
			newInt = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Int = newInt
		(*statUpdate)[constant.StatInt] = int32(character.GetTotalInt())
		return true

	case constant.StatTypeLuk:
		if character.GetTotalLuk()+uint16(amount) > constant.StatMaxStrDexIntLuk {
			return false
		}
		newLuk := character.BaseStats.Luk + uint16(amount)
		if newLuk > constant.StatMaxStrDexIntLuk {
			newLuk = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Luk = newLuk
		(*statUpdate)[constant.StatLuk] = int32(character.GetTotalLuk())
		return true

	default:
		return false
	}
}
