package wz

import lua "github.com/yuin/gopher-lua"

func (*Consume) LuaTypeName() string { return luaWzConsumeTypeName }
func (*Consume) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
		"consume_type": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushConsumeType(L, ud.Value)
			return 1
		},
		"bonus_stats": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			consume, ok := ud.Value.(*Consume)
			if !ok || consume == nil {
				L.Push(lua.LNil)
				return 1
			}
			t := L.NewTable()
			t.RawSetString("str", lua.LNumber(consume.ScrollIncStr))
			t.RawSetString("dex", lua.LNumber(consume.ScrollIncDex))
			t.RawSetString("int", lua.LNumber(consume.ScrollIncInt))
			t.RawSetString("luk", lua.LNumber(consume.ScrollIncLuk))
			t.RawSetString("max_hp", lua.LNumber(consume.ScrollIncMaxHP))
			t.RawSetString("max_mp", lua.LNumber(consume.ScrollIncMaxMP))
			t.RawSetString("pad", lua.LNumber(consume.ScrollIncPAD))
			t.RawSetString("mad", lua.LNumber(consume.ScrollIncMAD))
			t.RawSetString("pdd", lua.LNumber(consume.ScrollIncPDD))
			t.RawSetString("mdd", lua.LNumber(consume.ScrollIncMDD))
			t.RawSetString("acc", lua.LNumber(consume.ScrollIncACC))
			t.RawSetString("avoid", lua.LNumber(consume.ScrollIncAvoid))
			t.RawSetString("hands", lua.LNumber(consume.ScrollIncHands))
			t.RawSetString("speed", lua.LNumber(consume.ScrollIncSpeed))
			t.RawSetString("jump", lua.LNumber(consume.ScrollIncJump))
			L.Push(t)
			return 1
		},
		"rand_stat": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			consume, ok := ud.Value.(*Consume)
			if !ok || consume == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(consume.ScrollRandStat))
			return 1
		},
		"recover": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			consume, ok := ud.Value.(*Consume)
			if !ok || consume == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(consume.ScrollRecover))
			return 1
		},
	}
}
func (c *Consume) String() string       { return c.LuaTypeName() }
func (c *Consume) Type() lua.LValueType { return lua.LTUserData }
