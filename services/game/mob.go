package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Homing struct {
	skillWz    *wz.Skill
	skillLevel uint8
}

type Mob struct {
	lifeCore
	listener      MobListener
	wz            *wz.Mob
	fake          bool
	foothold      int16
	spawn         *MobSpawn
	buffs         *MobBuffs
	skills        *MobSkills
	expRate       int32
	dropRate      int32
	stealOutcome  *uint32
	homing        map[uint32]*Homing
	accDamage     map[int64]map[uint32]uint64
	sponge        Sponge
	lastHitAt     time.Time
	dropItemCount int
	deathStarted  bool
	carnivalTeam  constant.CarnivalTeam
}

func (m *Mob) SetHoming(causerOID uint32, h *Homing) {
	if m == nil || m.onMap == nil {
		return
	}
	if causerOID == 0 {
		return
	}

	if h == nil {
		before, ok := m.homing[causerOID]
		if !ok || before == nil {
			return
		}
		oldCauser := m.onMap.GetPlayer(causerOID)
		if oldCauser != nil && oldCauser.skills.homingTargetOID != nil && *oldCauser.skills.homingTargetOID == m.GetOID() {
			oldCauser.skills.homingTargetOID = nil
		}
		if m.onMap.listener != nil {
			m.onMap.listener.OnMobHomingRemoved(m.onMap, m, before, oldCauser)
		}
		delete(m.homing, causerOID)
		return
	}

	if h.skillWz == nil {
		return
	}
	lv := int(h.skillLevel)
	if lv < 1 {
		return
	}
	if h.skillWz.MaxLevel > 0 && lv > h.skillWz.MaxLevel {
		return
	}
	if h.skillWz.GetLevelData(lv) == nil {
		return
	}

	causer := m.onMap.GetPlayer(causerOID)
	if causer != nil && causer.skills.homingTargetOID != nil {
		prevOID := *causer.skills.homingTargetOID
		if prevOID != m.GetOID() {
			if prevMob := m.onMap.GetMob(prevOID); prevMob != nil {
				prevMob.SetHoming(causerOID, nil)
			}
		}
	}

	before, hadBefore := m.homing[causerOID]
	if hadBefore && before != nil {
		oldCauser := m.onMap.GetPlayer(causerOID)
		if oldCauser != nil && oldCauser.skills.homingTargetOID != nil && *oldCauser.skills.homingTargetOID == m.GetOID() {
			oldCauser.skills.homingTargetOID = nil
		}
		if m.onMap.listener != nil {
			m.onMap.listener.OnMobHomingRemoved(m.onMap, m, before, oldCauser)
		}
	}

	m.homing[causerOID] = h

	if causer != nil {
		oid := m.GetOID()
		causer.skills.homingTargetOID = &oid
	}
	if m.onMap.listener != nil {
		m.onMap.listener.OnMobHomingSet(m.onMap, m, h, causer)
	}
}

