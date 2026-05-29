package entity

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

func (m *Map) LuaTypeName() string {
	return "LuaMap"
}

func (m *Map) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"foothold_point": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			posTbl := L.CheckTable(2)
			var x, y int16
			if lx := posTbl.RawGetInt(1); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			} else if lx := posTbl.RawGetString("x"); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			}
			if ly := posTbl.RawGetInt(2); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			} else if ly := posTbl.RawGetString("y"); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			}
			out := mapInstance.FootholdPoint(types.Point[int16]{X: x, Y: y})
			if out == nil {
				L.Push(lua.LNil)
				return 1
			}
			res := L.NewTable()
			res.RawSetString("x", lua.LNumber(out.X))
			res.RawSetString("y", lua.LNumber(out.Y))
			L.Push(res)
			return 1
		},
		"npcs": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			npcs := mapInstance.GetNpcs()
			tbl := L.NewTable()
			for _, npc := range npcs {
				if npcObj, ok := npc.(*Npc); ok {
					tbl.RawSetInt(int(npcObj.OID), luax.NewLuable(L, npcObj))
				}
			}
			L.Push(tbl)
			return 1
		},
		"mobs": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			mobs := mapInstance.GetMobs()
			tbl := L.NewTable()
			for _, mob := range mobs {
				if mobObj, ok := mob.(*Mob); ok {
					tbl.RawSetInt(int(mobObj.OID), luax.NewLuable(L, mobObj))
				}
			}
			L.Push(tbl)
			return 1
		},
		"objects": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			filter := constant.ObjectTypeObject
			if L.GetTop() >= 2 {
				if lv := L.Get(2); lv.Type() == lua.LTNumber {
					filter = constant.ObjectType(lua.LVAsNumber(lv))
				}
			}
			objs := mapInstance.GetObjects(filter)
			result := L.NewTable()
			idx := 0
			for _, v := range objs {
				if luable, ok := v.(luax.Luable); ok {
					idx++
					result.RawSetInt(idx, luax.NewLuable(L, luable))
				}
			}
			L.Push(result)
			return 1
		},
		"characters": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			players := mapInstance.GetAllPlayers()
			tbl := L.NewTable()
			for _, player := range players {
				if char, ok := player.(*Character); ok {
					tbl.RawSetInt(int(char.GetID()), luax.NewLuable(L, char))
				}
			}
			L.Push(tbl)
			return 1
		},
		"items": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			items := mapInstance.GetItems()
			tbl := L.NewTable()
			for _, item := range items {
				if itemObj, ok := item.(Item); ok {
					fp := itemObj.GetFieldPlacement()
					if fp != nil && fp.ObjectCore != nil {
						tbl.RawSetInt(int(fp.OID), luax.NewLuable(L, itemObj))
					}
				}
			}
			L.Push(tbl)
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			spec := mapInstance.Wz
			if spec == nil {
				L.Push(lua.LNil)
				return 1
			}
			tbl := L.NewTable()
			tbl.RawSetString("id", lua.LNumber(spec.ID))
			tbl.RawSetString("name", lua.LString(spec.Name))
			tbl.RawSetString("return_map_id", lua.LNumber(spec.ReturnMapId))
			tbl.RawSetString("town", lua.LBool(spec.IsTown))
			L.Push(tbl)
			return 1
		},
		"recovery_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			L.Push(lua.LNumber(mapInstance.GetRecoveryRate()))
			return 1
		},
		"spawn_meso": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			count := int32(L.CheckInt(2))
			posTbl := L.CheckTable(3)
			var x, y int16
			if lx := posTbl.RawGetInt(1); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			} else if lx := posTbl.RawGetString("x"); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			}
			if ly := posTbl.RawGetInt(2); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			} else if ly := posTbl.RawGetString("y"); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			}
			var owner *Character
			if ownerLV := L.Get(4); ownerLV != lua.LNil {
				if ownerUd, ok := ownerLV.(*lua.LUserData); ok {
					if ch, ok := ownerUd.Value.(*Character); ok {
						owner = ch
					}
				}
			}
			if count <= 0 {
				return 0
			}
			pos := types.Point[int16]{X: x, Y: y}
			dropType := constant.DROP_TYPE_FFA
			ownerID := uint32(0)
			if owner != nil {
				dropType = constant.DROP_TYPE_OWNED
				ownerID = owner.GetID()
			}
			meso, err := mapInstance.SpawnMeso(count, pos, ownerID, dropType)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, meso))
			return 1
		},
		"spawn_item": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.GameWorld == nil {
				return 0
			}
			resources := mapInstance.GameWorld.GetResources()
			if resources == nil {
				return 0
			}

			argc := L.GetTop()
			if argc < 4 {
				L.ArgError(2, "spawn_item(itemIdOrName, count, position [, owner] [, isRandomizeStats]) requires item, count, and position")
				return 0
			}

			var itemId uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToItem(string(lv))
				if !ok {
					return 0
				}
				itemId = id
			case lua.LNumber:
				itemId = uint32(lv)
			default:
				L.ArgError(2, "item id (number) or item name (string) expected")
				return 0
			}

			if _, ok := resources.Items[itemId]; !ok {
				return 0
			}

			count := uint16(1)
			if n := L.CheckInt(3); n >= 1 {
				count = uint16(n)
			}

			posTbl := L.CheckTable(4)
			var x, y int16
			if lx := posTbl.RawGetInt(1); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			} else if lx := posTbl.RawGetString("x"); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			}
			if ly := posTbl.RawGetInt(2); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			} else if ly := posTbl.RawGetString("y"); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			}
			pos := types.Point[int16]{X: x, Y: y}

			var owner *Character
			if argc >= 5 && L.Get(5) != lua.LNil {
				if ownerUd, ok := L.Get(5).(*lua.LUserData); ok {
					if ch, ok := ownerUd.Value.(*Character); ok {
						owner = ch
					}
				}
			}
			isRandomizeStats := false
			if argc >= 6 {
				isRandomizeStats = lua.LVAsBool(L.Get(6))
			}

			item, err := NewItem(itemId, count, mapInstance.GameWorld)
			if err != nil {
				return 0
			}
			if isRandomizeStats {
				if eq, ok := item.(Equipment); ok {
					if em, ok := eq.GetModel().(wz.Equipment); ok {
						eq.GetEquipmentCore().RandomizeStats(em)
					}
				}
			}
			dropType := constant.DROP_TYPE_FFA
			ownerID := uint32(0)
			if owner != nil {
				dropType = constant.DROP_TYPE_OWNED
				ownerID = owner.GetID()
			}
			fp := &FieldPlacement{
				ObjectCore:   &ObjectCore{Position: pos},
				Owner:        ownerID,
				SpawnedPoint: pos,
				DropType:     dropType,
			}
			fp.ObjectCore.self = fp
			item.BindFieldPlacement(fp)

			if err := mapInstance.SpawnItem(item, ownerID, dropType); err != nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, item))
			return 1
		},
		"spawn_mob": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.GameWorld == nil {
				L.RaiseError("spawn_mob: map has no context")
				return 0
			}
			resources := mapInstance.GameWorld.GetResources()
			if resources == nil {
				L.RaiseError("spawn_mob: no resources")
				return 0
			}
			var mobID uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToMob(string(lv))
				if !ok {
					L.Push(lua.LNil)
					return 1
				}
				mobID = id
			case lua.LNumber:
				mobID = uint32(lv)
			default:
				L.ArgError(2, "mob id (number) or mob name (string) expected")
				return 0
			}
			x := int16(L.CheckInt(3))
			y := int16(L.CheckInt(4))
			pos := types.Point[int16]{X: x, Y: y}
			spawnType := constant.MOB_SPAWN_TYPE_ANIMATE
			link := uint32(0)
			if L.GetTop() >= 5 {
				spawnType = constant.MobSpawnType(L.CheckInt(5))
			}
			if L.GetTop() >= 6 {
				link = uint32(L.CheckInt(6))
			}
			mob, err := mapInstance.SpawnMob(mobID, pos, nil, spawnType, link)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mob))
			return 1
		},
		"spawn_npc": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.GameWorld == nil {
				L.RaiseError("spawn_npc: map has no context")
				return 0
			}
			resources := mapInstance.GameWorld.GetResources()
			if resources == nil {
				L.RaiseError("spawn_npc: no resources")
				return 0
			}
			var npcID uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToNpc(string(lv))
				if !ok {
					L.RaiseError("spawn_npc: unknown npc name %q", string(lv))
					return 0
				}
				npcID = id
			case lua.LNumber:
				npcID = uint32(lv)
			default:
				L.ArgError(2, "npc id (number) or npc name (string) expected")
				return 0
			}
			x := int16(L.CheckInt(3))
			y := int16(L.CheckInt(4))
			pos := types.Point[int16]{X: x, Y: y}
			npc, err := mapInstance.SpawnNpc(npcID, pos)
			if err != nil {
				L.RaiseError("spawn_npc: %v", err)
				return 0
			}
			L.Push(luax.NewLuable(L, npc))
			return 1
		},
		"remove_mob": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			oid := uint32(L.CheckInt(2))
			animType := constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT
			if L.GetTop() >= 3 {
				animType = constant.MobDieAnimationType(L.CheckInt(3))
			}
			err := mapInstance.RemoveMob(oid, animType)
			if err != nil {
				L.RaiseError("remove_mob: %v", err)
				return 0
			}
			return 0
		},
		"remove_mist": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			arg := L.Get(2)
			switch v := arg.(type) {
			case *lua.LUserData:
				mist, ok := v.Value.(*Mist)
				if !ok || mist == nil {
					L.ArgError(2, "Mist or mist OID expected")
					return 0
				}
				if mist.GetMap() != mapInstance || mist.OID == 0 {
					return 0
				}
				mapInstance.RemoveMist(mist.OID)
			case lua.LNumber:
				oid := uint32(v)
				if oid == 0 {
					return 0
				}
				mapInstance.RemoveMist(oid)
			default:
				L.ArgError(2, "Mist or mist OID expected")
				return 0
			}
			return 0
		},
	}
}

func (m *Map) String() string {
	return m.LuaTypeName()
}

func (m *Map) Type() lua.LValueType {
	return lua.LTUserData
}
