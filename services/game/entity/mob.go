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
	Listener     MobListener
	Wz           *wz.Mob
	Foothold     int16
	Spawn        *MobSpawn
	Buffs        *MobBuffContainer
	Skills       *MobSkillContainer
	ExpRate      int32
	DropRate     int32
	stealOutcome *uint32
	Homing       map[uint32]*Homing
	accDamage    map[int64]map[uint32]uint64
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
		SpawnType: constant.MobSpawnTypeNone,
	}, types.SEND_POLICY_ENCRYPT)
}

func (m *Mob) ApplyMobBuff(flag constant.MobBuffFlag, value int32, duration time.Duration, skillWz *wz.Skill, skillLevel uint8, causerOID uint32, stack uint8) {
	if m == nil {
		return
	}
	if stack < 1 {
		stack = 1
	}
	m.Buffs.Add(duration, skillWz, skillLevel, causerOID,
		map[constant.MobBuffFlag]int32{flag: value},
		map[constant.MobBuffFlag]uint8{flag: stack})
}

func (m *Mob) SpawnMist(skill *MobSkill, position types.Point[int16], mistType constant.MistType, bounds types.Rect[int32], duration time.Duration, initialDelay time.Duration, poisonTickMultiplier float64) *Mist {
	if m == nil || m.GameWorld == nil || skill == nil || skill.LevelData == nil || skill.LevelData.SkillWz == nil {
		return nil
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return nil
	}
	b := bounds
	if b.Left == 0 && b.Right == 0 && b.Top == 0 && b.Bottom == 0 {
		pos := m.GetPosition()
		b = skill.LevelData.Bounds.AtOrigin(types.Point[int32]{X: int32(pos.X), Y: int32(pos.Y)})
	}
	if poisonTickMultiplier <= 0 {
		poisonTickMultiplier = 1.0
	}
	causer := uint32(0)
	if m.Wz != nil {
		causer = m.Wz.ID
	}
	mist := &Mist{
		ObjectCore: ObjectCore{
			Position:  position,
			GameWorld: mapInstance.GameWorld,
			Map:       nil,
		},
		Causer:               causer,
		SkillWz:              skill.LevelData.SkillWz,
		MobLevelData:         skill.LevelData,
		SkillLevel:           skill.Slot.Level,
		MistType:             mistType,
		MobMist:              true,
		MobSkill:             true,
		SkillDelay:           0,
		Bounds:               b,
		ExpiresAt:            time.Time{},
		PoisonTickMultiplier: poisonTickMultiplier,
	}
	mist.ObjectCore.self = mist
	if initialDelay > 0 {
		mist.NextPoisonTickAt = time.Now().Add(initialDelay)
	}
	if duration > 0 {
		mist.ExpiresAt = time.Now().Add(duration)
	}
	mapInstance.AddMist(mist)
	return mist
}

