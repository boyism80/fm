package game

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	lua "github.com/yuin/gopher-lua"
)

func (g *Guild) LuaTypeName() string {
	return "LuaGuild"
}

func (g *Guild) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.guildID))
			return 1
		},
		"world_id": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "world_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.worldID))
			return 1
		},
		"name": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "name() is read-only")
				return 0
			}
			L.Push(lua.LString(guild.name))
			return 1
		},
		"leader_id": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "leader_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.leaderCharacterID))
			return 1
		},
		"gp": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "gp() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.gp))
			return 1
		},
		"level": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.Level()))
			return 1
		},
		"message": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			message := L.CheckString(2)
			msgType := constant.MsgPinkText
			if L.GetTop() >= 3 && L.Get(3) != lua.LNil {
				msgType = constant.ServerMessageType(L.CheckInt(3))
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				return 0
			}
			guild.gameWorld.Systems().Get[*GuildSystem]().SendMessageAsync(cfg.ActorContext, guild.guildID, msgType, message)
			return 0
		},
		"capacity": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "capacity() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guild.capacity))
			return 1
		},
		"alliance": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "alliance() takes no arguments")
				return 0
			}
			a := guild.Alliance()
			if a == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, a))
			return 1
		},
		"notice": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "notice() is read-only")
				return 0
			}
			L.Push(lua.LString(guild.notice))
			return 1
		},
		"members": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "members() takes no arguments")
				return 0
			}
			tbl := L.NewTable()
			for i, m := range guild.members {
				if m == nil {
					continue
				}
				tbl.RawSetInt(i+1, luax.NewLuable(L, m))
			}
			L.Push(tbl)
			return 1
		},
		"member": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(3, "member(character) takes exactly one argument")
				return 0
			}
			ch := luax.Check[Character](L, 2)
			m := guild.findMember(ch.GetID())
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, m))
			return 1
		},
		"rank": func(L *lua.LState) int {
			guild := luax.Check[Guild](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(3, "rank(character) takes exactly one argument")
				return 0
			}
			ch := luax.Check[Character](L, 2)
			m := guild.findMember(ch.GetID())
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(guildRankToUint32(m.rank)))
			return 1
		},
		"disband": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 2)
			return luaGuildDisband(L, ch)
		},
		"inc_capacity": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 2)
			extendedCap := L.CheckBool(3)
			return luaGuildIncCapacity(L, ch, extendedCap)
		},
	}
}

func (g *Guild) findMember(characterID uint32) *GuildMember {
	if g == nil {
		return nil
	}
	for _, m := range g.members {
		if m == nil || m.characterID != characterID {
			continue
		}
		return m
	}
	return nil
}

func (g *Guild) String() string {
	return g.LuaTypeName()
}

func (g *Guild) Type() lua.LValueType {
	return lua.LTUserData
}
