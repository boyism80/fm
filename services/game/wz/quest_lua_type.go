package wz

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

const luaWzQuestTypeName = "LuaWzQuest"

func (*Quest) LuaTypeName() string { return luaWzQuestTypeName }

func (*Quest) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			q, ok := ud.Value.(*Quest)
			if !ok || q == nil {
				L.ArgError(1, "WzQuest expected")
				return 0
			}
			L.Push(lua.LNumber(q.ID))
			return 1
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			q, ok := ud.Value.(*Quest)
			if !ok || q == nil {
				L.ArgError(1, "WzQuest expected")
				return 0
			}
			L.Push(lua.LString(q.Meta.Name))
			return 1
		},
		"complete_requirements": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			q, ok := ud.Value.(*Quest)
			if !ok || q == nil {
				L.ArgError(1, "WzQuest expected")
				return 0
			}
			L.Push(questRequirementsToLuaTable(L, q.Complete.Requirements))
			return 1
		},
		"start_requirements": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			q, ok := ud.Value.(*Quest)
			if !ok || q == nil {
				L.ArgError(1, "WzQuest expected")
				return 0
			}
			L.Push(questRequirementsToLuaTable(L, q.Start.Requirements))
			return 1
		},
		"party_ranks": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			q, ok := ud.Value.(*Quest)
			if !ok || q == nil {
				L.ArgError(1, "WzQuest expected")
				return 0
			}
			if len(q.PartyRanks) == 0 {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(partyRanksToLuaTable(L, q.PartyRanks))
			return 1
		},
	}
}

func (q *Quest) String() string       { return q.LuaTypeName() }
func (q *Quest) Type() lua.LValueType { return lua.LTUserData }

var _ luax.Luable = (*Quest)(nil)
