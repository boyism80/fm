package entity

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Homing struct {
	SkillWz    *wz.Skill
	SkillLevel uint8
}

type Mob struct {
	LifeCore
	Listener     MobListener
	Wz           *wz.Mob
	Fake         bool
	Foothold     int16
	Spawn        *MobSpawn
	Buffs        *MobBuffContainer
	Skills       *MobSkillContainer
	ExpRate      int32
	DropRate     int32
	stealOutcome *uint32
	Homing       map[uint32]*Homing
	accDamage    map[int64]map[uint32]uint64
	sponge       Sponge
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

func (m *Mob) IsFake() bool {
	return m != nil && m.Fake
}

func (m *Mob) SetFake(fake bool) {
	if m == nil || m.Fake == fake {
		return
	}
	m.Fake = fake
	if !fake {
		return
	}
	mapInstance := m.GetMap()
	if mapInstance == nil || mapInstance.listener == nil {
		return
	}
	mapInstance.listener.OnMobSpawned(mapInstance, m, constant.MobSpawnTypeFake, 0)
}

func (m *Mob) canReceiveMobBuff(flag constant.MobBuffFlag) bool {
	if !m.IsFake() {
		return true
	}
	switch flag {
	case constant.MobBuffStun, constant.MobBuffSpeed, constant.MobBuffPoison, constant.MobBuffVenom:
		return false
	default:
		return true
	}
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

func (m *Mob) Relink(spawnType constant.MobSpawnType, link uint32) bool {
	if m == nil || !m.IsAlive() {
		return false
	}
	mapInstance := m.GetMap()
	if mapInstance == nil || mapInstance.listener == nil {
		return false
	}
	mapInstance.listener.OnMobSpawned(mapInstance, m, spawnType, link)
	return true
}

func (m *Mob) ApplyMobBuff(flag constant.MobBuffFlag, value int32, duration time.Duration, skillWz *wz.Skill, skillLevel uint8, causerOID uint32, stack uint8) {
	if m == nil {
		return
	}
	if !m.canReceiveMobBuff(flag) {
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

func (m *Mob) highestDamageCharacter() *Character {
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return nil
	}

	var winBucketID int64 = soloDamageBucketID
	var winTotal uint64
	for bucketID, damagers := range m.accDamage {
		var total uint64
		for _, dmg := range damagers {
			total += dmg
		}
		if total > winTotal {
			winTotal = total
			winBucketID = bucketID
		}
	}
	if winTotal == 0 {
		return nil
	}

	damagers := m.accDamage[winBucketID]
	var topCID uint32
	var topDamage uint64
	for cid, dmg := range damagers {
		if dmg > topDamage {
			topDamage = dmg
			topCID = cid
		}
	}
	if topCID == 0 {
		return nil
	}
	return mapInstance.GetPlayer(topCID)
}

func (m *Mob) mobDropType(dropChar *Character) constant.DropType {
	if m == nil || m.Wz == nil {
		return constant.DropTypeOwnerOnly
	}
	if m.Wz.ExplosiveReward {
		return constant.DropTypeExplosive
	}
	if m.Wz.FfaLoot {
		return constant.DropTypeFFA
	}
	if dropChar != nil && dropChar.GetPartyID() != nil {
		return constant.DropTypeParty
	}
	return constant.DropTypeOwnerOnly
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

	dropOwner := m.highestDamageCharacter()
	if dropOwner == nil {
		dropOwner = attacker
	}
	var ownerID uint32
	if dropOwner != nil {
		ownerID = dropOwner.GetID()
	}
	dropType := m.mobDropType(dropOwner)

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

			if _, err := mapInstance.SpawnMeso(spawn.count, destPoint, ownerID, dropType, false); err != nil {
				log.Printf("Failed to spawn meso drop: %v", err)
			}
		} else {

			fp := &FieldPlacement{
				ObjectCore: &ObjectCore{
					OID:       0,
					Position:  destPoint,
					GameWorld: m.GameWorld,
				},
				Owner:        ownerID,
				SpawnedPoint: spawnPoint,
				DropType:     dropType,
			}
			fp.ObjectCore.self = fp
			spawn.item.BindFieldPlacement(fp)

			if err := mapInstance.SpawnItem(spawn.item, ownerID, dropType); err != nil {
				log.Printf("Failed to spawn item drop: %v", err)
			}
		}
	}
}

func (m *Mob) ApplyDamage(attacker *Character, amount uint32) bool {
	if amount == 0 || m.GetHp() == 0 {
		return false
	}

	if attacker != nil && attacker.GetInstantKill() {
		amount = m.GetHp()
	}

	damage := amount
	if damage > m.GetHp() {
		damage = m.GetHp()
	}

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

	spongeParent := m.sponge.parent

	m.AddHp(-int(damage))
	killed := m.GetHp() == 0

	m.sponge.applyDamageFromHit(damage)

	if !killed {
		if attacker != nil && spongeParent == nil {
			maxHp := m.GetMaxHp()
			if maxHp == 0 {
				return false
			}
			percent := min(m.GetHp()*100/maxHp, 100)
			attacker.Listener.OnShowMobHp(attacker, m, uint8(percent))
		}
		return false
	}

	return m.Kill(attacker, constant.MobDieAnimationTypeFadeOut)
}

func (m *Mob) runDieScript(attacker *Character) {
	if m == nil || m.Wz == nil {
		return
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	var attackerArg interface{} = lua.LNil
	if attacker != nil {
		attackerArg = attacker
	}

	mobID := m.Wz.ID
	scriptPath := fmt.Sprintf("script/mob/%d.lua", mobID)
	mapInstance.runMobLuaHook(root, scriptPath, fmt.Sprintf("on_mob_die_%d", mobID), m, attackerArg, mapInstance)
	mapInstance.runMobLuaHook(root, constant.CharacterHookScriptPath, "on_mob_die", m, attackerArg, mapInstance)
}

func (m *Mob) onDead(attacker *Character, dieAnim constant.MobDieAnimationType) bool {
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return true
	}

	m.grantKillExp()
	if attacker != nil {
		if m.Wz != nil {
			log.Printf("Mob %d (ID: %d) killed by character %d", m.OID, m.Wz.ID, attacker.GetID())
		}
		m.dropItems(attacker)
	} else if m.Wz != nil {
		log.Printf("Mob %d (ID: %d) killed with no attacker", m.OID, m.Wz.ID)
	}

	pos := m.Position
	revives := []uint32(nil)
	if !m.IsFake() && m.Wz != nil && len(m.Wz.Revives) > 0 {
		revives = m.Wz.Revives
	}

	m.runDieScript(attacker)
	if len(revives) > 0 {
		m.revive(pos, revives)
	}
	m.sponge.onDead(attacker)

	mapInstance.RemoveMob(m.OID, dieAnim)

	if attacker != nil {
		attacker.Listener.OnShowMobHp(attacker, m, 0)
	}
	return true
}

func (m *Mob) Kill(attacker *Character, dieAnim constant.MobDieAnimationType) bool {
	if m == nil {
		return false
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return false
	}
	if mapInstance.GetMob(m.OID) == nil {
		return false
	}
	if hp := m.GetHp(); hp > 0 {
		m.LifeCore.AddHp(-int(hp))
	}
	return m.onDead(attacker, dieAnim)
}

func (m *Mob) removeAfterDieAnimation() constant.MobDieAnimationType {
	if m == nil || m.Wz == nil || m.Wz.SelfDestructionAction < 0 {
		return constant.MobDieAnimationTypeFadeOut
	}
	return constant.MobDieAnimationType(m.Wz.SelfDestructionAction)
}

func (m *Mob) SpawnRevives(reviveIDs []uint32, pos types.Point[int16], spawnType constant.MobSpawnType, link uint32) map[uint32][]*Mob {
	spawned := make(map[uint32][]*Mob)
	if m == nil || len(reviveIDs) == 0 {
		return spawned
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return spawned
	}
	for _, reviveID := range reviveIDs {
		if reviveID == 0 {
			continue
		}
		mob, err := mapInstance.SpawnMob(reviveID, pos, nil, spawnType, link)
		if err != nil {
			continue
		}
		spawned[reviveID] = append(spawned[reviveID], mob)
	}
	return spawned
}

func (m *Mob) runReviveScript(pos types.Point[int16], revives []uint32) bool {
	if m == nil || m.Wz == nil {
		return false
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return false
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return false
	}

	mobID := m.Wz.ID
	scriptPath := fmt.Sprintf("script/mob/%d.lua", mobID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return false
	}

	hook := fmt.Sprintf("on_revive_%d", mobID)
	if thread.GetGlobal(hook).Type() != lua.LTFunction {
		thread.Close()
		return false
	}

	revivesTbl := thread.NewTable()
	for i, id := range revives {
		revivesTbl.RawSetInt(i+1, lua.LNumber(id))
	}
	luax.SetConfiguration(thread, luax.Configuration{
		MapActorPID: mapInstance.GetActorPID(),
	})
	if _, err := luax.Call(thread, hook, m, mapInstance, pos.X, pos.Y, revivesTbl); err != nil {
		log.Printf("mob revive script %s: %v", hook, err)
	}
	return true
}

func (m *Mob) revive(pos types.Point[int16], revives []uint32) {
	if m == nil || m.Wz == nil || len(revives) == 0 {
		return
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return
	}

	if m.runReviveScript(pos, revives) {
		return
	}

	for _, reviveID := range revives {
		if reviveID == 0 {
			continue
		}
		_, err := mapInstance.SpawnMob(reviveID, pos, nil, constant.MobSpawnTypeRevive, m.OID)
		if err != nil {
			log.Printf("Failed to spawn revive mob %d from mob %d: %v", reviveID, m.Wz.ID, err)
		}
	}
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
