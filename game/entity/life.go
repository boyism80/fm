package entity

import (
	"github.com/boyism80/fm/game/constant"
	lua "github.com/yuin/gopher-lua"
)

type LifeCore struct {
	ObjectCore
	Hp         uint32
	Mp         uint32
	BaseHp     uint32
	BaseMp     uint32
	BonusHp    int32
	BonusMp    int32
	Stance     uint8
	Invincible bool
}

type Life interface {
	GetHp() uint32
	SetHp(uint32, bool)
	GetMp() uint32
	SetMp(uint32, bool)
	GetMaxHp() uint32
	SetBaseHp(uint32, bool)
	GetMaxMp() uint32
	SetBaseMp(uint32, bool)
	AddHp(int)
	AddMp(int)
	AddHpMp(hpDelta, mpDelta int)
	GetBonusHp() int32
	SetBonusHp(int32)
	GetBonusMp() int32
	SetBonusMp(int32)
	GetInvincible() bool
	SetInvincible(bool)
	IsAlive() bool
}

func (life *LifeCore) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeLife
}

func (life *LifeCore) Is(typ constant.ObjectType) bool {
	return life.GetObjectType().Has(typ)
}

func (life *LifeCore) SendSpawnSyncToViewer(viewer *Character) {}

func (life *LifeCore) GetMaxHp() uint32 {
	t := int64(life.BaseHp) + int64(life.BonusHp)
	if t < 1 {
		return 1
	}
	if t > 0xffffffff {
		return 0xffffffff
	}
	return uint32(t)
}

func (life *LifeCore) GetMaxMp() uint32 {
	t := int64(life.BaseMp) + int64(life.BonusMp)
	if t < 0 {
		return 0
	}
	if t > 0xffffffff {
		return 0xffffffff
	}
	return uint32(t)
}

func (life *LifeCore) AddBaseHp(amount uint32) {
	life.BaseHp += amount
}

func (life *LifeCore) AddBonusHp(amount int32) {
	life.BonusHp += amount
	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

func (life *LifeCore) AddBaseMp(amount uint32) {
	life.BaseMp += amount
}

func (life *LifeCore) AddBonusMp(amount int32) {
	life.BonusMp += amount
	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

func (life *LifeCore) GetHp() uint32       { return life.Hp }
func (life *LifeCore) GetMp() uint32       { return life.Mp }
func (life *LifeCore) GetBonusHp() int32   { return life.BonusHp }
func (life *LifeCore) GetBonusMp() int32   { return life.BonusMp }
func (life *LifeCore) GetInvincible() bool { return life.Invincible }
func (life *LifeCore) IsAlive() bool       { return life.Hp > 0 }

func (life *LifeCore) SetHp(v uint32, _ bool) {
	maxHp := life.GetMaxHp()
	if v > maxHp {
		v = maxHp
	}
	life.Hp = v
}

func (life *LifeCore) SetMp(v uint32, _ bool) {
	maxMp := life.GetMaxMp()
	if v > maxMp {
		v = maxMp
	}
	life.Mp = v
}

func (life *LifeCore) SetBaseHp(v uint32, _ bool) {
	life.BaseHp = v
	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

func (life *LifeCore) SetBaseMp(v uint32, _ bool) {
	life.BaseMp = v
	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

func (life *LifeCore) SetBonusHp(v int32) {
	life.BonusHp = v
	if life.Hp > life.GetMaxHp() {
		life.Hp = life.GetMaxHp()
	}
}

func (life *LifeCore) SetBonusMp(v int32) {
	life.BonusMp = v
	if life.Mp > life.GetMaxMp() {
		life.Mp = life.GetMaxMp()
	}
}

func (life *LifeCore) SetInvincible(b bool) { life.Invincible = b }

func (life *LifeCore) AddHp(amount int) {
	n := int(life.Hp) + amount
	if n < 0 {
		n = 0
	}
	m := int(life.GetMaxHp())
	if n > m {
		n = m
	}
	life.Hp = uint32(n)
}

func (life *LifeCore) AddMp(amount int) {
	n := int(life.Mp) + amount
	if n < 0 {
		n = 0
	}
	m := int(life.GetMaxMp())
	if n > m {
		n = m
	}
	life.Mp = uint32(n)
}

func (life *LifeCore) AddHpMp(hpDelta, mpDelta int) {
	life.AddHp(hpDelta)
	life.AddMp(mpDelta)
}

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
				acc.SetBonusHp(v)
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
				acc.SetBonusMp(v)
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
	}
}

func (life *LifeCore) String() string {
	return life.LuaTypeName()
}

func (life *LifeCore) Type() lua.LValueType {
	return lua.LTUserData
}

var (
	_ Life = (*LifeCore)(nil)
	_ Life = (*Character)(nil)
	_ Life = (*Mob)(nil)
)
