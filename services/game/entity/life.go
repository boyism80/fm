package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type LifeCore struct {
	ObjectCore
	hp         uint32
	mp         uint32
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
	SetBonusHp(int32, bool)
	GetBonusMp() int32
	SetBonusMp(int32, bool)
	GetInvincible() bool
	SetInvincible(bool)
	IsAlive() bool
	GetStance() uint8
	SetStance(uint8)
}

func (life *LifeCore) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeLife
}

func (life *LifeCore) Is(typ constant.ObjectType) bool {
	return life.GetObjectType().Has(typ)
}

func (life *LifeCore) SendSpawnSyncToViewer(viewer *Character) {}

func (life *LifeCore) SendDestroySyncToViewer(viewer *Character) {}

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
	if life.GetHp() > life.GetMaxHp() {
		life.setHp(life.GetMaxHp())
	}
}

func (life *LifeCore) AddBaseMp(amount uint32) {
	life.BaseMp += amount
}

func (life *LifeCore) AddBonusMp(amount int32) {
	life.BonusMp += amount
	if life.GetMp() > life.GetMaxMp() {
		life.setMp(life.GetMaxMp())
	}
}

func (life *LifeCore) GetHp() uint32       { return life.hp }
func (life *LifeCore) GetMp() uint32       { return life.mp }
func (life *LifeCore) GetBonusHp() int32   { return life.BonusHp }
func (life *LifeCore) GetBonusMp() int32   { return life.BonusMp }
func (life *LifeCore) GetInvincible() bool { return life.Invincible }
func (life *LifeCore) IsAlive() bool       { return life.hp > 0 }

func (life *LifeCore) GetStance() uint8 {
	return life.Stance
}

func (life *LifeCore) SetStance(v uint8) {
	life.Stance = v
}

func (life *LifeCore) setHp(v uint32) {
	maxHp := life.GetMaxHp()
	if v > maxHp {
		v = maxHp
	}
	life.hp = v
}

func (life *LifeCore) setMp(v uint32) {
	maxMp := life.GetMaxMp()
	if v > maxMp {
		v = maxMp
	}
	life.mp = v
}

func (life *LifeCore) SetHp(v uint32, _ bool) {
	life.setHp(v)
}

func (life *LifeCore) SetMp(v uint32, _ bool) {
	life.setMp(v)
}

func (life *LifeCore) SetBaseHp(v uint32, _ bool) {
	life.BaseHp = v
	if life.GetHp() > life.GetMaxHp() {
		life.setHp(life.GetMaxHp())
	}
}

func (life *LifeCore) SetBaseMp(v uint32, _ bool) {
	life.BaseMp = v
	if life.GetMp() > life.GetMaxMp() {
		life.setMp(life.GetMaxMp())
	}
}

func (life *LifeCore) SetBonusHp(v int32, _ bool) {
	life.BonusHp = v
	if life.GetHp() > life.GetMaxHp() {
		life.setHp(life.GetMaxHp())
	}
}

func (life *LifeCore) SetBonusMp(v int32, _ bool) {
	life.BonusMp = v
	if life.GetMp() > life.GetMaxMp() {
		life.setMp(life.GetMaxMp())
	}
}

func (life *LifeCore) SetInvincible(b bool) { life.Invincible = b }

func (life *LifeCore) AddHp(amount int) {
	n := int(life.GetHp()) + amount
	if n < 0 {
		n = 0
	}
	m := int(life.GetMaxHp())
	if n > m {
		n = m
	}
	life.setHp(uint32(n))
}

func (life *LifeCore) AddMp(amount int) {
	n := int(life.GetMp()) + amount
	if n < 0 {
		n = 0
	}
	m := int(life.GetMaxMp())
	if n > m {
		n = m
	}
	life.setMp(uint32(n))
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