func (m *Mob) dropItems(attacker *Character) {
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return
	}

	resources := m.GameWorld.GetResources()
	mobDrops, ok := resources.Drops[m.Wz.ID]
	if !ok {
		return
	}

	var drops []struct {
		isMeso bool
		count  int32
		item   Item
	}

	dropRate := float32(m.GameWorld.GetDropRate())
	mesoRate := float32(m.GameWorld.GetMesoRate())
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

	for _, entry := range mobDrops {
		if m.stealOutcome != nil && entry.Item != 0 && *m.stealOutcome == entry.Item {
			continue
		}
		adjustedProb := entry.Prob * dropRate * dropRateMulF * mobDropF
		if adjustedProb > 1.0 {
			adjustedProb = 1.0
		}

		if rand.Float32() > adjustedProb {
			continue
		}

		if entry.Item == 0 {
			min := float64(entry.Money) * 0.75
			max := float64(entry.Money)
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
			if entry.Max != 0 && entry.Min != 0 {
				count = uint16(rand.Intn(int(entry.Max-entry.Min)+1) + int(entry.Min))
			}

			item, err := NewItem(entry.Item, count, m.GameWorld)
			if err != nil {
				continue
			}
			if eq, ok := item.(Equipment); ok {
				if em, ok := eq.GetModel().(wz.Equipment); ok {
					eq.GetEquipmentCore().RandomizeStats(em)
				}
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

	for i, spawn := range drops {
		destPoint := spawnPoint
		if len(drops) > 1 {
			offset := spacing * int16(i/2+1)
			if i%2 == 0 {
				destPoint.X += offset
			} else {
				destPoint.X -= offset
			}
		}

		if spawn.isMeso {

			if _, err := mapInstance.SpawnMeso(spawn.count, destPoint, attacker.GetID(), constant.DropTypeOwned); err != nil {
				log.Printf("Failed to spawn meso drop: %v", err)
			}
		} else {

			fp := &FieldPlacement{
				ObjectCore: &ObjectCore{
					OID:       0,
					Position:  destPoint,
					GameWorld: m.GameWorld,
				},
				Owner:        attacker.GetID(),
				SpawnedPoint: spawnPoint,
				DropType:     constant.DropTypeOwned,
			}
			fp.ObjectCore.self = fp
			spawn.item.BindFieldPlacement(fp)

			if err := mapInstance.SpawnItem(spawn.item, attacker.GetID(), constant.DropTypeOwned); err != nil {
				log.Printf("Failed to spawn item drop: %v", err)
			}
		}
	}
}

func (m *Mob) ApplyDamage(attacker *Character, amount uint32) bool {
	if amount == 0 || m.GetHp() == 0 {
		return false
	}

	damage := amount
	if damage > m.GetHp() {
		damage = m.GetHp()
	}

	m.AddHp(-int(damage))
	if attacker != nil {
		bucketID := int64(-1)
		if pid := attacker.GetPartyID(); pid != nil {
			bucketID = int64(*pid)
		}
		bucket := m.accDamage[bucketID]
		if bucket == nil {
			bucket = make(map[uint32]uint64)
			m.accDamage[bucketID] = bucket
		}
		bucket[attacker.GetID()] += uint64(damage)
	}
	isDead := m.GetHp() == 0

	if !isDead {
		if attacker != nil {
			maxHp := m.GetMaxHp()
			if maxHp == 0 {
				return false
			}
			percent := min(m.GetHp()*100/maxHp, 100)
			attacker.Listener.OnShowMobHp(attacker, m, uint8(percent))
		}
		return false
	}

	mapInstance := m.GetMap()
	if mapInstance != nil {
		m.distributeKillExperience()
	}

	if attacker != nil {
		log.Printf("Mob %d (ID: %d) killed by character %d", m.OID, m.Wz.ID, attacker.GetID())
		m.dropItems(attacker)
	} else {
		log.Printf("Mob %d (ID: %d) killed with no attacker", m.OID, m.Wz.ID)
	}

	if mapInstance != nil {
		mapInstance.RemoveMob(m.OID, constant.MobDieAnimationTypeFadeOut)
	}

	if attacker != nil {
		attacker.Listener.OnShowMobHp(attacker, m, 0)
	}

	return true
}

func (m *Mob) AddHp(amount int) {
	if m == nil || amount == 0 {
		return
	}

	before := m.GetHp()
	m.LifeCore.AddHp(amount)
	after := m.GetHp()
	if after <= before {
		return
	}
	if m.Listener == nil {
		return
	}

	recovered := after - before
	if recovered > uint32(2147483647) {
		recovered = uint32(2147483647)
	}
	m.Listener.OnMobDamaged(m, -int32(recovered))
}

func (m *Mob) ExpForDamage(damage uint64) uint32 {
	if m == nil || m.Wz == nil || damage == 0 {
		return 0
	}
	maxHp := uint64(m.GetMaxHp())
	if maxHp == 0 {
		return 0
	}
	if damage > maxHp {
		damage = maxHp
	}
	base := uint64(m.Wz.EXP)
	num := base * damage
	den := maxHp
	exp := (num + den - 1) / den
	const maxU32 = uint64(0xffffffff)
	if exp > maxU32 {
		return ^uint32(0)
	}
	return uint32(exp)
}
