package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (g *Guild) LuaTypeName() string {
	return "LuaGuild"
}

func (g *Guild) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.GuildID))
			return 1
		},
		"world_id": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "world_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.WorldID))
			return 1
		},
		"name": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "name() is read-only")
				return 0
			}
			L.Push(lua.LString(guild.Name))
			return 1
		},
		"leader_id": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "leader_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.LeaderCharacterID))
			return 1
		},
		"revision": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "revision() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.Revision))
			return 1
		},
		"gp": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "gp() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.GP))
			return 1
		},
		"capacity": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "capacity() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.Capacity))
			return 1
		},
		"notice": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "notice() is read-only")
				return 0
			}
			L.Push(lua.LString(guild.Notice))
			return 1
		},
		"members": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "members() takes no arguments")
				return 0
			}
			tbl := L.NewTable()
			for i, m := range guild.Members {
				if m == nil {
					continue
				}
				tbl.RawSetInt(i+1, luax.NewLuable(L, m))
			}
			L.Push(tbl)
			return 1
		},
		"member": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(3, "member(character) takes exactly one argument")
				return 0
			}
			ch, ok := LuaCheckCharacterArg(L, 2)
			if !ok {
				return 0
			}
			m := guild.findMember(ch.GetID())
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, m))
			return 1
		},
		"rank": func(L *lua.LState) int {
			guild, ok := LuaCheckGuild(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(3, "rank(character) takes exactly one argument")
				return 0
			}
			ch, ok := LuaCheckCharacterArg(L, 2)
			if !ok {
				return 0
			}
			m := guild.findMember(ch.GetID())
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(guildRankToUint32(m.Rank)))
			return 1
		},
		"disband": func(L *lua.LState) int {
			ch, ok := LuaCheckCharacterArg(L, 2)
			if !ok {
				return 0
			}
			return luaGuildDisband(L, ch)
		},
		"inc_capacity": func(L *lua.LState) int {
			ch, ok := LuaCheckCharacterArg(L, 2)
			if !ok {
				return 0
			}
			extendedCap := L.CheckBool(3)
			return luaGuildIncCapacity(L, ch, extendedCap)
		},
	}
}

func (g *Guild) findMember(characterID uint32) *GuildMember {
	if g == nil {
		return nil
	}
	for _, m := range g.Members {
		if m == nil || m.CharacterID != characterID {
			continue
		}
		return m
	}
	return nil
}

func LuaCheckGuild(L *lua.LState, idx int) (*Guild, bool) {
	ud := L.CheckUserData(idx)
	guild, ok := ud.Value.(*Guild)
	if !ok || guild == nil {
		L.ArgError(idx, "Guild expected")
		return nil, false
	}
	return guild, true
}

func LuaCheckCharacterArg(L *lua.LState, idx int) (*Character, bool) {
	ud := L.CheckUserData(idx)
	ch, ok := ud.Value.(*Character)
	if !ok || ch == nil {
		L.ArgError(idx, "Character expected")
		return nil, false
	}
	return ch, true
}

func (g *Guild) String() string {
	return g.LuaTypeName()
}

func (g *Guild) Type() lua.LValueType {
	return lua.LTUserData
}
