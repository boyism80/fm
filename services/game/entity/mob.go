package entity

import (
	"log"
	"math/rand"
	"time"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

type Homing struct {
	SkillWz    *wz.Skill
	SkillLevel uint8
}

type Mob struct {
	LifeCore
	Wz           *wz.Mob
	Foothold     int16
	Spawn        *MobSpawn
	mobBuffs     *MobBuffContainer
	ExpRate      int32
	DropRate     int32
	stealOutcome *uint32
	Homing       map[uint32]*Homing
}

func (m *Mob) SetHoming(causerOID uint32, h *Homing) {
	if m == nil || m.Map == nil {
		return
	}
	if causerOID == 0 {
		return
	}

	if h == nil {
		before, ok := m.Homing[causerOID]
		if !ok || before == nil {
			return
		}
		oldCauser := m.Map.GetPlayer(causerOID)
		if oldCauser != nil && oldCauser.HomingTargetOID != nil && *oldCauser.HomingTargetOID == m.GetOID() {
			oldCauser.HomingTargetOID = nil
		}
		if m.Map.listener != nil {
			m.Map.listener.OnMobHomingRemoved(m.Map, m, before, oldCauser)
		}
		delete(m.Homing, causerOID)
		return
	}

	if h.SkillWz == nil {
		return
	}
	lv := int(h.SkillLevel)
	if lv < 1 {
		return
	}
	if h.SkillWz.MaxLevel > 0 && lv > h.SkillWz.MaxLevel {
		return
	}
	if h.SkillWz.GetLevelData(lv) == nil {
		return
	}

	causer := m.Map.GetPlayer(causerOID)
	if causer != nil && causer.HomingTargetOID != nil {
		prevOID := *causer.HomingTargetOID
		if prevOID != m.GetOID() {
			if prevMob := m.Map.GetMob(prevOID); prevMob != nil {
				prevMob.SetHoming(causerOID, nil)
			}
		}
	}

	before, hadBefore := m.Homing[causerOID]
	if hadBefore && before != nil {
		oldCauser := m.Map.GetPlayer(causerOID)
		if oldCauser != nil && oldCauser.HomingTargetOID != nil && *oldCauser.HomingTargetOID == m.GetOID() {
			oldCauser.HomingTargetOID = nil
		}
		if m.Map.listener != nil {
			m.Map.listener.OnMobHomingRemoved(m.Map, m, before, oldCauser)
		}
	}

	m.Homing[causerOID] = h

	if causer != nil {
		oid := m.GetOID()
		causer.HomingTargetOID = &oid
	}
	if m.Map.listener != nil {
		m.Map.listener.OnMobHomingSet(m.Map, m, h, causer)
	}
}

func (m *Mob) ClearAllHoming() {
	if m == nil || len(m.Homing) == 0 {
		return
	}
	causerOIDs := make([]uint32, 0, len(m.Homing))
	for id := range m.Homing {
		causerOIDs = append(causerOIDs, id)
	}
	for _, id := range causerOIDs {
		m.SetHoming(id, nil)
	}
}

func (m *Mob) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMob
}

func (m *Mob) Is(typ constant.ObjectType) bool {
	return m.GetObjectType().Has(typ)
}

func (m *Mob) SendSpawnSyncToViewer(viewer *Character) {
	if m == nil || viewer == nil {
		return
	}
	viewer.Send(&response.SpawnMob{
		Mob:       m.ToDTO(),
		SpawnType: constant.MOB_SPAWN_TYPE_NONE,
	}, types.SEND_POLICY_ENCRYPT)
}

func (m *Mob) ensureMobBuffs() *MobBuffContainer {
	if m.mobBuffs == nil {
		m.mobBuffs = NewMobBuffContainer(m)
	}
	return m.mobBuffs
}

func (m *Mob) GetMobBuffValue(flag constant.MobBuffFlag) int32 {
	if m.mobBuffs == nil {
		return 0
	}
	return m.mobBuffs.getValue(flag)
}

func (m *Mob) GetMobBuffStack(flag constant.MobBuffFlag) uint8 {
	if m.mobBuffs == nil {
		return 0
	}
	return m.mobBuffs.getStack(flag)
}

func (m *Mob) SetMobBuffStack(flag constant.MobBuffFlag, stack uint8) bool {
	if m.mobBuffs == nil {
		return false
	}
	return m.mobBuffs.setStack(flag, stack)
}

func (m *Mob) ApplyMobBuff(flag constant.MobBuffFlag, value int32, durationMs int64, skillWz *wz.Skill, skillLevel uint8, causerOID uint32, stack uint8) {
	if stack < 1 {
		stack = 1
	}
	now := time.Now()
	m.ensureMobBuffs().AddSkillBuff(now, durationMs, skillWz, skillLevel, causerOID,
		map[constant.MobBuffFlag]int32{flag: value},
		map[constant.MobBuffFlag]uint8{flag: stack})
}

func (m *Mob) CancelMobBuff(flag constant.MobBuffFlag) {
	if m.mobBuffs == nil {
		return
	}
	m.mobBuffs.RemoveBuffForFlag(flag)
}

func (m *Mob) HasBuff(flag constant.MobBuffFlag) bool {
	if m.mobBuffs == nil {
		return false
	}
	return m.mobBuffs.hasFlag(flag)
}

func (m *Mob) GetCauserCharacterID(flag constant.MobBuffFlag) uint32 {
	if m.mobBuffs == nil {
		return 0
	}
	return m.mobBuffs.causerForFlag(flag)
}

func (m *Mob) ClearAllMobBuffTimers() {
	m.mobBuffs = nil
}

