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
		ch.Stats.TotalStr() >= required.Str &&
		ch.Stats.TotalDex() >= required.Dex &&
		ch.Stats.TotalInt() >= required.Int &&
		ch.Stats.TotalLuk() >= required.Luk &&
		ch.Stats.Population >= required.Popularity
}

func (ch *Character) Equip(actx actor.Context, slot int16, parts constant.EquipmentPartsType) error {
	before := ch.Inventory.Equipped[parts]
	after, _ := ch.Inventory.GetItem(constant.InventoryTypeEquipment, slot).(Equipment)
	if err := ch.Inventory.Equip(slot, parts); err != nil {
		return err
	}
	ch.onEquipmentChanged(actx, parts, before, after)
	return nil
}

func (ch *Character) Unequip(actx actor.Context, parts constant.EquipmentPartsType, slot int16) error {
	before := ch.Inventory.Equipped[parts]
	if err := ch.Inventory.UnequipToSlot(parts, slot); err != nil {
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
		if item := ch.Inventory.Equipped[parts]; item == nil || item.GetFlags()&constant.ItemFlagLock != 0 {
			return false
		}
		equipment := ch.Inventory.RemoveEquipped(parts)
		if equipment == nil {
			return false
		}
		ch.onEquipmentChanged(actx, parts, equipment, nil)
		dropped = equipment
	} else {
		item := ch.Inventory.GetItem(invType, slot)
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
		if !ch.Inventory.RemoveItem(invType, slot, count) {
			return false
		}
	}
	restricted := !dropped.prepareFieldDrop()
	if restricted && dropped.GetModel().GetTraits()&constant.ItemTraitPartyQuest == 0 {
		return true
	}
	if err := m.SpawnItem(dropped, ItemSpawn{
		Position:             ch.Position,
		From:                 ch.Position,
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
	if amount < constant.DropMesoMin || amount > constant.DropMesoMax || amount > ch.Inventory.Meso {
		return false
	}
	m := ch.GetMap()
	if m == nil {
		return false
	}

	ch.Inventory.SetMeso(ch.Inventory.Meso - amount)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMeso: ch.Inventory.Meso,
	}, true)

	if _, err := m.SpawnMeso(amount, ch.Position, ch.Position, ch.GetID(), constant.DropTypeFFA, true); err != nil {
		log.Printf("Failed to spawn meso on map: %v", err)
	}
	return true
}
