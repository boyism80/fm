package entity

import (
	"github.com/boyism80/fm/core/types"
	lua "github.com/yuin/gopher-lua"
)

type Object struct {
	OID      uint32
	Position types.Vector2[int16]
	Context  GameContext // GameContext for accessing game resources
}

// Luable interface implementation
func (obj *Object) LuaTypeName() string {
	return "LuaObject"
}

func (obj *Object) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"position": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			obj, ok := ud.Value.(*Object)
			if !ok {
				L.ArgError(1, "Object expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return x, y
				L.Push(lua.LNumber(obj.Position.X))
				L.Push(lua.LNumber(obj.Position.Y))
				return 2
			} else if argc == 3 {
				// Setter: position(x, y)
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
			obj, ok := ud.Value.(*Object)
			if !ok {
				L.ArgError(1, "Object expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return oid
				L.Push(lua.LNumber(obj.OID))
				return 1
			} else {
				L.ArgError(2, "oid() is read-only")
				return 0
			}
		},
	}
}

func (obj *Object) String() string {
	return obj.LuaTypeName()
}

func (obj *Object) Type() lua.LValueType {
	return lua.LTUserData
}
