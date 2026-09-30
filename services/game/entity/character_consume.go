package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) UseItemEffect(itemID uint32) bool {
	model := ch.itemModel(itemID)
	consume, ok := model.(*wz.Consume)
	if !ok || consume == nil {
		return false
	}
	ch.addItemBuff(consume)
	ch.recover(consume)
	return true
}

func (ch *Character) UseConsume(consume *Consume) bool {
	if ch == nil || consume == nil {
		return false
	}
	wzConsume, ok := consume.GetModel().(*wz.Consume)
	if !ok || wzConsume == nil {
		return false
	}

	applyWZ, scriptOK := ch.callActiveConsumeScript(consume)
	if !scriptOK {
		return false
	}

	if applyWZ {
		var targets []*Character
		if wzConsume.Party {
			if pid := ch.GetPartyID(); pid != nil {
				if m := ch.GetMap(); m != nil {
					if members := m.GetPartyMembers(*pid); len(members) > 0 {
						targets = members
					}
				}
			}
		}
		if len(targets) == 0 {
			targets = []*Character{ch}
		}

		for _, character := range targets {
			if character == nil {
				continue
			}
			if len(wzConsume.CureDebuffs) > 0 {
				character.RemoveDebuff(wzConsume.CureDebuffs...)
			}
			character.addItemBuff(wzConsume)
			character.recover(wzConsume)
			if wzConsume.ExpInc > 0 {
				character.AddExp(uint32(wzConsume.ExpInc))
			}
		}
	}

	return true
}

func (ch *Character) addItemBuff(consumeItem *wz.Consume) bool {
	if ch == nil || consumeItem == nil {
		return false
	}
	buffValues := consumeItem.BuffSpecValues()
	if len(buffValues) == 0 {
		return false
	}
	ch.Buffs.AddItemBuff(consumeItem, consumeItem.BuffDuration, buffValues, true, true)
	return true
}

func (ch *Character) recover(consumeItem *wz.Consume) bool {
	if ch == nil || consumeItem == nil {
		return false
	}
	healMul := ch.PotionHealMultiplierPercent()
	hpChange := int(consumeItem.ActiveEffect.HP) * healMul / 100
	mpChange := int(consumeItem.ActiveEffect.MP) * healMul / 100
	hpRate := int(consumeItem.ActiveEffect.HPRate)
	mpRate := int(consumeItem.ActiveEffect.MPRate)
	if hpRate > 0 {
		hpChange += int(ch.GetMaxHp()) * hpRate / 100
	}
	if mpRate > 0 {
		mpChange += int(ch.GetMaxMp()) * mpRate / 100
	}
	if hpChange == 0 && mpChange == 0 {
		return false
	}

	stats := make(map[constant.Stat]int32)
	if hpChange != 0 {
		newHP := int(ch.GetHp()) + hpChange
		if newHP < 1 {
			newHP = 1
		}
		maxHp := int(ch.GetMaxHp())
		if newHP > maxHp {
			newHP = maxHp
		}
		if newHP != int(ch.GetHp()) {
			ch.SetHp(uint32(newHP), false)
			stats[constant.StatHP] = int32(ch.GetHp())
		}
	}
	if mpChange != 0 {
		newMP := int(ch.GetMp()) + mpChange
		if newMP < 0 {
			newMP = 0
		}
		maxMp := int(ch.GetMaxMp())
		if newMP > maxMp {
			newMP = maxMp
		}
		if newMP != int(ch.GetMp()) {
			ch.SetMp(uint32(newMP), false)
			stats[constant.StatMP] = int32(ch.GetMp())
		}
	}
	if len(stats) == 0 {
		return false
	}
	ch.Listener.OnUpdateStats(ch, stats, false)
	return true
}

func (ch *Character) callActiveConsumeScript(consume *Consume) (applyWZ bool, scriptOK bool) {
	applyWZ, scriptOK = true, true
	if ch == nil || consume == nil {
		return applyWZ, scriptOK
	}
	consumeWz, ok := consume.GetModel().(*wz.Consume)
	if !ok || consumeWz == nil || consumeWz.ID == 0 {
		return applyWZ, scriptOK
	}
	itemID := consumeWz.ID

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return applyWZ, scriptOK
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return applyWZ, scriptOK
	}

	scriptPath := fmt.Sprintf("script/item/%d.lua", itemID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return applyWZ, scriptOK
	}
	hookName := "on_active_item"
	applyResult := applyWZ
	scriptOK = false
	luax.CallAsync(root, thread, hookName, ch, consume).Then(func(value interface{}) (interface{}, error) {
		scriptOK = true
		vals := luax.ResultValues(value)
		if len(vals) > 0 && vals[0] == lua.LFalse {
			applyResult = false
		}
		return nil, nil
	}).OnError(func(err error) {
		scriptOK = true
		log.Printf("item script failed %s: %v", scriptPath, err)
	})
	return applyResult, scriptOK
}
