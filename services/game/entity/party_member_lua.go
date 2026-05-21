package entity

import lua "github.com/yuin/gopher-lua"

func (m *PartyMember) LuaTypeName() string {
	return "LuaPartyMember"
}

func (m *PartyMember) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*PartyMember)
			if !ok || mm == nil {
				L.ArgError(1, "PartyMember expected")
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
			mm, ok := ud.Value.(*PartyMember)
			if !ok || mm == nil {
				L.ArgError(1, "PartyMember expected")
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
			mm, ok := ud.Value.(*PartyMember)
			if !ok || mm == nil {
				L.ArgError(1, "PartyMember expected")
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
			mm, ok := ud.Value.(*PartyMember)
			if !ok || mm == nil {
				L.ArgError(1, "PartyMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "class_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mm.ClassID))
			return 1
		},
		"role": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*PartyMember)
			if !ok || mm == nil {
				L.ArgError(1, "PartyMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "role() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mm.Role))
			return 1
		},
		"map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*PartyMember)
			if !ok || mm == nil {
				L.ArgError(1, "PartyMember expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "map_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mm.MapID))
			return 1
		},
		"channel_index": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mm, ok := ud.Value.(*PartyMember)
			if !ok || mm == nil {
				L.ArgError(1, "PartyMember expected")
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

func (m *PartyMember) String() string {
	return m.LuaTypeName()
}

func (m *PartyMember) Type() lua.LValueType {
	return lua.LTUserData
}
