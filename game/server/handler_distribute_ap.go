package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
	lua "github.com/yuin/gopher-lua"
)

// DistributeAP handles DISTRIBUTE_AP packet requests (0x46)
type DistributeAP struct {
	gs     *GameServer
	opcode byte
}

func (DistributeAP) New(gs *GameServer) *DistributeAP {
	return &DistributeAP{
		gs:     gs,
		opcode: 0x46,
	}
}

func (h *DistributeAP) GetOpcode() byte {
	return h.opcode
}

func (h *DistributeAP) Handle(ctx *core.ClientContext, req *request.DistributeAP) error {
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

	// Send empty stat update packet first (for client synchronization)
	stats := map[constant.Stat]int32{}
	character.Listener.OnUpdateStats(character, stats, true)

	// Check if character has remaining AP
	if character.AbilityPoint == 0 {
		return nil
	}

	// Process stat distribution
	statUpdate := map[constant.Stat]int32{}
	success := false

	switch constant.StatType(req.StatType) {
	case constant.STAT_TYPE_STR:
		if character.GetTotalStr() >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		newStr := character.BaseStats.Str + 1
		if newStr > constant.STAT_MAX_STR_DEX_INT_LUK {
			newStr = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Str = newStr
		statUpdate[constant.STAT_STR] = int32(character.GetTotalStr())
		success = true

	case constant.STAT_TYPE_DEX:
		if character.GetTotalDex() >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		newDex := character.BaseStats.Dex + 1
		if newDex > constant.STAT_MAX_STR_DEX_INT_LUK {
			newDex = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Dex = newDex
		statUpdate[constant.STAT_DEX] = int32(character.GetTotalDex())
		success = true

	case constant.STAT_TYPE_INT:
		if character.GetTotalInt() >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		newInt := character.BaseStats.Int + 1
		if newInt > constant.STAT_MAX_STR_DEX_INT_LUK {
			newInt = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Int = newInt
		statUpdate[constant.STAT_INT] = int32(character.GetTotalInt())
		success = true

	case constant.STAT_TYPE_LUK:
		if character.GetTotalLuk() >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		newLuk := character.BaseStats.Luk + 1
		if newLuk > constant.STAT_MAX_STR_DEX_INT_LUK {
			newLuk = constant.STAT_MAX_STR_DEX_INT_LUK
		}
		character.BaseStats.Luk = newLuk
		statUpdate[constant.STAT_LUK] = int32(character.GetTotalLuk())
		success = true

	case constant.STAT_TYPE_HP:
		if character.HpApUsed >= constant.HP_AP_USED_MAX || character.GetMaxHp() >= constant.STAT_MAX_HP_MP {
			return nil
		}
		hpIncrease := h.scriptAPToHP(ctx, character)
		if hpIncrease == 0 {
			hpIncrease = 10
		}
		character.AddBaseHp(hpIncrease)
		character.HpApUsed++
		statUpdate[constant.STAT_MAX_HP] = int32(character.GetMaxHp())
		success = true

	case constant.STAT_TYPE_MP:
		if character.HpApUsed >= constant.HP_AP_USED_MAX || character.GetMaxMp() >= constant.STAT_MAX_HP_MP {
			return nil
		}
		mpIncrease := h.scriptAPToMP(ctx, character)
		if mpIncrease == 0 {
			mpIncrease = 5
		}
		character.AddBaseMp(mpIncrease)
		character.HpApUsed++
		statUpdate[constant.STAT_MAX_MP] = int32(character.GetMaxMp())
		success = true

	default:
		// Invalid stat type - send empty stat update
		character.Listener.OnUpdateStats(character, stats, true)
		return nil
	}

	if success {
		// Decrease AP
		character.AbilityPoint = character.AbilityPoint - 1
		statUpdate[constant.STAT_AVAILABLE_AP] = int32(character.AbilityPoint)

		// Send stat update packet
		character.Listener.OnUpdateStats(character, statUpdate, true)
	}

	return nil
}

func (h *DistributeAP) scriptAPToHP(ctx *core.ClientContext, character interface{}) uint32 {
	return h.callAPToStatScript(ctx, character, "on_ap_to_hp")
}

func (h *DistributeAP) scriptAPToMP(ctx *core.ClientContext, character interface{}) uint32 {
	return h.callAPToStatScript(ctx, character, "on_ap_to_mp")
}

func (h *DistributeAP) callAPToStatScript(ctx *core.ClientContext, character interface{}, funcName string) uint32 {
	if ctx.LogicActorPID == nil {
		return 0
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		return 0
	}
	result, thread, err := luax.Call(root, "script/script.lua", funcName, character)
	if thread != nil {
		thread.Close()
	}
	if err != nil || result == nil || result.Type() != lua.LTNumber {
		return 0
	}
	n := float64(lua.LVAsNumber(result))
	if n <= 0 {
		return 0
	}
	if n >= float64(0xffffffff) {
		return 0xffffffff
	}
	return uint32(n)
}