func (m *Mob) ClearAllHoming() {
	if m == nil || len(m.homing) == 0 {
		return
	}
	causerOIDs := make([]uint32, 0, len(m.homing))
	for id := range m.homing {
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

func (m *Mob) SetScaleLevel(scale int) {
	if m == nil || m.wz == nil || scale <= 0 {
		return
	}
	base := int(m.wz.Level)
	if base <= 0 || scale == base {
		return
	}
	ratio := float64(scale) / float64(base)
	hp := uint32(math.Max(1, math.Ceil(float64(m.baseHp)*ratio)))
	mp := uint32(math.Ceil(float64(m.baseMp) * ratio))
	m.SetBaseHp(hp, false)
	m.SetHp(hp, false)
	m.SetBaseMp(mp, false)
	m.SetMp(mp, false)
}

func (m *Mob) SetFake(fake bool) {
	if m == nil || m.fake == fake {
		return
	}
	m.fake = fake
	if !fake {
		return
	}
	mapInstance := m.GetMap()
	if mapInstance == nil || mapInstance.listener == nil {
		return
	}
	mapInstance.listener.OnMobSpawned(mapInstance, m, constant.MobSpawnTypeFake, 0)
}

func (m *Mob) canReceiveMobBuff(flag constant.MobBuffFlag, skillWz *wz.Skill) bool {
	if skillWz != nil && skillWz.MobSkill {
		return true
	}
	if m.fake {
		switch flag {
		case constant.MobBuffStun, constant.MobBuffSpeed, constant.MobBuffPoison, constant.MobBuffVenom:
			return false
		}
	}
	if m.wz == nil {
		return true
	}
	if skillWz != nil && m.resistsElement(skillWz) {
		return false
	}
	if flag == constant.MobBuffDoom && m.wz.NoDoom {
		return false
	}
	if m.wz.Boss {
		switch flag {
		case constant.MobBuffSpeed, constant.MobBuffWatk, constant.MobBuffBlind:
			return true
		default:
			return false
		}
	}
	return true
}

func (m *Mob) resistsElement(skillWz *wz.Skill) bool {
	elements := make([]string, 0, 2)
	if skillWz.ElemAttr != "" {
		elements = append(elements, strings.ToLower(skillWz.ElemAttr[:1]))
	}
	if skillWz.ExtraElement != "" {
		elements = append(elements, skillWz.ExtraElement)
	}
	for _, element := range elements {
		switch m.wz.ElemResist[element] {
		case 1, 2:
			return true
		}
	}
	return false
}

func (m *Mob) SendSpawnSyncToViewer(viewer *Character) {
	if m == nil || viewer == nil {
		return
	}
	viewer.Send(&response.SpawnMob{
		Mob:       m.BuildS2C(),
		SpawnType: constant.MobSpawnTypeNone,
	}, types.SEND_POLICY_ENCRYPT)
}

func (m *Mob) SendDestroySyncToViewer(viewer *Character) {
	if m == nil || viewer == nil {
		return
	}
	viewer.Send(&response.DieMob{
		OID:           m.oid,
		AnimationType: constant.MobDieAnimationTypeDisappear,
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

func (m *Mob) SpawnMist(skill *MobSkill, spec mistSpec) *Mist {
	if skill.levelData == nil || skill.levelData.SkillWz == nil {
		return nil
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return nil
	}
	if spec.bounds == (types.Rect[int32]{}) {
		spec.bounds = skill.levelData.Bounds.AtOrigin(types.Point[int32]{X: int32(m.position.X), Y: int32(m.position.Y)})
	}
	causer := uint32(0)
	if m.wz != nil {
		causer = m.wz.ID
	}

	mist := &Mist{
		objectCore:   ObjectCore{position: m.position},
		causer:       causer,
		skillWz:      skill.levelData.SkillWz,
		mobLevelData: skill.levelData,
		skillLevel:   skill.slot.Level,
		mobMist:      true,
	}
	mist.spread(mapInstance, spec)
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
	if m == nil || m.wz == nil {
		return constant.DropTypeOwnerOnly
	}
	if m.wz.ExplosiveReward {
		return constant.DropTypeExplosive
	}
	if m.wz.FfaLoot {
		return constant.DropTypeFFA
	}
	if dropChar != nil && dropChar.party.ID() != nil {
		return constant.DropTypeParty
	}
	return constant.DropTypeOwnerOnly
}

func (m *Mob) dropItems(attacker *Character) {
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return
	}

	resources := m.gameWorld.GetResources()
	mobDrops := resources.MobDrops[m.wz.ID]
	if len(mobDrops) == 0 && len(resources.GlobalDrops) == 0 {
		return
	}

	type rolledDrop struct {
		isMeso      bool
		count       int32
		item        Item
		quest       uint32
		questTagged bool
	}
	type dropRule struct {
		wz.Drop
		quest  uint32
		global bool
	}
	var drops []rolledDrop
	var entries []dropRule
	for _, entry := range mobDrops {
		entries = append(entries, dropRule{Drop: entry})
	}
	for _, entry := range resources.GlobalDrops {
		if !entry.IsAvailableAt(mapInstance.TemplateID()) || !entry.IsAvailableFor(m.wz) {
			continue
		}
		if entry.Quest != 0 && !attacker.quests.NeedsItem(entry.Quest, entry.Item) {
			continue
		}
		entries = append(entries, dropRule{Drop: entry.Drop, quest: entry.Quest, global: true})
	}

	dropRate := float32(m.gameWorld.GetDropRate())
	mesoRate := float32(m.gameWorld.GetMesoRate())
	dropRateMul := max(int16(attacker.family.DropRate()), attacker.stats.bonus.dropRate)
	mesoAmountMul := int16(100)
	if attacker.stats.bonus.mesoMultiplier > 0 {
		mesoAmountMul = attacker.stats.bonus.mesoMultiplier
	}
	if _, v, ok := attacker.buffs.GetBuffValue(constant.BuffFlagMesoUp); ok {
		mesoAmountMul += int16(v - 100)
	}
	dropRateMulF := float32(dropRateMul) / 100.0
	mesoAmountMulF := float32(mesoAmountMul) / 100.0

	mobDrop := m.dropRate
	if mobDrop <= 0 {
		mobDrop = 100
	}
	mobDrop += m.buffs.GetValue(constant.MobBuffShowdown)
	mobDropF := float32(mobDrop) / 100.0

	for _, entry := range entries {
		if !entry.global && m.stealOutcome != nil && entry.Item != 0 && *m.stealOutcome == entry.Item {
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
			drops = append(drops, rolledDrop{isMeso: true, count: count})
		} else {

			count := uint16(1)
			if entry.Max != 0 && entry.Min != 0 {
				count = uint16(rand.Intn(int(entry.Max-entry.Min)+1) + int(entry.Min))
			}

			item, err := NewItem(entry.Item, count, m.gameWorld)
			if err != nil {
				continue
			}
			if eq, ok := item.(Equipment); ok {
				if em, ok := eq.GetModel().(wz.Equipment); ok {
					eq.GetEquipmentCore().RandomizeStats(em)
				}
			}
			drops = append(drops, rolledDrop{count: int32(count), item: item, quest: entry.quest, questTagged: !entry.global})
		}
	}

	spawnPoint := m.position

	dropOwner := m.highestDamageCharacter()
	if dropOwner == nil {
		dropOwner = attacker
	}
	var ownerID uint32
	if dropOwner != nil {
		ownerID = dropOwner.GetID()
	}
	dropType := m.mobDropType(dropOwner)
	spacing := int16(25)
	if dropType == constant.DropTypeExplosive {
		spacing = 40
	}

	d := int16(1)
	for _, spawn := range drops {
		destPoint := spawnPoint
		if d%2 == 0 {
			destPoint.X += spacing * (d + 1) / 2
		} else {
			destPoint.X -= spacing * (d / 2)
		}

		if spawn.isMeso {
			if _, err := mapInstance.SpawnMeso(spawn.count, destPoint, spawnPoint, ownerID, dropType, false); err != nil {
				log.Printf("Failed to spawn meso drop: %v", err)
			}
		} else {
			if err := mapInstance.SpawnItem(spawn.item, ItemSpawn{
				Position:    destPoint,
				From:        spawnPoint,
				Owner:       ownerID,
				DropType:    dropType,
				QuestTagged: spawn.questTagged,
				Quest:       spawn.quest,
			}); err != nil {
				log.Printf("Failed to spawn item drop: %v", err)
			}
		}
		d++
	}
}

func (m *Mob) TakeDamage(attacker *Character, amount uint32) bool {
	if amount == 0 || m.GetHp() == 0 {
		return false
	}
	if attacker != nil && m.carnivalTeam != constant.CarnivalTeamNone {
		if team := attacker.carnival.Team(); team != nil && team.teamID == m.carnivalTeam {
			return false
		}
	}

	if attacker != nil && attacker.gm.instantKill {
		amount = m.GetHp()
	}

	damage := amount
	if damage > m.GetHp() {
		damage = m.GetHp()
	}

	if attacker != nil {
		bucketID := int64(-1)
		if pid := attacker.party.ID(); pid != nil {
			bucketID = int64(*pid)
		}
		bucket := m.accDamage[bucketID]
		if bucket == nil {
			bucket = make(map[uint32]uint64)
			m.accDamage[bucketID] = bucket
		}
		bucket[attacker.GetID()] += uint64(damage)
	}

	m.MarkHit()

	spongeParent := m.sponge.parent

	m.AddHp(-int(damage))
	killed := m.GetHp() == 0

	m.sponge.damage(damage)

	if !killed {
		if attacker != nil && spongeParent == nil {
			if m.wz == nil || !m.wz.DamagedByMob {
				maxHp := m.GetMaxHp()
				if maxHp == 0 {
					return false
				}
				percent := min(m.GetHp()*100/maxHp, 100)
				m.showHp(attacker, uint8(percent))
			}
		}
		return false
	}

	return m.Kill(attacker, constant.MobDieAnimationTypeFadeOut)
}

func (m *Mob) showHp(attacker *Character, percent uint8) {
	mapInstance := m.GetMap()
	if mapInstance == nil || mapInstance.wz.FieldType != wz.FieldTypeDojo {
		attacker.listener.OnShowMobHp(attacker, m, percent)
		return
	}
	for _, player := range mapInstance.GetAllPlayers() {
		ch, ok := player.(*Character)
		if !ok {
			continue
		}
		ch.listener.OnShowMobHp(ch, m, percent)
	}
}

func (m *Mob) runDieScript(attacker *Character) {
	if m == nil || m.wz == nil {
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

	mobID := m.wz.ID
	scriptPath := fmt.Sprintf("script/mob/%d.lua", mobID)
	mapInstance.runMobLuaHook(root, scriptPath, "on_mob_die", m, attackerArg, mapInstance)
}

func (m *Mob) onDead(attacker *Character, dieAnim constant.MobDieAnimationType, reward bool) bool {
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return true
	}

	if reward {
		m.grantKillExp()
	}
	if reward && attacker != nil {
		if m.wz != nil {
			attacker.quests.RecordMobKill(m.wz.ID)
			if match := attacker.carnival.Match(); match != nil && m.wz != nil {
				cpAmount := m.wz.CP
				if cpAmount <= 0 {
					cpAmount = m.wz.Point
				}
				if cpAmount > 0 {
					match.OnMobKilled(attacker, cpAmount)
				}
			}
		}
		m.dropItems(attacker)
	}

	pos := m.position
	revives := []uint32(nil)
	if !m.fake && m.wz != nil && len(m.wz.Revives) > 0 {
		revives = m.wz.Revives
	}

	m.runDieScript(attacker)
	if len(revives) > 0 {
		m.revive(pos, revives)
	}
	m.sponge.onDead(attacker)

	sm := mapInstance.StateMachine()
	mapInstance.RemoveMob(m.oid, dieAnim)
	if sm != nil {
		sm.CallHook("on_mob_die", m)
	}

	if attacker != nil {
		m.showHp(attacker, 0)
	}
	if sm := mapInstance.StateMachine(); sm != nil && !mapInstance.HasAliveMobs() {
		sm.CallHook("on_all_monsters_dead")
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
	if mapInstance.GetMob(m.oid) == nil {
		return false
	}
	if m.deathStarted {
		return false
	}
	m.deathStarted = true
	if hp := m.GetHp(); hp > 0 {
		m.lifeCore.AddHp(-int(hp))
	}
	return m.onDead(attacker, dieAnim, true)
}

func (m *Mob) SelfDestruct(attacker *Character) bool {
	if !m.wz.HasSelfDestruction || attacker.Spectating() || attacker.IsHidden() || !attacker.IsAlive() {
		return false
	}
	mapInstance := m.GetMap()
	if mapInstance == nil || mapInstance.GetMob(m.oid) == nil || m.deathStarted || !m.IsAlive() {
		return false
	}
	m.deathStarted = true
	m.lifeCore.AddHp(-int(m.GetHp()))
	return m.onDead(attacker, constant.MobDieAnimationType(m.wz.SelfDestructionAction), false)
}

// Mob bodies have no hitbox in WZ, so a body attack reaches only targets this close.
const mobTouchRange int64 = 100

func (m *Mob) canAttack(target *Mob, damage uint32) bool {
	if !m.IsAlive() || !target.IsAlive() {
		return false
	}
	sourcePosition := m.GetPosition()
	targetPosition := target.GetPosition()
	dx := int64(targetPosition.X) - int64(sourcePosition.X)
	dy := int64(targetPosition.Y) - int64(sourcePosition.Y)
	for _, attack := range m.wz.Attacks {
		maxDamage := attack.PADamage
		if attack.Magic {
			maxDamage = attack.MADamage
		}
		if maxDamage <= 0 {
			if attack.Magic {
				maxDamage = m.wz.MADamage
			} else {
				maxDamage = m.wz.PADamage
			}
		}
		if attack.Magic {
			maxDamage += int(m.buffs.GetValue(constant.MobBuffMatk))
			maxDamage += int(m.buffs.GetValue(constant.MobBuffMagicAttackUp))
			if target.wz.Invincible || target.buffs.Has(constant.MobBuffMagicImmunity) || target.buffs.Has(constant.MobBuffDamageImmunity) {
				continue
			}
		} else {
			maxDamage += int(m.buffs.GetValue(constant.MobBuffWatk))
			maxDamage += int(m.buffs.GetValue(constant.MobBuffWeaponAttackUp))
			if target.wz.Invincible || target.buffs.Has(constant.MobBuffWeaponImmunity) || target.buffs.Has(constant.MobBuffDamageImmunity) {
				continue
			}
		}
		if maxDamage <= 0 || uint64(damage) > uint64(maxDamage) {
			continue
		}
		if attack.RangeR > 0 && dx*dx+dy*dy <= int64(attack.RangeR)*int64(attack.RangeR) {
			return true
		}
		if attack.HasRange {
			x := int32(dx)
			y := int32(dy)
			left := min(attack.RangeBounds.Left, -attack.RangeBounds.Right)
			right := max(attack.RangeBounds.Right, -attack.RangeBounds.Left)
			if x >= left && x <= right && y >= attack.RangeBounds.Top && y <= attack.RangeBounds.Bottom {
				return true
			}
		}
	}

	if m.wz.BodyAttack == 0 || dx*dx+dy*dy > mobTouchRange*mobTouchRange {
		return false
	}
	if target.wz.Invincible || target.buffs.Has(constant.MobBuffWeaponImmunity) || target.buffs.Has(constant.MobBuffDamageImmunity) {
		return false
	}
	maxDamage := m.wz.PADamage + int(m.buffs.GetValue(constant.MobBuffWatk)) + int(m.buffs.GetValue(constant.MobBuffWeaponAttackUp))
	return maxDamage > 0 && uint64(damage) <= uint64(maxDamage)
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
		mob, err := mapInstance.SpawnMob(reviveID, pos, nil, spawnType, link, m.carnivalTeam)
		if err != nil {
			continue
		}
		spawned[reviveID] = append(spawned[reviveID], mob)
	}
	return spawned
}

func (m *Mob) runReviveScript(pos types.Point[int16], revives []uint32) bool {
	if m == nil || m.wz == nil {
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

	mobID := m.wz.ID
	scriptPath := fmt.Sprintf("script/mob/%d.lua", mobID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return false
	}

	hook := "on_revive"
	if !luax.HasFunc(thread, hook) {
		luax.Close(thread)
		return false
	}

	revivesTbl := thread.NewTable()
	for i, id := range revives {
		revivesTbl.RawSetInt(i+1, lua.LNumber(id))
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: mapInstance.LogicActorPID(),
	})
	luax.CallAsync(nil, root, thread, hook, m, mapInstance, pos.X, pos.Y, revivesTbl).OnError(func(err error) {
		log.Printf("mob revive script %s: %v", hook, err)
	})
	return true
}

func (m *Mob) revive(pos types.Point[int16], revives []uint32) {
	if m == nil || m.wz == nil || len(revives) == 0 {
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
		_, err := mapInstance.SpawnMob(reviveID, pos, nil, constant.MobSpawnTypeRevive, m.oid, m.carnivalTeam)
		if err != nil {
			log.Printf("Failed to spawn revive mob %d from mob %d: %v", reviveID, m.wz.ID, err)
		}
	}
}

func (m *Mob) SetHp(v uint32, notify bool) {
	if m == nil {
		return
	}
	m.setHp(v)
	if !notify || m.onMap == nil {
		return
	}
	maxHp := m.GetMaxHp()
	percent := uint8(0)
	if maxHp > 0 {
		percent = uint8(uint64(m.GetHp()) * 100 / uint64(maxHp))
	}
	packet := &response.ShowMobHp{OID: m.oid, Percentage: percent}
	for _, player := range m.onMap.GetAllPlayers() {
		ch, ok := player.(*Character)
		if !ok {
			continue
		}
		ch.Send(packet, types.SEND_POLICY_ENCRYPT)
	}
}

func (m *Mob) AddHp(amount int) {
	if m == nil || amount == 0 {
		return
	}

	before := m.GetHp()
	m.lifeCore.AddHp(amount)
	after := m.GetHp()
	if after <= before {
		return
	}

	recovered := after - before
	if recovered > uint32(2147483647) {
		recovered = uint32(2147483647)
	}
	m.listener.OnMobDamaged(m, -int32(recovered))
}

func (m *Mob) ExpForDamage(damage uint64) uint32 {
	if m == nil || m.wz == nil || damage == 0 {
		return 0
	}
	maxHp := uint64(m.GetMaxHp())
	if maxHp == 0 {
		return 0
	}
	if damage > maxHp {
		damage = maxHp
	}
	base := uint64(m.wz.EXP)
	num := base * damage
	den := maxHp
	exp := (num + den - 1) / den
	const maxU32 = uint64(0xffffffff)
	if exp > maxU32 {
		return ^uint32(0)
	}
	return uint32(exp)
}
