package game

import (
	"errors"
	"math/rand"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

const mountFatigueTimer = "mount_fatigue"

var (
	ErrCannotRide       = errors.New("cannot ride")
	ErrMountFoodInvalid = errors.New("invalid mount food")
)

type Mount struct {
	owner   *Character
	level   uint32
	exp     uint32
	fatigue uint32
	hp      uint32
}

func (m *Mount) SetLevel(level uint32) {
	level = min(max(level, 1), constant.MountMaxLevel)
	if m.level == level {
		return
	}

	levelUp := level > m.level
	m.level = level
	m.exp = constant.MountExpByLevel[level-1]
	m.owner.listener.OnMountUpdated(m.owner, m.level, m.exp, m.fatigue, levelUp)
}

func (m *Mount) SetFatigue(fatigue uint32) {
	fatigue = min(fatigue, constant.MountMaxFatigue)
	if m.fatigue == fatigue {
		return
	}

	m.fatigue = fatigue
	m.owner.listener.OnMountUpdated(m.owner, m.level, m.exp, m.fatigue, false)
}

func (m *Mount) Ride(skill *SkillEntry, vehicle int32) error {
	if field := m.owner.GetMap(); field == nil || field.wz.Limits(constant.FieldLimitMount) {
		return ErrCannotRide
	}

	m.owner.buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagPowerguard, constant.BuffFlagManaReflection})
	m.owner.buffs.AddBuff(skill.wz, 0, uint8(skill.Level()), m.owner.GetID(), map[constant.BuffFlag]int32{constant.BuffFlagMonsterRiding: vehicle}, true)
	m.start()
	return nil
}

func (m *Mount) riding() *SkillEntry {
	buff, ok := m.owner.buffs.GetEntity(constant.BuffFlagMonsterRiding).(*SkillBuff)
	if !ok {
		return nil
	}
	return m.owner.skills.Get(buff.wz.ID)
}

func (m *Mount) start() {
	skill := m.riding()
	if skill == nil {
		return
	}
	m.owner.skills.CallHook(nil, skill, "on_ride")
}

func (m *Mount) startFatigue() {
	m.owner.AddTimer(mountFatigueTimer, constant.MountFatigueInterval, true, m.tire)
}

func (m *Mount) Dismount() {
	m.owner.RemoveTimer(mountFatigueTimer)
	m.owner.buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagMonsterRiding})
}

func (m *Mount) changeEquipment(parts constant.EquipmentPartsType) {
	skill := m.riding()
	if skill == nil {
		return
	}
	m.owner.skills.CallHook(nil, skill, "on_rider_equipment_changed", int(parts))
}

func (m *Mount) tire() {
	_, vehicle, riding := m.owner.buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	if !riding {
		m.owner.RemoveTimer(mountFatigueTimer)
		return
	}

	fatigue := 1
	if armor, ok := m.owner.gameWorld.GetResources().Items[uint32(vehicle)].(*wz.Armor); ok {
		if tamingMob := m.owner.gameWorld.GetResources().TamingMobs[armor.TamingMob]; tamingMob != nil {
			fatigue = tamingMob.Fatigue
		}
	}
	m.fatigue = min(m.fatigue+uint32(fatigue), constant.MountMaxFatigue)
	m.owner.listener.OnMountUpdated(m.owner, m.level, m.exp, m.fatigue, false)
	if m.fatigue >= constant.MountMaxFatigue {
		m.Dismount()
	}
}

func (m *Mount) Feed(slot int16, itemID uint32) error {
	item := m.owner.inventory.GetItem(constant.InventoryTypeConsume, slot)
	if item == nil || item.GetModel().GetID() != itemID {
		return ErrMountFoodInvalid
	}
	food, ok := item.GetModel().(*wz.Consume)
	if !ok || food.MountFatigue == 0 {
		return ErrMountFoodInvalid
	}
	if m.owner.inventory.equipped[constant.EquipmentPartsTamingMob] == nil || m.owner.inventory.equipped[constant.EquipmentPartsSaddle] == nil {
		return ErrMountFoodInvalid
	}

	m.owner.inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
	levelUp := false
	if m.fatigue > 0 {
		for _, exp := range constant.MountFoodExps {
			if m.level >= exp.MinLevel {
				m.exp += uint32(exp.Min + rand.Intn(exp.Max-exp.Min+1))
				break
			}
		}
		for m.level < constant.MountMaxLevel && m.exp >= constant.MountExpByLevel[m.level] {
			m.level++
			levelUp = true
		}
	}
	m.fatigue = uint32(max(int(m.fatigue)+food.MountFatigue, 0))
	m.owner.listener.OnMountUpdated(m.owner, m.level, m.exp, m.fatigue, levelUp)
	return nil
}

func (m *Mount) takeDamage(damage int32) {
	skill := m.riding()
	if skill == nil {
		return
	}
	m.owner.skills.CallHook(nil, skill, "on_rider_damaged", int(damage))
}
