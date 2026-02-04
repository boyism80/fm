package entity

import (
	"github.com/boyism80/fm/game/constant"
	lua "github.com/yuin/gopher-lua"
)

type Life struct {
	Object
	Hp         uint16
	BaseHp     uint16 // Base maximum HP (from level, AP, or WZ data)
	BonusHp    int16  // Bonus maximum HP (from buffs/equipment, can be negative)
	Mp         uint16
	BaseMp     uint16 // Base maximum MP (from level, AP, or WZ data)
	BonusMp    int16  // Bonus maximum MP (from buffs/equipment, can be negative)
	Stance     uint8
	Invincible bool
}

// GetMaxHp returns the current maximum HP (BaseHp + BonusHp)
func (life *Life) GetMaxHp() uint16 {
	total := int32(life.BaseHp) + int32(life.BonusHp)
	if total < 1 {
		return 1
	}
	if total > int32(constant.STAT_MAX_HP_MP) {
		return constant.STAT_MAX_HP_MP
	}
	return uint16(total)
}

// GetMaxMp returns the current maximum MP (BaseMp + BonusMp)
func (life *Life) GetMaxMp() uint16 {
	total := int32(life.BaseMp) + int32(life.BonusMp)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_HP_MP) {
		return constant.STAT_MAX_HP_MP
	}
	return uint16(total)
}

// AddBaseHp adds to base maximum HP (permanent change)
func (life *Life) AddBaseHp(amount uint16) {
	newHp := life.BaseHp + amount
	if newHp > constant.STAT_MAX_HP_MP {
		newHp = constant.STAT_MAX_HP_MP
	}
	life.BaseHp = newHp
	// Adjust current HP if it exceeds new max
	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

// AddBonusHp adds to bonus maximum HP (temporary change from buffs/equipment)
func (life *Life) AddBonusHp(amount int16) {
	life.BonusHp += amount
	// Adjust current HP if it exceeds new max
	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

// AddBaseMp adds to base maximum MP (permanent change)
func (life *Life) AddBaseMp(amount uint16) {
	newMp := life.BaseMp + amount
	if newMp > constant.STAT_MAX_HP_MP {
		newMp = constant.STAT_MAX_HP_MP
	}
	life.BaseMp = newMp
	// Adjust current MP if it exceeds new max
	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

// AddBonusMp adds to bonus maximum MP (temporary change from buffs/equipment)
func (life *Life) AddBonusMp(amount int16) {
	life.BonusMp += amount
	// Adjust current MP if it exceeds new max
	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
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
				maxHp := life.GetMaxHp()
				if hp > int(maxHp) {
					hp = int(maxHp)
				}
				life.Hp = uint16(hp)
				return 0
			} else {
				L.ArgError(2, "hp() requires 0 or 1 arguments")
				return 0
			}
		},
		"get_max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			L.Push(lua.LNumber(life.GetMaxHp()))
			return 1
		},
		"set_base_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			baseHp := L.CheckInt(2)
			if baseHp < 0 {
				baseHp = 0
			}
			life.BaseHp = uint16(baseHp)
			// Adjust HP if it exceeds new max
			if life.Hp > life.GetMaxHp() {
				life.Hp = life.GetMaxHp()
			}
			return 0
		},
		"add_bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			bonusHp := L.CheckInt(2)
			life.AddBonusHp(int16(bonusHp))
			return 0
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
				maxMp := life.GetMaxMp()
				if mp > int(maxMp) {
					mp = int(maxMp)
				}
				life.Mp = uint16(mp)
				return 0
			} else {
				L.ArgError(2, "mp() requires 0 or 1 arguments")
				return 0
			}
		},
		"get_max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			L.Push(lua.LNumber(life.GetMaxMp()))
			return 1
		},
		"set_base_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			baseMp := L.CheckInt(2)
			if baseMp < 0 {
				baseMp = 0
			}
			life.BaseMp = uint16(baseMp)
			// Adjust MP if it exceeds new max
			if life.Mp > life.GetMaxMp() {
				life.Mp = life.GetMaxMp()
			}
			return 0
		},
		"add_bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			bonusMp := L.CheckInt(2)
			life.AddBonusMp(int16(bonusMp))
			return 0
		},
		"get_bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			L.Push(lua.LNumber(life.BonusHp))
			return 1
		},
		"set_bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			bonusHp := L.CheckInt(2)
			life.BonusHp = int16(bonusHp)
			// Adjust HP if it exceeds new max
			if life.Hp > life.GetMaxHp() {
				life.Hp = life.GetMaxHp()
			}
			return 0
		},
		"get_bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			L.Push(lua.LNumber(life.BonusMp))
			return 1
		},
		"set_bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			life, ok := ud.Value.(*Life)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			bonusMp := L.CheckInt(2)
			life.BonusMp = int16(bonusMp)
			// Adjust MP if it exceeds new max
			if life.Mp > life.GetMaxMp() {
				life.Mp = life.GetMaxMp()
			}
			return 0
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
