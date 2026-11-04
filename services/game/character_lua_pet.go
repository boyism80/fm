package game

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaPetFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"pet_id": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.GetModel().GetID()))
			return 1
		},
		"pet_closeness": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.closeness))
			return 1
		},
		"add_pet_closeness": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.AddCloseness(L.CheckInt(2))
			L.Push(lua.LTrue)
			return 1
		},
		"set_pet_closeness": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.SetCloseness(L.CheckInt(2))
			L.Push(lua.LTrue)
			return 1
		},
		"set_pet_fullness": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.SetFullness(L.CheckInt(2))
			L.Push(lua.LTrue)
			return 1
		},
		"hunger_pet": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.hunger()
			L.Push(lua.LTrue)
			return 1
		},
		"pet_level": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.level))
			return 1
		},
		"pet_fullness": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.fullness))
			return 1
		},
		"rename_pet": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.pets.active.Rename(L.CheckString(2)) == nil))
			return 1
		},
		"feed_pet_cash": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.pets.active.FeedCash(uint32(L.CheckInt(2))) == nil))
			return 1
		},
		"change_pet_skill": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			sn := uint64(L.CheckNumber(2))
			L.Push(lua.LBool(ch.pets.ChangeSkill(sn, uint32(L.CheckInt(3))) == nil))
			return 1
		},
		"expired_pets": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			tbl := L.NewTable()
			for slot, pet := range ch.pets.Expired() {
				entry := L.NewTable()
				entry.RawSetString("slot", lua.LNumber(slot))
				entry.RawSetString("id", lua.LNumber(pet.GetModel().GetID()))
				entry.RawSetString("name", lua.LString(pet.name))
				tbl.Append(entry)
			}
			L.Push(tbl)
			return 1
		},
		"revive_pet": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.pets.Revive(int16(L.CheckInt(2))) == nil))
			return 1
		},
	}
}
