package entity

import (
	"errors"
	"math/rand"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

const mountFatigueTimer = "mount_fatigue"

var (
	ErrCannotRide       = errors.New("cannot ride")
	ErrMountFoodInvalid = errors.New("invalid mount food")
)

type Mount struct {
	owner   *Character
	Level   uint32
	Exp     uint32
	Fatigue uint32
	HP      uint32
}

func (m *Mount) SetLevel(level uint32) {
	level = min(max(level, 1), constant.MountMaxLevel)
	if m.Level == level {
		return
	}

	levelUp := level > m.Level
	m.Level = level
	m.Exp = constant.MountExpByLevel[level-1]
	m.owner.Listener.OnMountUpdated(m.owner, levelUp)
}

func (m *Mount) SetFatigue(fatigue uint32) {
	fatigue = min(fatigue, constant.MountMaxFatigue)
	if m.Fatigue == fatigue {
		return
	}

	m.Fatigue = fatigue
	m.owner.Listener.OnMountUpdated(m.owner, false)
}

func (m *Mount) Ride(skill *SkillEntry, vehicle int32) error {
	if m.owner.GetMap().Wz.Limits(constant.FieldLimitMount) {
		return ErrCannotRide
	}

	m.owner.Buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagPowerguard, constant.BuffFlagManaReflection})
	m.owner.Buffs.AddBuff(skill.Wz, 0, uint8(skill.Level()), m.owner.GetID(), map[constant.BuffFlag]int32{constant.BuffFlagMonsterRiding: vehicle}, true)
	m.start()
	return nil
}

func (m *Mount) riding() *SkillEntry {
	buff, ok := m.owner.Buffs.GetEntity(constant.BuffFlagMonsterRiding).(*SkillBuff)
	if !ok {
		return nil
	}
	return m.owner.Skills.Get(buff.Wz.ID)
}

func (m *Mount) start() {
	skill := m.riding()
	if skill == nil {
		return
	}
	m.owner.CallSkillHook(nil, skill, "on_ride")
}

func (m *Mount) startFatigue() {
	m.owner.AddTimer(mountFatigueTimer, constant.MountFatigueInterval, true, m.tire)
}

func (m *Mount) Dismount() {
	m.owner.RemoveTimer(mountFatigueTimer)
	m.owner.Buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagMonsterRiding})
}

func (m *Mount) changeEquipment(parts constant.EquipmentPartsType) {
	skill := m.riding()
	if skill == nil {
		return
	}
	m.owner.CallSkillHook(nil, skill, "on_rider_equipment_changed", int(parts))
}

func (m *Mount) tire() {
	_, vehicle, riding := m.owner.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	if !riding {
		m.owner.RemoveTimer(mountFatigueTimer)
		return
	}

	fatigue := 1
	if armor, ok := m.owner.GameWorld.GetResources().Items[uint32(vehicle)].(*wz.Armor); ok {
		if tamingMob := m.owner.GameWorld.GetResources().TamingMobs[armor.TamingMob]; tamingMob != nil {
			fatigue = tamingMob.Fatigue
		}
	}
	m.Fatigue = min(m.Fatigue+uint32(fatigue), constant.MountMaxFatigue)
	m.owner.Listener.OnMountUpdated(m.owner, false)
	if m.Fatigue >= constant.MountMaxFatigue {
		m.Dismount()
	}
}

func (m *Mount) Feed(slot int16, itemID uint32) error {
	item := m.owner.Inventory.GetItem(constant.InventoryTypeConsume, slot)
	if item == nil || item.GetModel().GetID() != itemID {
		return ErrMountFoodInvalid
	}
	food, ok := item.GetModel().(*wz.Consume)
	if !ok || food.MountFatigue == 0 {
		return ErrMountFoodInvalid
	}
	if m.owner.GetTimerEntry(mountFatigueTimer) == nil {
		return ErrMountFoodInvalid
	}

	m.owner.Inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
	levelUp := false
	if m.Fatigue > 0 {
		for _, exp := range constant.MountFoodExps {
			if m.Level >= exp.MinLevel {
				m.Exp += uint32(exp.Min + rand.Intn(exp.Max-exp.Min+1))
				break
			}
		}
		for m.Level < constant.MountMaxLevel && m.Exp >= constant.MountExpByLevel[m.Level] {
			m.Level++
			levelUp = true
		}
	}
	m.Fatigue = uint32(max(int(m.Fatigue)+food.MountFatigue, 0))
	m.owner.Listener.OnMountUpdated(m.owner, levelUp)
	return nil
}

func (m *Mount) takeDamage(damage int32) {
	skill := m.riding()
	if skill == nil {
		return
	}
	m.owner.CallSkillHook(nil, skill, "on_rider_damaged", int(damage))
}
