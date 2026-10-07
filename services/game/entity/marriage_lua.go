package entity

import lua "github.com/yuin/gopher-lua"

func (m *Marriage) LuaTypeName() string {
	return "LuaMarriage"
}

func (m *Marriage) LuaBuiltinFuncs() map[string]lua.LGFunction {
	check := func(L *lua.LState) *Marriage {
		ud := L.CheckUserData(1)
		marriage, ok := ud.Value.(*Marriage)
		if !ok || marriage == nil {
			L.ArgError(1, "Marriage expected")
			return nil
		}
		return marriage
	}
	wishes := func(L *lua.LState, values []string) *lua.LTable {
		tbl := L.NewTable()
		for _, value := range values {
			tbl.Append(lua.LString(value))
		}
		return tbl
	}
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).ID))
			return 1
		},
		"status": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).Status))
			return 1
		},
		"married": func(L *lua.LState) int {
			L.Push(lua.LBool(check(L).Status == MarriageStatusMarried))
			return 1
		},
		"groom_id": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).GroomID))
			return 1
		},
		"bride_id": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).BrideID))
			return 1
		},
		"groom_name": func(L *lua.LState) int {
			L.Push(lua.LString(check(L).GroomName))
			return 1
		},
		"bride_name": func(L *lua.LState) int {
			L.Push(lua.LString(check(L).BrideName))
			return 1
		},
		"partner_id": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).PartnerID(uint32(L.CheckNumber(2)))))
			return 1
		},
		"ticket": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).TicketItemID))
			return 1
		},
		"wished": func(L *lua.LState) int {
			L.Push(lua.LBool(check(L).Wished(uint32(L.CheckNumber(2)))))
			return 1
		},
		"reserved": func(L *lua.LState) int {
			L.Push(lua.LBool(check(L).Reserved()))
			return 1
		},
		"groom_wishes": func(L *lua.LState) int {
			L.Push(wishes(L, check(L).GroomWishes))
			return 1
		},
		"bride_wishes": func(L *lua.LState) int {
			L.Push(wishes(L, check(L).BrideWishes))
			return 1
		},
		"divorce_requested_at": func(L *lua.LState) int {
			marriage := check(L)
			if marriage.DivorceRequestedAt.IsZero() {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(marriage.DivorceRequestedAt.Unix()))
			return 1
		},
	}
}

func (m *Marriage) String() string {
	return m.LuaTypeName()
}

func (m *Marriage) Type() lua.LValueType {
	return lua.LTUserData
}
