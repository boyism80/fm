package entity

import lua "github.com/yuin/gopher-lua"

type Life struct {
	Object
	Hp         uint16
	MaxHp      uint16
	Mp         uint16
	MaxMp      uint16
	Stance     uint8
	Invincible bool
}

// Luable interface implementation
func (life *Life) LuaTypeName() string {
	return "LuaLife"
}

func (life *Life) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return hp
				L.Push(lua.LNumber(life.Hp))
				return 1
			} else if argc == 2 {
				// Setter: hp(value)
				hp := L.CheckInt(2)
				if hp < 0 {
					hp = 0
				}
				if hp > int(life.MaxHp) {
					hp = int(life.MaxHp)
				}
				life.Hp = uint16(hp)
				return 0
			} else {
				L.ArgError(2, "hp() requires 0 or 1 arguments")
				return 0
			}
		},
		"max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return max_hp
				L.Push(lua.LNumber(life.MaxHp))
				return 1
			} else if argc == 2 {
				// Setter: max_hp(value)
				maxHp := L.CheckInt(2)
				if maxHp < 0 {
					maxHp = 0
				}
				life.MaxHp = uint16(maxHp)
				// HP가 새로운 MaxHP를 초과하면 조정
				if life.Hp > life.MaxHp {
					life.Hp = life.MaxHp
				}
				return 0
			} else {
				L.ArgError(2, "max_hp() requires 0 or 1 arguments")
				return 0
			}
		},
		"mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return mp
				L.Push(lua.LNumber(life.Mp))
				return 1
			} else if argc == 2 {
				// Setter: mp(value)
				mp := L.CheckInt(2)
				if mp < 0 {
					mp = 0
				}
				if mp > int(life.MaxMp) {
					mp = int(life.MaxMp)
				}
				life.Mp = uint16(mp)
				return 0
			} else {
				L.ArgError(2, "mp() requires 0 or 1 arguments")
				return 0
			}
		},
		"max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return max_mp
				L.Push(lua.LNumber(life.MaxMp))
				return 1
			} else if argc == 2 {
				// Setter: max_mp(value)
				maxMp := L.CheckInt(2)
				if maxMp < 0 {
					maxMp = 0
				}
				life.MaxMp = uint16(maxMp)
				// MP가 새로운 MaxMP를 초과하면 조정
				if life.Mp > life.MaxMp {
					life.Mp = life.MaxMp
				}
				return 0
			} else {
				L.ArgError(2, "max_mp() requires 0 or 1 arguments")
				return 0
			}
		},
		"invincible": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return invincible
				L.Push(lua.LBool(life.Invincible))
				return 1
			} else if argc == 2 {
				// Setter: invincible(value)
				invincible := L.CheckBool(2)
				life.Invincible = invincible
				return 0
			} else {
				L.ArgError(2, "invincible() requires 0 or 1 arguments")
				return 0
			}
		},
	}
}

func (life *Life) String() string {
	return life.LuaTypeName()
}

func (life *Life) Type() lua.LValueType {
	return lua.LTUserData
}