func (m *Mob) RemoveExpiredMobBuffs(now time.Time) {
	if m == nil || m.mobBuffs == nil {
		return
	}
	m.mobBuffs.removeExpiredEntities(now)
}

func (m *Mob) getMobBuffMaskAndEntries() (mask uint32, entries []mobBuffForPacket) {
	if m.mobBuffs == nil {
		return 0, nil
	}
	return m.mobBuffs.flattenForSpawnPacket()
}

func (m *Mob) dropItems(attacker *Character) {
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return
	}

	resources := m.Context.GetResources()
	mobDrops, ok := resources.Drops[m.Wz.ID]
	if !ok {
		return
	}

	var drops []struct {
		isMeso bool
		count  int32
		item   Item
	}

	dropRate := float32(m.Context.GetDropRate())
	mesoRate := float32(m.Context.GetMesoRate())
	dropRateMul := int16(100)
	mesoAmountMul := int16(100)
	if attacker.BonusStats.DropRate > 0 {
		dropRateMul = attacker.BonusStats.DropRate
	}
	if attacker.BonusStats.MesoMultiplier > 0 {
		mesoAmountMul = attacker.BonusStats.MesoMultiplier
	}
	dropRateMulF := float32(dropRateMul) / 100.0
	mesoAmountMulF := float32(mesoAmountMul) / 100.0

	mobDrop := m.DropRate
	if mobDrop <= 0 {
		mobDrop = 100
	}
	mobDropF := float32(mobDrop) / 100.0

	for _, drop := range mobDrops {
		if m.stealOutcome != nil && drop.Item != 0 && *m.stealOutcome == drop.Item {
			continue
		}
		adjustedProb := drop.Prob * dropRate * dropRateMulF * mobDropF
		if adjustedProb > 1.0 {
			adjustedProb = 1.0
		}

		if rand.Float32() > adjustedProb {
			continue
		}

		if drop.Item == 0 {
			min := float64(drop.Money) * 0.75
			max := float64(drop.Money)
			count := int32(min + rand.Float64()*(max-min))
			if count == 0 {
				continue
			}
			count = int32(float32(count) * mesoRate)
			if count == 0 {
				continue
			}
			count = int32(float32(count) * mesoAmountMulF)
			if count == 0 {
				continue
			}
			drops = append(drops, struct {
				isMeso bool
				count  int32
				item   Item
			}{isMeso: true, count: count, item: nil})
		} else {

			count := uint16(1)
			if drop.Max != 0 && drop.Min != 0 {
				count = uint16(rand.Intn(int(drop.Max-drop.Min)+1) + int(drop.Min))
			}

			item, err := NewItem(drop.Item, count, m.Context)
			if err != nil {
				continue
			}
			drops = append(drops, struct {
				isMeso bool
				count  int32
				item   Item
			}{isMeso: false, count: int32(count), item: item})
		}
	}

	spawnPoint := m.Position
	spacing := int16(15)

	for i, drop := range drops {
		destPoint := spawnPoint
		if len(drops) > 1 {
			offset := spacing * int16(i/2+1)
			if i%2 == 0 {
				destPoint.X += offset
			} else {
				destPoint.X -= offset
			}
		}

		if drop.isMeso {

			if _, err := mapInstance.SpawnMeso(drop.count, destPoint, attacker.GetID(), constant.DROP_TYPE_OWNED); err != nil {
				log.Printf("Failed to spawn meso drop: %v", err)
			}
		} else {

			d := &Drop{
				ObjectCore: &ObjectCore{
					OID:      0,
					Position: destPoint,
					Context:  m.Context,
				},
				Owner:        attacker.GetID(),
				SpawnedPoint: spawnPoint,
				DropType:     constant.DROP_TYPE_OWNED,
			}
			d.ObjectCore.self = d
			drop.item.BindDrop(d)

			if err := mapInstance.SpawnItem(drop.item, attacker.GetID(), constant.DROP_TYPE_OWNED); err != nil {
				log.Printf("Failed to spawn item drop: %v", err)
			}
		}
	}
}

func (m *Mob) ApplyDamage(attacker *Character, amount uint32) bool {
	if amount == 0 || m.Hp == 0 {
		return false
	}

	damage := amount
	if damage > uint32(m.Hp) {
		damage = uint32(m.Hp)
	}

	m.Hp -= damage
	isDead := m.Hp == 0

	if !isDead {
		if attacker != nil {
			maxHp := m.GetMaxHp()
			if maxHp == 0 {
				return false
			}
			percent := min(uint32(m.Hp)*100/uint32(maxHp), 100)
			attacker.Listener.OnShowMobHp(attacker, m, uint8(percent))
		}
		return false
	}

	if attacker != nil {
		log.Printf("Mob %d (ID: %d) killed by character %d", m.OID, m.Wz.ID, attacker.GetID())

		exp := uint32(m.Wz.EXP)
		if attacker.BonusStats.ExpRate > 0 {
			exp = exp * uint32(attacker.BonusStats.ExpRate) / 100
		}
		mobExp := m.ExpRate
		if mobExp <= 0 {
			mobExp = 100
		}
		exp = exp * uint32(mobExp) / 100
		attacker.AddExp(exp)

		m.dropItems(attacker)
	} else {
		log.Printf("Mob %d (ID: %d) killed with no attacker", m.OID, m.Wz.ID)
	}

	mapInstance := m.GetMap()
	if mapInstance != nil {
		mapInstance.RemoveMob(m.OID, constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT)
	}

	if attacker != nil {
		attacker.Listener.OnShowMobHp(attacker, m, 0)
	}

	return true
}
