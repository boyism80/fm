package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) ApplyConsumeEffect(consume *Consume) bool {
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
			character.applyConsumeCureDebuffs(wzConsume)
			character.applyConsumeBuff(wzConsume)
			character.applyConsumeRecovery(wzConsume)
			character.applyConsumeExp(wzConsume)
		}
	}

	return true
}

func (ch *Character) applyConsumeCureDebuffs(consumeItem *wz.Consume) {
	if ch == nil || consumeItem == nil || len(consumeItem.CureDebuffs) == 0 {
		return
	}
	ch.RemoveDebuff(consumeItem.CureDebuffs...)
}

func (ch *Character) applyConsumeExp(consumeItem *wz.Consume) {
	if ch == nil || consumeItem == nil || consumeItem.ExpInc <= 0 {
		return
	}
	ch.AddExp(uint32(consumeItem.ExpInc))
}

func (ch *Character) applyConsumeBuff(consumeItem *wz.Consume) bool {
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

func (ch *Character) applyConsumeRecovery(consumeItem *wz.Consume) bool {
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
			stats[constant.STAT_HP] = int32(ch.GetHp())
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
			stats[constant.STAT_MP] = int32(ch.GetMp())
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
	hookName := fmt.Sprintf("on_active_item_%d", itemID)
	ret, err := luax.Call(thread, hookName, ch, consume)
	if err != nil {
		log.Printf("item script failed %s: %v", scriptPath, err)
		return true, false
	}
	if ret != nil && ret == lua.LFalse {
		return false, true
	}
	return true, true
}
