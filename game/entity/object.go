package entity

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type ObjectCore struct {
	OID      uint32
	Position types.Vector2[int16]
	Context  GameContext
	Map      *Map
}

func (obj *ObjectCore) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeObject
}

type Object interface {
	GetOID() uint32
	GetPosition() types.Vector2[int16]
	SetPosition(x, y int16)
	GetContext() GameContext
	GetMap() *Map
	GetObjectType() constant.ObjectType
	Is(typ constant.ObjectType) bool
}

func (obj *ObjectCore) GetOID() uint32 {
	return obj.OID
}

func (obj *ObjectCore) GetPosition() types.Vector2[int16] {
	return obj.Position
}

func (obj *ObjectCore) SetPosition(x, y int16) {
	obj.Position.X = x
	obj.Position.Y = y
}

func (obj *ObjectCore) GetContext() GameContext {
	return obj.Context
}

func (obj *ObjectCore) Is(typ constant.ObjectType) bool {
	return obj.GetObjectType().Has(typ)
}

func (obj *ObjectCore) GetMap() *Map {
	return obj.Map
}

func (obj *ObjectCore) LuaTypeName() string {
	return "LuaObject"
}

func (obj *ObjectCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"position": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			obj, ok := ud.Value.(Object)
			if !ok {
				L.ArgError(1, "Object expected")
				return 0
			}
			pos := obj.GetPosition()

			argc := L.GetTop()
			if argc == 1 {

				L.Push(lua.LNumber(pos.X))
				L.Push(lua.LNumber(pos.Y))
				return 2
			} else if argc == 3 {

				x := L.CheckInt(2)
				y := L.CheckInt(3)
				obj.SetPosition(int16(x), int16(y))
				return 0
			} else {
				L.ArgError(2, "position() requires 0 or 2 arguments")
				return 0
			}
		},
		"oid": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			obj, ok := ud.Value.(Object)
			if !ok {
				L.ArgError(1, "Object expected")
				return 0
			}
			oid := obj.GetOID()

			argc := L.GetTop()
			if argc == 1 {

				L.Push(lua.LNumber(oid))
				return 1
			} else {
				L.ArgError(2, "oid() is read-only")
				return 0
			}
		},
		"map": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			obj, ok := ud.Value.(Object)
			if !ok {
				L.Push(lua.LNil)
				return 1
			}
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
			obj, ok := ud.Value.(Object)
			if !ok {
				L.Push(lua.LFalse)
				return 1
			}
			typeArg := constant.ObjectType(L.CheckInt(2))
			L.Push(lua.LBool(obj.Is(typeArg)))
			return 1
		},
	}
}

func (obj *ObjectCore) String() string {
	return obj.LuaTypeName()
}

func (obj *ObjectCore) Type() lua.LValueType {
	return lua.LTUserData
}

func (d *Drop) Is(typ constant.ObjectType) bool {
	return d.GetObjectType().Has(typ)
}

var (
	_ Object = (*ObjectCore)(nil)
	_ Object = (*LifeCore)(nil)
	_ Object = (*Character)(nil)
	_ Object = (*Mob)(nil)
	_ Object = (*Npc)(nil)
	_ Object = (*Drop)(nil)
)
