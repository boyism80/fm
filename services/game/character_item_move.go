package game

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

func (ch *Character) canEquip(required wz.RequiredStats) bool {
	return ch.GetLevel() >= required.Level &&
		ch.stats.TotalStr() >= required.Str &&
		ch.stats.TotalDex() >= required.Dex &&
		ch.stats.TotalInt() >= required.Int &&
		ch.stats.TotalLuk() >= required.Luk &&
		ch.stats.population >= required.Popularity
}

func (ch *Character) Equip(actx actor.Context, slot int16, parts constant.EquipmentPartsType) error {
	before := ch.inventory.equipped[parts]
	after, _ := ch.inventory.GetItem(constant.InventoryTypeEquipment, slot).(Equipment)
	ch.damageHistory.changing()
	if err := ch.inventory.Equip(slot, parts); err != nil {
		return err
	}
	ch.onEquipmentChanged(actx, parts, before, after)
	return nil
}

func (ch *Character) Unequip(actx actor.Context, parts constant.EquipmentPartsType, slot int16) error {
	before := ch.inventory.equipped[parts]
	ch.damageHistory.changing()
	if err := ch.inventory.UnequipToSlot(parts, slot); err != nil {
		return err
	}
	ch.onEquipmentChanged(actx, parts, before, nil)
	return nil
}

func (ch *Character) onEquipmentChanged(actx actor.Context, parts constant.EquipmentPartsType, before, after Equipment) {
	m := ch.GetMap()
	if m == nil {
		return
	}
	root := m.GetLuaRoot()
	if root == nil {
		return
	}
	var beforeArg, afterArg interface{}
	if before != nil {
		beforeArg = luax.NewLuable(root, before.(luax.Luable))
	}
	if after != nil {
		afterArg = luax.NewLuable(root, after.(luax.Luable))
	}
	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_equipment_changed: %v", err)
		return
	}
	luax.CallAsync(actx, root, thread, "on_equipment_changed", ch, int32(parts), beforeArg, afterArg).OnError(func(err error) {
		log.Printf("Failed to call script on_equipment_changed: %v", err)
	})
}

func (ch *Character) DropItem(actx actor.Context, invType constant.InventoryType, slot int16, count uint16) bool {
	if ch.Spectating() {
		return false
	}
	m := ch.GetMap()
	if m == nil {
		return false
	}

	var dropped Item
	if slot < 0 {
		parts := constant.EquipmentPartsType(slot)
		if item := ch.inventory.equipped[parts]; item == nil || item.GetFlags()&constant.ItemFlagLock != 0 {
			return false
		}
		equipment := ch.inventory.RemoveEquipped(parts)
		if equipment == nil {
			return false
		}
		ch.onEquipmentChanged(actx, parts, equipment, nil)
		dropped = equipment
	} else {
		item := ch.inventory.GetItem(invType, slot)
		if item == nil {
			return false
		}
		if _, ok := item.(*Pet); ok {
			return false
		}
		if item.GetFlags()&constant.ItemFlagLock != 0 {
			return false
		}
		count = min(count, item.GetCount())
		if count == 0 {
			return false
		}
		dropped = item.Clone(count)
		if !ch.inventory.RemoveItem(invType, slot, count) {
			return false
		}
	}
	restricted := !dropped.prepareFieldDrop()
	if restricted && dropped.GetModel().GetTraits()&constant.ItemTraitPartyQuest == 0 {
		return true
	}
	if err := m.SpawnItem(dropped, ItemSpawn{
		Position:             ch.position,
		From:                 ch.position,
		Owner:                ch.GetID(),
		DropType:             constant.DropTypeFFA,
		PlayerDrop:           true,
		RestrictedPlayerDrop: restricted,
	}); err != nil {
		log.Printf("Failed to spawn item on map: %v", err)
	}
	return true
}

func (ch *Character) DropMeso(amount int32) bool {
	if ch.Spectating() {
		return false
	}
	if amount < constant.DropMesoMin || amount > constant.DropMesoMax || amount > ch.inventory.meso {
		return false
	}
	m := ch.GetMap()
	if m == nil {
		return false
	}

	ch.inventory.SetMeso(ch.inventory.meso - amount)
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMeso: ch.inventory.meso,
	}, true)

	if _, err := m.SpawnMeso(amount, ch.position, ch.position, ch.GetID(), constant.DropTypeFFA, true); err != nil {
		log.Printf("Failed to spawn meso on map: %v", err)
	}
	return true
}
