package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (p *Party) LuaTypeName() string {
	return "LuaParty"
}

func (p *Party) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			party, ok := ud.Value.(*Party)
			if !ok || party == nil {
				L.ArgError(1, "Party expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(party.PartyID))
			return 1
		},
		"world_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			party, ok := ud.Value.(*Party)
			if !ok || party == nil {
				L.ArgError(1, "Party expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "world_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(party.WorldID))
			return 1
		},
		"leader_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			party, ok := ud.Value.(*Party)
			if !ok || party == nil {
				L.ArgError(1, "Party expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "leader_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(party.LeaderCharacterID))
			return 1
		},
		"revision": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			party, ok := ud.Value.(*Party)
			if !ok || party == nil {
				L.ArgError(1, "Party expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "revision() is read-only")
				return 0
			}
			L.Push(lua.LNumber(party.Revision))
			return 1
		},
		"state": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			party, ok := ud.Value.(*Party)
			if !ok || party == nil {
				L.ArgError(1, "Party expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "state() is read-only")
				return 0
			}
			L.Push(lua.LNumber(party.State))
			return 1
		},
		"members": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			party, ok := ud.Value.(*Party)
			if !ok || party == nil {
				L.ArgError(1, "Party expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "members() takes no arguments")
				return 0
			}
			tbl := L.NewTable()
			for i, m := range party.Members {
				if m == nil {
					continue
				}
				tbl.RawSetInt(i+1, luax.NewLuable(L, m))
			}
			L.Push(tbl)
			return 1
		},
	}
}

func (p *Party) String() string {
	return p.LuaTypeName()
}

func (p *Party) Type() lua.LValueType {
	return lua.LTUserData
}
