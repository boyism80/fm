package entity

import (
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
	}
}

func (r *Reactor) String() string {
	return r.LuaTypeName()
}

func (r *Reactor) Type() lua.LValueType {
	return lua.LTUserData
}
