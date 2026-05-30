package entity

import lua "github.com/yuin/gopher-lua"

func (life *LifeCore) LuaTypeName() string {
	return "LuaLife"
}

func (life *LifeCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(acc.GetHp()))
				return 1
			case 2, 3:
				hp := L.CheckInt(2)
				if hp < 0 {
					hp = 0
				}
				maxHp := acc.GetMaxHp()
				if hp > int(maxHp) {
					hp = int(maxHp)
				}
				notify := argc == 2 || L.ToBool(3)
				acc.SetHp(uint32(hp), notify)
				return 0
			default:
				L.ArgError(2, "hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(acc.GetMp()))
				return 1
			case 2, 3:
				mp := L.CheckInt(2)
				if mp < 0 {
					mp = 0
				}
				maxMp := acc.GetMaxMp()
				if mp > int(maxMp) {
					mp = int(maxMp)
				}
				notify := argc == 2 || L.ToBool(3)
				acc.SetMp(uint32(mp), notify)
				return 0
			default:
				L.ArgError(2, "mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(life.GetMaxHp()))
				return 1
			case 2, 3:
				hp := L.CheckInt(2)
				if hp < 0 {
					hp = 0
				}
				notify := argc == 2 || L.ToBool(3)
				life.SetBaseHp(uint32(hp), notify)
				return 0
			default:
				L.ArgError(2, "max_hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(acc.GetMaxMp()))
				return 1
			case 2, 3:
				mp := L.CheckInt(2)
				if mp < 0 {
					mp = 0
				}
				notify := argc == 2 || L.ToBool(3)
				acc.SetBaseMp(uint32(mp), notify)
				return 0
			default:
				L.ArgError(2, "max_mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(acc.GetBonusHp()))
				return 1
			case 2:
				v := int32(L.CheckInt(2))
				acc.SetBonusHp(v, true)
				return 0
			default:
				L.ArgError(2, "bonus_hp() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(acc.GetBonusMp()))
				return 1
			case 2:
				v := int32(L.CheckInt(2))
				acc.SetBonusMp(v, true)
				return 0
			default:
				L.ArgError(2, "bonus_mp() requires 0 or 1 arguments")
				return 0
			}
		},
		"add_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			acc.AddHp(L.CheckInt(2))
			return 0
		},
		"add_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			acc.AddMp(L.CheckInt(2))
			return 0
		},
		"add_mp_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			acc.AddHpMp(L.CheckInt(2), L.CheckInt(3))
			return 0
		},
		"invincible": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(acc.GetInvincible()))
				return 1
			case 2:
				acc.SetInvincible(L.CheckBool(2))
				return 0
			default:
				L.ArgError(2, "invincible() requires 0 or 1 arguments")
				return 0
			}
		},
		"is_alive": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			L.Push(lua.LBool(acc.IsAlive()))
			return 1
		},
		"stance": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(acc.GetStance()))
				return 1
			case 2:
				acc.SetStance(uint8(L.CheckInt(2)))
				return 0
			default:
				L.ArgError(2, "stance() requires 0 or 1 argument")
				return 0
			}
		},
	}
}

func (life *LifeCore) String() string {
	return life.LuaTypeName()
}

func (life *LifeCore) Type() lua.LValueType {
	return lua.LTUserData
}
