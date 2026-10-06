package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/response"
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
	if ch.GetHp() == 0 {
		return false
	}
	wzConsume, ok := consume.GetModel().(*wz.Consume)
	if !ok || wzConsume == nil {
		return false
	}

	applyWZ, err := ch.callActiveConsumeScript(consume)
	if err != nil {
		log.Printf("item script failed %v", err)
		return false
	}

	if applyWZ {
		if wzConsume.CP > 0 || wzConsume.NuffSkillID > 0 {
			ch.useCarnivalItem(wzConsume)
			return true
		}

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
			if character == nil || character.GetHp() == 0 {
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

func (ch *Character) UseCatchItem(slot int16, itemID uint32, mobOID uint32) {
	if itemID/10000 != 227 {
		return
	}
	item := ch.Inventory.Tabs[constant.InventoryTypeConsume].Get(uint8(slot))
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != itemID {
		return
	}
	consume, ok := item.GetModel().(*wz.Consume)
	if ok == false {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	mob := mapInstance.GetMob(mobOID)
	if mob == nil || mob.Wz == nil || mob.Wz.ID != consume.MobID {
		return
	}

	if consume.MobHP > 0 && mob.GetHp() > mob.GetMaxHp()/2 {
		ch.Broadcast(&response.ShowMagnet{MobID: mob.OID, Success: 0}, &ObjectBroadcastOption{})
		return
	}

	spec := ExchangeSpec{Cost: ExchangeSide{Items: map[uint32]uint16{itemID: 1}}}
	if consume.CreateID > 0 {
		spec.Reward.Items = map[uint32]uint16{consume.CreateID: 1}
	}
	if ch.Exchange(spec) != ExchangeOK {
		return
	}

	ch.Broadcast(&response.ShowMagnet{MobID: mob.OID, Success: 1}, &ObjectBroadcastOption{})
	mob.Kill(ch, constant.MobDieAnimationTypeFadeOut)
}

func (ch *Character) UseReturnScroll(ctx actor.Context, slot int16, itemID uint32) error {
	item := ch.Inventory.Tabs[constant.InventoryTypeConsume].Get(uint8(slot))
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != itemID {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	consume, ok := item.GetModel().(*wz.Consume)
	if ok == false || consume.MoveTo <= 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil || m.GetLuaRoot() == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	thread, err := luax.NewThread(m.GetLuaRoot(), constant.CharacterQueryScriptPath)
	if err != nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return err
	}
	allowed, err := luax.Call(thread, "can_use_return_scroll", ch, itemID, consume.MoveTo)
	if err != nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return err
	}
	if allowed != lua.LTrue {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	targetID := uint32(consume.MoveTo)
	if consume.MoveTo == wz.ConsumeMoveToReturnMap {
		targetID = uint32(m.Wz.ReturnMapId)
	}
	target := ch.GameWorld.GetMapSystem().Find(m.StateMachine(), targetID)
	if target == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	err = ch.GameWorld.GetMapSystem().Warp(ctx, ch, target, 0, func(actor.Context) {
		item := ch.Inventory.Tabs[constant.InventoryTypeConsume].Get(uint8(slot))
		if item != nil && item.GetModel().GetID() == itemID {
			ch.Inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
		}
		ch.Listener.OnUpdateStats(ch, nil, true)
	})
	if err != nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return err
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

func (ch *Character) callActiveConsumeScript(consume *Consume) (applyWZ bool, err error) {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return true, nil
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return true, nil
	}

	scriptPath := fmt.Sprintf("script/item/%d.lua", consume.GetModel().GetID())
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return true, nil
	}
	ret, err := luax.Call(thread, "on_active_item", ch, consume)
	if err != nil {
		return false, fmt.Errorf("%s: %w", scriptPath, err)
	}
	return ret != lua.LFalse, nil
}
