package entity

import lua "github.com/yuin/gopher-lua"

var lifeBuiltinFuncs = map[string]lua.LGFunction{}

type Life struct {
	Object
	Hp      uint16
	MaxHp   uint16
	Mp      uint16
	MaxMp   uint16
	IsAlive bool
}

func (l *Life) LuaTypeName() string {
	return "Life"
}
func (l *Life) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return lifeBuiltinFuncs
}
