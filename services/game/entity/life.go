package entity

import (
	"github.com/boyism80/fm/services/game/constant"
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

var (
	_ Life = (*LifeCore)(nil)
	_ Life = (*Character)(nil)
	_ Life = (*Mob)(nil)
)
