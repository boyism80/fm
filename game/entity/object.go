package entity

import (
	"github.com/boyism80/fm/common/types"
	lua "github.com/yuin/gopher-lua"
)

var objectBuiltinFuncs = map[string]lua.LGFunction{}

type Object struct {
	Position types.Vector2[int16]
}

func (o *Object) LuaTypeName() string {
	return "Object"
}
func (o *Object) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return objectBuiltinFuncs
}
