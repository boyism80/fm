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
	owner        *Character
	Level        uint32
	Exp          uint32
	Fatigue      uint32
	BattleshipHP uint32
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

func (m *Mount) Ride(skill *SkillEntry) error {
	var vehicle int32
	switch constant.SkillID(skill.Wz.ID) {
	case constant.SkillMonsterRider, constant.SkillMonsterRiderCygnus:
		tamingMob := m.owner.Inventory.Equipped[constant.EquipmentPartsTamingMob]
		if tamingMob == nil || m.owner.Inventory.Equipped[constant.EquipmentPartsSaddle] == nil {
			return ErrCannotRide
		}
		if m.Fatigue >= constant.MountMaxFatigue {
			return ErrCannotRide
		}
		vehicle = int32(tamingMob.GetModel().GetID())
	case constant.SkillBattleship:
		vehicle = constant.BattleshipVehicle
	default:
		return ErrCannotRide
	}
	if m.owner.GetMap().Wz.Limits(constant.FieldLimitMount) {
		return ErrCannotRide
	}

	m.owner.Buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagPowerguard, constant.BuffFlagManaReflection})
	m.owner.Buffs.AddBuff(skill.Wz, 0, uint8(skill.Level()), m.owner.GetID(), map[constant.BuffFlag]int32{constant.BuffFlagMonsterRiding: vehicle}, true)
	m.start()
	return nil
}

func (m *Mount) start() {
	_, vehicle, riding := m.owner.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	if !riding {
		return
	}

	if vehicle != constant.BattleshipVehicle {
		m.owner.AddTimer(mountFatigueTimer, constant.MountFatigueInterval, true, m.tire)
		return
	}
	skill := m.owner.Skills.Get(uint32(constant.SkillBattleship))
	if skill == nil {
		return
	}
	maxHP := m.maxBattleshipHP(skill)
	if m.BattleshipHP == 0 || m.BattleshipHP > maxHP {
		m.BattleshipHP = maxHP
	}
	if m.BattleshipHP < maxHP {
		m.owner.Listener.OnSkillCooldown(m.owner, constant.BattleshipGauge, uint16(m.BattleshipHP))
	}
}

func (m *Mount) maxBattleshipHP(skill *SkillEntry) uint32 {
	return uint32(max(200*(int(m.owner.GetLevel())+2*skill.Level()-120), 1))
}

func (m *Mount) dismount() {
	m.owner.RemoveTimer(mountFatigueTimer)
	m.owner.Buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagMonsterRiding})
}

func (m *Mount) dismountRider(parts constant.EquipmentPartsType) {
	if parts != constant.EquipmentPartsTamingMob && parts != constant.EquipmentPartsSaddle {
		return
	}
	_, vehicle, riding := m.owner.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	if !riding || vehicle == constant.BattleshipVehicle {
		return
	}
	m.dismount()
}

func (m *Mount) tire() {
	_, vehicle, riding := m.owner.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	if !riding || vehicle == constant.BattleshipVehicle {
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
		m.dismount()
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
	_, vehicle, riding := m.owner.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	if !riding || vehicle == constant.BattleshipVehicle {
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

func (m *Mount) hitBattleship(damage int32) {
	_, vehicle, riding := m.owner.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	if !riding || vehicle != constant.BattleshipVehicle {
		return
	}

	m.BattleshipHP = uint32(max(int(m.BattleshipHP)-int(damage), 0))
	m.owner.Listener.OnSkillCooldown(m.owner, constant.BattleshipGauge, uint16(m.BattleshipHP))
	if m.BattleshipHP > 0 {
		return
	}
	m.dismount()
	if skill := m.owner.Skills.Get(uint32(constant.SkillBattleship)); skill != nil {
		if level := skill.Wz.GetLevelData(skill.Level()); level != nil {
			skill.StartCooldown(level.Cooldown)
		}
	}
}
