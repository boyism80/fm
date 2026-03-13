package entity

import (
	"github.com/boyism80/fm/game/constant"
	lua "github.com/yuin/gopher-lua"
)

type Life struct {
	Object
	Hp         uint16
	BaseHp     uint16
	BonusHp    int16
	Mp         uint16
	BaseMp     uint16
	BonusMp    int16
	Stance     uint8
	Invincible bool
}

type LifeAccessor interface {
	GetHp() uint16
	SetHp(uint16, bool)
	GetMp() uint16
	SetMp(uint16, bool)
	GetMaxHp() uint16
	SetMaxHp(uint16, bool)
	GetMaxMp() uint16
	SetMaxMp(uint16, bool)
	AddHp(int)
	AddMp(int)
	AddHpMp(hpDelta, mpDelta int)
	GetBonusHp() int16
	SetBonusHp(int16)
	GetBonusMp() int16
	SetBonusMp(int16)
	GetInvincible() bool
	SetInvincible(bool)
	IsAlive() bool
}

func (life *Life) GetObject() *Object {
	return &life.Object
}

func (life *Life) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeLife
}

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

func (life *Life) AddBaseHp(amount uint16) {
	newHp := life.BaseHp + amount
	if newHp > constant.STAT_MAX_HP_MP {
		newHp = constant.STAT_MAX_HP_MP
	}
	life.BaseHp = newHp

	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

func (life *Life) AddBonusHp(amount int16) {
	life.BonusHp += amount

	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

func (life *Life) AddBaseMp(amount uint16) {
	newMp := life.BaseMp + amount
	if newMp > constant.STAT_MAX_HP_MP {
		newMp = constant.STAT_MAX_HP_MP
	}
	life.BaseMp = newMp

	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

func (life *Life) AddBonusMp(amount int16) {
	life.BonusMp += amount

	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

func (life *Life) GetHp() uint16       { return life.Hp }
func (life *Life) GetMp() uint16       { return life.Mp }
func (life *Life) GetBonusHp() int16   { return life.BonusHp }
func (life *Life) GetBonusMp() int16   { return life.BonusMp }
func (life *Life) GetInvincible() bool { return life.Invincible }
func (life *Life) IsAlive() bool       { return life.Hp > 0 }

func (life *Life) SetHp(v uint16, _ bool) {
	maxHp := life.GetMaxHp()
	if v > maxHp {
		v = maxHp
	}
	life.Hp = v
}

func (life *Life) SetMp(v uint16, _ bool) {
	maxMp := life.GetMaxMp()
	if v > maxMp {
		v = maxMp
	}
	life.Mp = v
}

func (life *Life) SetMaxHp(v uint16, _ bool) {
	if v > constant.STAT_MAX_HP_MP {
		v = constant.STAT_MAX_HP_MP
	}
	life.BaseHp = v
	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

func (life *Life) SetMaxMp(v uint16, _ bool) {
	if v > constant.STAT_MAX_HP_MP {
		v = constant.STAT_MAX_HP_MP
	}
	life.BaseMp = v
	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

func (life *Life) SetBonusHp(v int16) {
	life.BonusHp = v
	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

func (life *Life) SetBonusMp(v int16) {
	life.BonusMp = v
	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

func (life *Life) SetInvincible(b bool) { life.Invincible = b }

func (life *Life) AddHp(amount int) {
	newHp := int(life.Hp) + amount
	if newHp < 0 {
		newHp = 0
	}
	maxHp := life.GetMaxHp()
	if newHp > int(maxHp) {
		newHp = int(maxHp)
	}
	life.Hp = uint16(newHp)
}

func (life *Life) AddMp(amount int) {
	newMp := int(life.Mp) + amount
	if newMp < 0 {
		newMp = 0
	}
	maxMp := life.GetMaxMp()
	if newMp > int(maxMp) {
		newMp = int(maxMp)
	}
	life.Mp = uint16(newMp)
}

func (life *Life) AddHpMp(hpDelta, mpDelta int) {
	life.AddHp(hpDelta)
	life.AddMp(mpDelta)
}

func (life *Life) LuaTypeName() string {
	return "LuaLife"
}

func (life *Life) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
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
				acc.SetHp(uint16(hp), notify)
				return 0
			default:
				L.ArgError(2, "hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
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
				acc.SetMp(uint16(mp), notify)
				return 0
			default:
				L.ArgError(2, "mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(acc.GetMaxHp()))
				return 1
			case 2, 3:
				hp := L.CheckInt(2)
				if hp < 0 {
					hp = 0
				}
				maxCap := int(constant.STAT_MAX_HP_MP)
				if hp > maxCap {
					hp = maxCap
				}
				notify := argc == 2 || L.ToBool(3)
				acc.SetMaxHp(uint16(hp), notify)
				return 0
			default:
				L.ArgError(2, "max_hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
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
				maxCap := int(constant.STAT_MAX_HP_MP)
				if mp > maxCap {
					mp = maxCap
				}
				notify := argc == 2 || L.ToBool(3)
				acc.SetMaxMp(uint16(mp), notify)
				return 0
			default:
				L.ArgError(2, "max_mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
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
				acc.SetBonusHp(int16(L.CheckInt(2)))
				return 0
			default:
				L.ArgError(2, "bonus_hp() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
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
				acc.SetBonusMp(int16(L.CheckInt(2)))
				return 0
			default:
				L.ArgError(2, "bonus_mp() requires 0 or 1 arguments")
				return 0
			}
		},
		"add_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			acc.AddHp(L.CheckInt(2))
			return 0
		},
		"add_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			acc.AddMp(L.CheckInt(2))
			return 0
		},
		"add_mp_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			acc.AddHpMp(L.CheckInt(2), L.CheckInt(3))
			return 0
		},
		"invincible": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			acc, ok := ud.Value.(LifeAccessor)
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
			acc, ok := ud.Value.(LifeAccessor)
			if !ok {
				L.ArgError(1, "Life expected")
				return 0
			}
			L.Push(lua.LBool(acc.IsAlive()))
			return 1
		},
	}
}

func (life *Life) String() string {
	return life.LuaTypeName()
}

func (life *Life) Type() lua.LValueType {
	return lua.LTUserData
}

var (
	_ LifeAccessor = (*Life)(nil)
	_ LifeAccessor = (*Character)(nil)
	_ LifeAccessor = (*Mob)(nil)
)
