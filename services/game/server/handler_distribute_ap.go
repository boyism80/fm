package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type DistributeAP struct {
	gs *GameServer
}

func (DistributeAP) New(gs *GameServer) *DistributeAP {
	return &DistributeAP{
		gs: gs,
	}
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

	stats := map[constant.Stat]int32{}
	character.Listener.OnUpdateStats(character, stats, true)

	if character.AbilityPoint == 0 {
		return nil
	}

	statUpdate := map[constant.Stat]int32{}
	success := false

	switch constant.StatType(req.StatType) {
	case constant.StatTypeStr:
		if character.GetTotalStr() >= constant.StatMaxStrDexIntLuk {
			return nil
		}
		newStr := character.BaseStats.Str + 1
		if newStr > constant.StatMaxStrDexIntLuk {
			newStr = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Str = newStr
		statUpdate[constant.StatStr] = int32(character.GetTotalStr())
		success = true

	case constant.StatTypeDex:
		if character.GetTotalDex() >= constant.StatMaxStrDexIntLuk {
			return nil
		}
		newDex := character.BaseStats.Dex + 1
		if newDex > constant.StatMaxStrDexIntLuk {
			newDex = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Dex = newDex
		statUpdate[constant.StatDex] = int32(character.GetTotalDex())
		success = true

	case constant.StatTypeInt:
		if character.GetTotalInt() >= constant.StatMaxStrDexIntLuk {
			return nil
		}
		newInt := character.BaseStats.Int + 1
		if newInt > constant.StatMaxStrDexIntLuk {
			newInt = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Int = newInt
		statUpdate[constant.StatInt] = int32(character.GetTotalInt())
		success = true

	case constant.StatTypeLuk:
		if character.GetTotalLuk() >= constant.StatMaxStrDexIntLuk {
			return nil
		}
		newLuk := character.BaseStats.Luk + 1
		if newLuk > constant.StatMaxStrDexIntLuk {
			newLuk = constant.StatMaxStrDexIntLuk
		}
		character.BaseStats.Luk = newLuk
		statUpdate[constant.StatLuk] = int32(character.GetTotalLuk())
		success = true

	case constant.StatTypeHP:
		if character.HpApUsed >= constant.HpAPUsedMax || character.GetMaxHp() >= constant.StatMaxHPMP {
			return nil
		}
		h.callAPToStatScript(ctx, character, "get_ap_to_hp", func(hpIncrease uint32) {
			if hpIncrease == 0 {
				hpIncrease = 10
			}
			character.AddBaseHp(hpIncrease, false)
			character.HpApUsed++
			statUpdate := map[constant.Stat]int32{
				constant.StatHP:    int32(character.GetHp()),
				constant.StatMaxHP: int32(character.GetMaxHp()),
			}
			h.commitAP(character, statUpdate)
		})
		return nil

	case constant.StatTypeMP:
		if character.HpApUsed >= constant.HpAPUsedMax || character.GetMaxMp() >= constant.StatMaxHPMP {
			return nil
		}
		h.callAPToStatScript(ctx, character, "get_ap_to_mp", func(mpIncrease uint32) {
			if mpIncrease == 0 {
				mpIncrease = 5
			}
			character.AddBaseMp(mpIncrease, false)
			character.HpApUsed++
			statUpdate := map[constant.Stat]int32{
				constant.StatMP:    int32(character.GetMp()),
				constant.StatMaxMP: int32(character.GetMaxMp()),
			}
			h.commitAP(character, statUpdate)
		})
		return nil

	default:

		character.Listener.OnUpdateStats(character, stats, true)
		return nil
	}

	if success {

		character.AbilityPoint = character.AbilityPoint - 1
		statUpdate[constant.StatAvailableAP] = int32(character.AbilityPoint)

		character.Listener.OnUpdateStats(character, statUpdate, true)
	}

	return nil
}

func (h *DistributeAP) commitAP(character *entity.Character, statUpdate map[constant.Stat]int32) {
	character.AbilityPoint = character.AbilityPoint - 1
	statUpdate[constant.StatAvailableAP] = int32(character.AbilityPoint)
	character.Listener.OnUpdateStats(character, statUpdate, true)
}

func (h *DistributeAP) callAPToStatScript(ctx *core.ClientContext, character *entity.Character, funcName string, fn func(uint32)) {
	if character == nil || character.GetMap() == nil {
		fn(0)
		return
	}
	root := character.GetMap().GetLuaRoot()
	if root == nil {
		fn(0)
		return
	}
	thread, err := luax.NewThread(root, constant.CharacterQueryScriptPath)
	if err != nil {
		fn(0)
		return
	}
	luax.CallAsync(root, thread, funcName, character).Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) == 0 || vals[0] == nil || vals[0].Type() != lua.LTNumber {
			fn(0)
			return nil, nil
		}
		n := float64(lua.LVAsNumber(vals[0]))
		if n <= 0 {
			fn(0)
			return nil, nil
		}
		if n >= float64(0xffffffff) {
			fn(0xffffffff)
			return nil, nil
		}
		fn(uint32(n))
		return nil, nil
	}).OnError(func(err error) {
		fn(0)
	})
}
