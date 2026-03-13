package entity

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Object struct {
	OID      uint32
	Position types.Vector2[int16]
	Context  GameContext
	Map      *Map
}

func (obj *Object) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeObject
}

type ObjectProvider interface {
	GetObject() *Object
	GetObjectType() constant.ObjectType
	Is(typ constant.ObjectType) bool
}

func (obj *Object) GetObject() *Object {
	return obj
}

func (obj *Object) Is(typ constant.ObjectType) bool {
	return obj.GetObjectType().Has(typ)
}

func (obj *Object) GetMap() *Map {
	return obj.Map
}

func (obj *Object) LuaTypeName() string {
	return "LuaObject"
}

func (obj *Object) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"position": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			provider, ok := ud.Value.(ObjectProvider)
			if !ok {
				L.ArgError(1, "Object expected")
				return 0
			}
			obj := provider.GetObject()

			argc := L.GetTop()
			if argc == 1 {

				L.Push(lua.LNumber(obj.Position.X))
				L.Push(lua.LNumber(obj.Position.Y))
				return 2
			} else if argc == 3 {

				x := L.CheckInt(2)
				y := L.CheckInt(3)
				obj.Position.X = int16(x)
				obj.Position.Y = int16(y)
				return 0
			} else {
				L.ArgError(2, "position() requires 0 or 2 arguments")
				return 0
			}
		},
		"oid": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			provider, ok := ud.Value.(ObjectProvider)
			if !ok {
				L.ArgError(1, "Object expected")
				return 0
			}
			obj := provider.GetObject()

			argc := L.GetTop()
			if argc == 1 {

				L.Push(lua.LNumber(obj.OID))
				return 1
			} else {
				L.ArgError(2, "oid() is read-only")
				return 0
			}
		},
		"map": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			provider, ok := ud.Value.(ObjectProvider)
			if !ok {
				L.Push(lua.LNil)
				return 1
			}
			obj := provider.GetObject()
			mapInstance := obj.GetMap()
			if mapInstance == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mapInstance))
			return 1
		},
		"is": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			provider, ok := ud.Value.(ObjectProvider)
			if !ok {
				L.Push(lua.LFalse)
				return 1
			}
			typeArg := constant.ObjectType(L.CheckInt(2))
			L.Push(lua.LBool(provider.Is(typeArg)))
			return 1
		},
	}
}

func (obj *Object) String() string {
	return obj.LuaTypeName()
}

func (obj *Object) Type() lua.LValueType {
	return lua.LTUserData
}

func (d *Drop) GetObject() *Object {
	if d == nil {
		return nil
	}
	return d.Object
}

func (d *Drop) Is(typ constant.ObjectType) bool {
	return d.GetObjectType().Has(typ)
}

var (
	_ ObjectProvider = (*Object)(nil)
	_ ObjectProvider = (*Life)(nil)
	_ ObjectProvider = (*Character)(nil)
	_ ObjectProvider = (*Mob)(nil)
	_ ObjectProvider = (*Npc)(nil)
	_ ObjectProvider = (*Drop)(nil)
)
