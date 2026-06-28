package entity

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

func (r *Reactor) LuaTypeName() string {
	return "LuaReactor"
}

func (r *Reactor) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			reactor, ok := ud.Value.(*Reactor)
			if !ok || reactor == nil {
				L.ArgError(1, "Reactor expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			if reactor.Wz == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(reactor.Wz.ID))
			return 1
		},
		"state": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			reactor, ok := ud.Value.(*Reactor)
			if !ok || reactor == nil {
				L.ArgError(1, "Reactor expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "state() is read-only")
				return 0
			}
			L.Push(lua.LNumber(reactor.State))
			return 1
		},
		"react_item_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			reactor, ok := ud.Value.(*Reactor)
			if !ok || reactor == nil {
				L.ArgError(1, "Reactor expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "react_item_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(reactor.ReactItemID()))
			return 1
		},
		"react_item_quantity": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			reactor, ok := ud.Value.(*Reactor)
			if !ok || reactor == nil {
				L.ArgError(1, "Reactor expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "react_item_quantity() is read-only")
				return 0
			}
			L.Push(lua.LNumber(reactor.ReactItemQuantity()))
			return 1
		},
		"hit": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			reactor, ok := ud.Value.(*Reactor)
			if !ok || reactor == nil {
				L.ArgError(1, "Reactor expected")
				return 0
			}
			var trigger *Character
			if L.GetTop() >= 2 {
				arg := L.Get(2)
				if state, ok := arg.(lua.LNumber); ok {
					reactor.ForceHitState(byte(state))
					return 0
				}
				if arg != lua.LNil {
					triggerUd, ok := arg.(*lua.LUserData)
					if !ok {
						L.ArgError(2, "Character or state expected")
						return 0
					}
					trigger, _ = triggerUd.Value.(*Character)
				}
			}
			reactor.Hit(trigger, constant.ReactorHitAirLeft, 0)
			return 0
		},
		"trigger": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			reactor, ok := ud.Value.(*Reactor)
			if !ok || reactor == nil {
				L.ArgError(1, "Reactor expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "trigger() is read-only")
				return 0
			}
			trigger := reactor.GetTrigger()
			if trigger == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, trigger))
			return 1
		},
		"drop_items": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			reactor, ok := ud.Value.(*Reactor)
			if !ok || reactor == nil {
				L.ArgError(1, "Reactor expected")
				return 0
			}
			reactor.DropItems()
			return 0
		},
	}
}

func (r *Reactor) String() string {
	return r.LuaTypeName()
}

func (r *Reactor) Type() lua.LValueType {
	return lua.LTUserData
}
