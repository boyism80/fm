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
		"point_below": func(L *lua.LState) int {
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
			out := mapInstance.PointBelow(types.Point[int16]{X: x, Y: y})
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
			tbl := L.NewTable()
			if L.GetTop() >= 2 {
				mobWZID := uint32(L.CheckInt(2))
				limit := 0
				if v, ok := L.Get(3).(lua.LNumber); ok {
					limit = int(v)
				}
				for _, mob := range mapInstance.MobsByTemplate(mobWZID, limit) {
					tbl.RawSetInt(int(mob.OID), luax.NewLuable(L, mob))
				}
			} else {
				for _, mob := range mapInstance.GetMobs() {
					if mobObj, ok := mob.(*Mob); ok {
						tbl.RawSetInt(int(mobObj.OID), luax.NewLuable(L, mobObj))
					}
				}
			}
			L.Push(tbl)
			return 1
		},
		"mob_by_template": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			mobID := uint32(L.CheckInt(2))
			mobs := mapInstance.MobsByTemplate(mobID, 1)
			if len(mobs) == 0 {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mobs[0]))
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
			wz := mapInstance.Wz
			if wz == nil {
				L.Push(lua.LNil)
				return 1
			}
			tbl := L.NewTable()
			tbl.RawSetString("id", lua.LNumber(wz.ID))
			tbl.RawSetString("name", lua.LString(wz.Name))
			tbl.RawSetString("return_map_id", lua.LNumber(wz.ReturnMapId))
			tbl.RawSetString("town", lua.LBool(wz.IsTown))
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
			dropType := constant.DropTypeFFA
			ownerID := uint32(0)
			if owner != nil {
				if owner.GetPartyID() != nil {
					dropType = constant.DropTypeParty
				} else {
					dropType = constant.DropTypeOwnerOnly
				}
				ownerID = owner.GetID()
			}
			meso, err := mapInstance.SpawnMeso(count, pos, pos, ownerID, dropType, false)
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
			dropType := constant.DropTypeFFA
			ownerID := uint32(0)
			if owner != nil {
				if owner.GetPartyID() != nil {
					dropType = constant.DropTypeParty
				} else {
					dropType = constant.DropTypeOwnerOnly
				}
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
			spawnType := constant.MobSpawnTypeAnimate
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
		"property": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			key := L.CheckString(2)
			if L.GetTop() == 2 {
				val, exists := mapInstance.GetProperty(key)
				if !exists {
					L.Push(lua.LNil)
					return 1
				}
				switch v := val.(type) {
				case bool:
					L.Push(lua.LBool(v))
				case int:
					L.Push(lua.LNumber(v))
				case int32:
					L.Push(lua.LNumber(v))
				case int64:
					L.Push(lua.LNumber(v))
				case uint32:
					L.Push(lua.LNumber(v))
				case float64:
					L.Push(lua.LNumber(v))
				case string:
					L.Push(lua.LString(v))
				default:
					L.Push(lua.LNil)
				}
				return 1
			}
			val := L.Get(3)
			switch val.Type() {
			case lua.LTNil:
				mapInstance.SetProperty(key, nil)
			case lua.LTBool:
				mapInstance.SetProperty(key, lua.LVAsBool(val))
			case lua.LTNumber:
				mapInstance.SetProperty(key, float64(lua.LVAsNumber(val)))
			case lua.LTString:
				mapInstance.SetProperty(key, string(val.(lua.LString)))
			default:
				L.ArgError(3, "property value must be nil, bool, number, or string")
				return 0
			}
			return 0
		},
		"respawn": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.GameWorld == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			includeNegativeMobTime := false
			if L.GetTop() >= 2 {
				includeNegativeMobTime = lua.LVAsBool(L.Get(2))
			}
			cfg, _ := luax.GetConfiguration(L)
			return mapInstance.GameWorld.GetMapSystem().RespawnFromLua(L, mapInstance, cfg.ActorContext, includeNegativeMobTime)
		},
		"block_gen": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			argc := L.GetTop()
			if argc < 2 {
				L.ArgError(2, "block_gen(enabled) or block_gen(enabled, mobId) requires enabled")
				return 0
			}
			enabled := lua.LVAsBool(L.Get(2))
			if argc == 2 {
				mapInstance.SetAllMobGenEnabled(enabled)
				return 0
			}
			if argc == 3 {
				mobWZID := uint32(L.CheckInt(3))
				mapInstance.SetMobGenEnabled(mobWZID, enabled)
				return 0
			}
			L.ArgError(4, "block_gen(enabled) or block_gen(enabled, mobId)")
			return 0
		},
		"remove_npc": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			arg := L.Get(2)
			var oid uint32
			switch v := arg.(type) {
			case *lua.LUserData:
				npc, ok := v.Value.(*Npc)
				if !ok || npc == nil {
					L.ArgError(2, "Npc or npc OID expected")
					return 0
				}
				oid = npc.OID
			case lua.LNumber:
				oid = uint32(v)
			default:
				L.ArgError(2, "Npc or npc OID expected")
				return 0
			}
			if oid == 0 {
				return 0
			}
			err := mapInstance.RemoveNpc(oid)
			if err != nil {
				L.RaiseError("remove_npc: %v", err)
				return 0
			}
			return 0
		},
		"reactor": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			reactorID := uint32(L.CheckInt(2))
			reactor := mapInstance.ReactorByTemplate(reactorID)
			if reactor == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, reactor))
			return 1
		},
		"reactor_by_name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			name := L.CheckString(2)
			reactor := mapInstance.ReactorByName(name)
			if reactor == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, reactor))
			return 1
		},
		"kill_all_mobs": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			animType := constant.MobDieAnimationTypeFadeOut
			if L.GetTop() >= 2 {
				animType = constant.MobDieAnimationType(L.CheckInt(2))
			}
			mapInstance.KillAllMonsters(animType)
			return 0
		},
		"remove_mob": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			oid := uint32(L.CheckInt(2))
			animType := constant.MobDieAnimationTypeFadeOut
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
		"music": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			song := L.CheckString(2)
			mapInstance.ChangeMusic(song)
			return 0
		},
		"message": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			message := L.CheckString(2)
			mapInstance.MapMessage(message)
			return 0
		},
		"clear_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			mapInstance.ClearEffect()
			return 0
		},
		"show_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			path := L.CheckString(2)
			mapInstance.ShowEffect(path)
			return 0
		},
		"play_sound": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			path := L.CheckString(2)
			mapInstance.PlaySound(path)
			return 0
		},
		"players_in_area": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			index := L.CheckInt(2)
			L.Push(lua.LNumber(mapInstance.PlayersInArea(index)))
			return 1
		},
		"reload_reactors": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			count := mapInstance.ReloadReactors()
			L.Push(lua.LNumber(count))
			return 1
		},
		"reset": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.GameWorld == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			cfg, _ := luax.GetConfiguration(L)
			return mapInstance.GameWorld.GetMapSystem().ResetFromLua(L, mapInstance, cfg.ActorContext)
		},
		"portal": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			name := L.CheckString(2)
			portal := mapInstance.FindPortalByName(name)
			if portal == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, portal))
			return 1
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
