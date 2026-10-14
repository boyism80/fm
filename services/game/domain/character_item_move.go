package domain

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
)

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
	root := ch.GetMap().GetLuaRoot()
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

	var dropped Item
	if slot < 0 {
		parts := constant.EquipmentPartsType(slot)
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
		count = min(count, item.GetCount())
		if count == 0 {
			return false
		}
		dropped = item.Clone(count)
		if !ch.Inventory.RemoveItem(invType, slot, count) {
			return false
		}
	}

	m := ch.GetMap()
	if m == nil {
		return true
	}
	if err := m.SpawnItem(dropped, ItemSpawn{
		Position:   ch.Position,
		From:       ch.Position,
		Owner:      ch.GetID(),
		DropType:   constant.DropTypeFFA,
		PlayerDrop: true,
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

	ch.Inventory.SetMeso(ch.Inventory.Meso - amount)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMeso: ch.Inventory.Meso,
	}, true)

	m := ch.GetMap()
	if m == nil {
		return true
	}
	if _, err := m.SpawnMeso(amount, ch.Position, ch.Position, ch.GetID(), constant.DropTypeFFA, true); err != nil {
		log.Printf("Failed to spawn meso on map: %v", err)
	}
	return true
}
