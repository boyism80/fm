package game

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaItemFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"exchange": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			if argc < 2 || argc > 3 {
				L.ArgError(2, "exchange(cost [, reward]) requires 1 or 2 side tables")
				return 0
			}

			var nameToItem func(string) (uint32, bool)
			if ch.gameWorld != nil {
				if resources := ch.gameWorld.GetResources(); resources != nil {
					nameToItem = resources.NameToItem
				}
			}

			cost, ok := parseLuaExchangeSide(L, L.Get(2), 2, nameToItem)
			if !ok {
				return 0
			}
			var reward ExchangeSide
			if argc >= 3 {
				reward, ok = parseLuaExchangeSide(L, L.Get(3), 3, nameToItem)
				if !ok {
					return 0
				}
			}

			result := ch.Exchange(ExchangeSpec{Cost: cost, Reward: reward})
			L.Push(lua.LNumber(result))
			return 1
		},
		"tag_equipment": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.inventory.TagEquipment(int16(L.CheckInt(2)))))
			return 1
		},
		"extend_equipment": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.inventory.ExtendEquipment(int16(L.CheckInt(2)), L.CheckInt(3))))
			return 1
		},
		"lock_item": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.inventory.LockItem(constant.InventoryType(L.CheckInt(2)), int16(L.CheckInt(3)))))
			return 1
		},
		"use_karma": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			scissors := constant.ItemTraitKarma
			if L.OptBool(4, false) {
				scissors = constant.ItemTraitPlatinumKarma
			}
			L.Push(lua.LBool(ch.inventory.UseKarma(constant.InventoryType(L.CheckInt(2)), int16(L.CheckInt(3)), scissors)))
			return 1
		},
		"add_cash": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			nxCash := int32(L.CheckInt(2))
			maplePoint := int32(L.OptInt(3, 0))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.AddCash(cfg.ActorContext, nxCash, maplePoint)
			return 0
		},
		"create_cash_coupons": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			kind := internal.CashCouponKind(L.CheckInt(2))
			value := uint32(L.CheckInt(3))
			count := uint32(L.OptInt(4, 1))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.CreateCashCoupons(cfg.ActorContext, kind, value, count)
			return 0
		},
		"clear_inventory": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "clear_inventory() takes no arguments")
				return 0
			}
			L.Push(lua.LNumber(ch.inventory.ClearInventory()))
			return 1
		},
		"register_card": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cardID := uint32(L.CheckInt(2))
			card, ok := ch.itemModel(cardID).(*wz.Consume)
			if !ok || !card.MonsterBook {
				L.Push(lua.LFalse)
				return 1
			}
			for range L.OptInt(3, 1) {
				ch.monsterBook.Register(cardID)
			}
			L.Push(lua.LTrue)
			return 1
		},
		"reset_monster_book": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.monsterBook.Reset()
			return 0
		},
		"reset_teleport_stones": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.teleportStones.Reset()
			return 0
		},
		"equipped": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			switch L.GetTop() {
			case 1:
				L.ArgError(2, "equipped(part) requires equipment part (e.g. EquipmentPart.Weapon)")
				return 0
			default:
				part := constant.EquipmentPartsType(L.CheckInt(2))
				if eq, ok := ch.inventory.equipped[part]; ok && eq != nil {
					L.Push(luax.NewLuable(L, eq.(luax.Luable)))
					return 1
				}
				L.Push(lua.LNil)
				return 1
			}
		},
		"equip": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() < 2 {
				L.ArgError(2, "equip(slot | item | name) requires one argument")
				return 0
			}
			var slot int16
			switch lv := L.Get(2).(type) {
			case lua.LNumber:
				slot = int16(lv)
			case lua.LString:
				if ch.gameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.gameWorld.GetResources()
				if resources == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				itemID, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LBool(false))
					return 1
				}
				invType, slots := ch.inventory.FindSlots(itemID)
				if invType != constant.InventoryTypeEquipment || len(slots) == 0 {
					L.Push(lua.LBool(false))
					return 1
				}
				slot = slots[0]
			default:
				if itemUD, ok := L.Get(2).(*lua.LUserData); ok && itemUD.Value != nil {
					if item, ok := itemUD.Value.(Item); ok {
						var found bool
						slot, found = ch.inventory.FindSlot(constant.InventoryTypeEquipment, item)
						if !found {
							L.Push(lua.LBool(false))
							return 1
						}
						break
					}
				}
				L.ArgError(2, "equip(slot | item | name): slot (number), item (equipment), or name (string) expected")
				return 0
			}
			item := ch.inventory.tabs[constant.InventoryTypeEquipment].items[slot]
			if item == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			model, ok := item.GetModel().(wz.Equipment)
			if !ok {
				L.Push(lua.LBool(false))
				return 1
			}
			parts := model.GetParts()
			if len(parts) == 0 {
				L.Push(lua.LBool(false))
				return 1
			}
			err := ch.inventory.Equip(slot, parts[0])
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"unequip": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() < 2 {
				L.ArgError(2, "unequip(parts | item | name) requires one argument")
				return 0
			}
			var parts constant.EquipmentPartsType
			switch lv := L.Get(2).(type) {
			case lua.LNumber:
				parts = constant.EquipmentPartsType(lv)
			case lua.LString:
				if ch.gameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.gameWorld.GetResources()
				if resources == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				itemID, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LBool(false))
					return 1
				}
				var found bool
				for p, eq := range ch.inventory.equipped {
					if eq != nil && eq.GetModel().GetID() == itemID {
						parts = p
						found = true
						break
					}
				}
				if !found {
					L.Push(lua.LBool(false))
					return 1
				}
			default:
				if itemUD, ok := L.Get(2).(*lua.LUserData); ok && itemUD.Value != nil {
					if item, ok := itemUD.Value.(Item); ok {
						var found bool
						for p, eq := range ch.inventory.equipped {
							if eq == item {
								parts = p
								found = true
								break
							}
						}
						if !found {
							L.Push(lua.LBool(false))
							return 1
						}
						break
					}
				}
				L.ArgError(2, "unequip(parts | item | name): parts (EquipmentPart), item (equipment), or name (string) expected")
				return 0
			}
			if L.GetTop() >= 3 && !L.ToBool(3) {
				L.Push(lua.LBool(ch.inventory.RemoveEquipped(parts) != nil))
				return 1
			}
			err := ch.inventory.Unequip(parts)
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"empty_slots": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(2, "empty_slots(inventory_type) requires inventory type")
				return 0
			}
			inven := ch.inventory.tabs[constant.InventoryType(L.CheckInt(2))]
			if inven == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(inven.EmptySlotCount()))
			return 1
		},
		"put_key": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 4 {
				L.ArgError(2, "put_key(slot, type, action) requires slot, type, and action")
				return 0
			}
			ch.keyLayout.SetKey(L.CheckInt(2), byte(L.CheckInt(3)), int32(L.CheckInt(4)))
			ch.listener.OnKeyMap(ch)
			return 0
		},
		"use_item": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(2, "use_item(item_id) requires an item id")
				return 0
			}
			L.Push(lua.LBool(ch.UseItemEffect(uint32(L.CheckInt(2)))))
			return 1
		},
		"sync_item": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 3 {
				L.ArgError(2, "sync_item(inventory_type, slot)")
				return 0
			}
			invType := constant.InventoryType(L.CheckInt(2))
			slot := int16(L.CheckInt(3))
			item := ch.inventory.GetItem(invType, slot)
			if item == nil {
				return 0
			}
			ch.listener.OnInventorySlotUpdated(ch, invType, slot, item)
			return 0
		},
		"mob_drops": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(2, "mob_drops(mob_id) requires a mob id")
				return 0
			}
			t := L.NewTable()
			if ch.gameWorld == nil {
				L.Push(t)
				return 1
			}
			resources := ch.gameWorld.GetResources()
			if resources == nil {
				L.Push(t)
				return 1
			}
			mobDrops := resources.MobDrops[uint32(L.CheckInt(2))]
			for i, d := range mobDrops {
				row := L.NewTable()
				row.RawSetString("item", lua.LNumber(d.Item))
				row.RawSetString("money", lua.LNumber(d.Money))
				row.RawSetString("prob", lua.LNumber(d.Prob))
				row.RawSetString("min", lua.LNumber(d.Min))
				row.RawSetString("max", lua.LNumber(d.Max))
				row.RawSetString("quest", lua.LNumber(resources.QuestItems[d.Item]))
				t.RawSetInt(i+1, row)
			}
			L.Push(t)
			return 1
		},
		"item": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			switch L.GetTop() {
			case 2:
				var itemID uint32
				switch lv := L.Get(2).(type) {
				case lua.LNumber:
					itemID = uint32(lv)
				case lua.LString:
					if ch.gameWorld == nil {
						L.Push(L.NewTable())
						return 1
					}
					resources := ch.gameWorld.GetResources()
					if resources == nil {
						L.Push(L.NewTable())
						return 1
					}
					id, ok := resources.NameToItem(string(lv))
					if !ok {
						L.Push(L.NewTable())
						return 1
					}
					itemID = id
				default:
					L.ArgError(2, "item id (number) or item name (string) expected")
					return 0
				}
				slotTbl := L.NewTable()
				model := ch.itemModel(itemID)
				if model == nil {
					L.Push(slotTbl)
					return 1
				}
				if inven := ch.inventory.tabs[model.GetInventoryType()]; inven != nil {
					for slot, item := range inven.items {
						if item != nil && item.GetModel().GetID() == itemID {
							L.Push(luax.NewLuable(L, item))
							slotTbl.RawSet(lua.LNumber(slot), L.Get(-1))
							L.Pop(1)
						}
					}
				}
				L.Push(slotTbl)
				return 1
			case 3:
				invType := constant.InventoryType(L.CheckInt(2))
				slot := int16(L.CheckInt(3))
				item := ch.inventory.GetItem(invType, slot)
				if item == nil {
					L.Push(lua.LNil)
				} else {
					L.Push(luax.NewLuable(L, item))
				}
				return 1
			default:
				L.ArgError(2, "item(item_id) or item(inven_type, slot) expected")
				return 0
			}
		},
		"items": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				result := L.NewTable()
				for invType, inven := range ch.inventory.tabs {
					if inven == nil {
						continue
					}
					invTbl := L.NewTable()
					for slot, item := range inven.items {
						if item == nil {
							continue
						}
						L.Push(luax.NewLuable(L, item))
						invTbl.RawSet(lua.LNumber(slot), L.Get(-1))
						L.Pop(1)
					}
					result.RawSet(lua.LNumber(invType), invTbl)
				}
				L.Push(result)
				return 1
			default:
				invType := constant.InventoryType(L.CheckInt(2))
				slotTbl := L.NewTable()
				if inven := ch.inventory.tabs[invType]; inven != nil {
					for slot, item := range inven.items {
						if item == nil {
							continue
						}
						L.Push(luax.NewLuable(L, item))
						slotTbl.RawSet(lua.LNumber(slot), L.Get(-1))
						L.Pop(1)
					}
				}
				L.Push(slotTbl)
				return 1
			}
		},
		"mkitem": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.gameWorld.GetResources()
			if resources == nil {
				L.Push(lua.LNil)
				return 1
			}

			argc := L.GetTop()
			count := uint16(1)
			isRandomizeStats := false
			switch argc {
			case 1:
				L.ArgError(2, "mkitem(itemIdOrName [, count] [, isRandomizeStats]) requires at least one argument")
				return 0
			case 2:
			case 3:
				if n := int(L.CheckInt(3)); n >= 1 {
					count = uint16(n)
				}
			default:
				if n := L.CheckInt(3); n >= 1 {
					count = uint16(n)
				}
				if argc >= 4 {
					isRandomizeStats = lua.LVAsBool(L.Get(4))
				}
			}

			var itemId uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LNil)
					return 1
				}
				itemId = id
			case lua.LNumber:
				itemId = uint32(lv)
			default:
				L.ArgError(2, "item id (number) or item name (string) expected")
				return 0
			}

			if _, ok := resources.Items[itemId]; !ok {
				L.Push(lua.LNil)
				return 1
			}

			if resources.Items[itemId] == nil {
				L.Push(lua.LNil)
				return 1
			}
			item, err := NewItem(itemId, count, ch.gameWorld)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			if isRandomizeStats {
				if eq, ok := item.(Equipment); ok {
					if em, ok := eq.GetModel().(wz.Equipment); ok {
						eq.GetEquipmentCore().RandomizeStats(em)
					}
				}
			}

			added, err := ch.inventory.AddItem(item, true)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}

			L.Push(luax.NewLuable(L, added[0]))
			return 1
		},
		"rmitem": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.ArgError(2, "rmitem(itemIdOrName [, count]) or rmitem(inv_type, slot, count)")
				return 0
			case 4:
				invType := constant.InventoryType(L.CheckInt(2))
				slot := int16(L.CheckInt(3))
				count := uint16(L.CheckInt(4))
				if count == 0 {
					L.Push(lua.LBool(false))
					return 1
				}
				ok := ch.inventory.RemoveItem(invType, slot, count)
				L.Push(lua.LBool(ok))
				return 1
			case 2, 3:
				if ch.gameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.gameWorld.GetResources()
				if resources == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				var itemId uint32
				switch lv := L.Get(2).(type) {
				case lua.LString:
					id, ok := resources.NameToItem(string(lv))
					if !ok {
						L.Push(lua.LBool(false))
						return 1
					}
					itemId = id
				case lua.LNumber:
					itemId = uint32(lv)
				default:
					L.ArgError(2, "item id (number) or item name (string) expected")
					return 0
				}
				count := uint16(1)
				if argc == 3 {
					if n := L.CheckInt(3); n < 1 {
						L.Push(lua.LBool(false))
						return 1
					}
					count = uint16(L.CheckInt(3))
				}
				removed := ch.inventory.RemoveByItemIDCount(itemId, count)
				L.Push(lua.LBool(removed))
				return 1
			default:
				L.ArgError(2, "rmitem(itemIdOrName [, count]) or rmitem(inv_type, slot, count)")
				return 0
			}
		},
	}
}
