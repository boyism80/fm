package entity

import lua "github.com/yuin/gopher-lua"

func (m *GuildMember) LuaTypeName() string {
	return "LuaGuildMember"
}

func (m *GuildMember) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*GuildMember)
			if !ok || mm == nil {
				L.ArgError(1, "GuildMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mm.CharacterID))
			return 1
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*GuildMember)
			if !ok || mm == nil {
				L.ArgError(1, "GuildMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "name() is read-only")
				return 0
			}
			L.Push(lua.LString(mm.CharacterName))
			return 1
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*GuildMember)
			if !ok || mm == nil {
				L.ArgError(1, "GuildMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mm.Level))
			return 1
		},
		"class_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*GuildMember)
			if !ok || mm == nil {
				L.ArgError(1, "GuildMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "class_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mm.ClassID))
			return 1
		},
		"rank": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*GuildMember)
			if !ok || mm == nil {
				L.ArgError(1, "GuildMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "rank() is read-only")
				return 0
			}
			L.Push(lua.LNumber(guildRankToUint32(mm.Rank)))
			return 1
		},
		"channel_index": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*GuildMember)
			if !ok || mm == nil {
				L.ArgError(1, "GuildMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "channel_index() is read-only")
				return 0
			}
			if mm.ChannelIndex == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(*mm.ChannelIndex))
			return 1
		},
	}
}

func (m *GuildMember) String() string {
	return m.LuaTypeName()
}

func (m *GuildMember) Type() lua.LValueType {
	return lua.LTUserData
}
