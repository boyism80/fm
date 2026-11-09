package game

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Hit struct {
	Type       constant.IncomingHitType
	Damage     int32
	MobOID     uint32
	SkillID    uint32
	SkillLevel uint8
	Reflect    uint8
	Direction  uint8
}

const maxIncomingDamage = 200000

const blockRange = 400

type AttackKind uint8

const (
	AttackClose AttackKind = iota
	AttackRanged
	AttackMagic
	AttackEnergy
)

// Attack runs one attack packet: validation, skill costs, server attack hooks, damage, then the client update.
func (ch *Character) Attack(actx actor.Context, kind AttackKind, payload dto.AttackPayload) error {
	if ch.Spectating() {
		return nil
	}
	info := payload.ToAttackInfo()
	// An attack the client made on the field it left names that field's objects.
	if info.FieldKey != ch.fieldKey {
		return nil
	}
	if !ch.IsAlive() {
		if ch.GetMap() != nil {
			ch.cheat.Register(actx, CheatAttackWhileDead, info, "attacked while dead")
		}
		return nil
	}
	if ch.GetMap() == nil {
		return fmt.Errorf("attacker has no map")
	}

	ch.cheat.Pace(actx, info)
	var skill *SkillEntry
	if info.Skill != 0 {
		skill = ch.skills.Get(info.Skill)
	}
	bullet, attacked, err := ch.strike(actx, kind, info, skill)
	if err != nil {
		return err
	}
	if !attacked {
		ch.listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if kind != AttackEnergy {
		ch.markCriticals(info)
	}
	skillLevel := uint8(ch.skills.Level(info.Skill))
	switch kind {
	case AttackClose:
		ch.listener.OnAttack(ch, payload, skillLevel)
	case AttackRanged:
		ch.listener.OnRangedAttack(ch, payload, skillLevel, bullet)
	case AttackMagic:
		ch.listener.OnMagicAttack(ch, payload, skillLevel)
	case AttackEnergy:
		ch.listener.OnEnergyAttack(ch, payload)
	}
	if skill != nil && !followUpAttack(skill) {
		ch.skills.CallHook(actx, skill, "on_activated")
	}
	return nil
}

// markCriticals flags the hits the others see as critical, from the RNG stream the client was seeded with, as the
// original server's critModify did; the client draws a critical number for each flagged hit.
func (ch *Character) markCriticals(info dto.AttackInfo) {
	var randoms [7]int64
	for i := range randoms {
		randoms[i] = ch.random[0].CRand32Random()
	}
	prob, param := ch.criticalOdds(info)
	if prob <= 0 || param <= 0 {
		return
	}
	skill := constant.SkillID(info.Skill)
	if info.Display != criticalSkillDisplay {
		switch skill {
		case constant.SkillBlizzard3211003, constant.SkillShadowMeso, constant.SkillAssassinate, constant.SkillMesoExplosion:
			return
		}
	}
	index := 4
	if skill == constant.SkillSnipe || skill == constant.SkillNinjaStorm {
		index++
	}
	for i := range info.Damages {
		for j := range info.Damages[i].DamagePairs {
			info.Damages[i].DamagePairs[j].Unknown = float64(randoms[index%7]%10000000)*0.0000100000010000001 < float64(prob)
			index++
		}
	}
}

// criticalSkillDisplay is the attack display of a skill that brings its own critical odds.
const criticalSkillDisplay = 58

// criticalOdds is the critical chance in percent and the critical damage the attack can show: the skill's own for
// display 58, otherwise the bow or claw critical passive, plus Sharp Eyes.
func (ch *Character) criticalOdds(info dto.AttackInfo) (prob int, param int) {
	if info.Display == criticalSkillDisplay {
		if data := ch.skills.LevelData(constant.SkillID(info.Skill)); data != nil {
			prob, param = data.Prop, data.CriticalDamage
		}
	} else if weapon := ch.inventory.equipped[constant.EquipmentPartsWeapon]; weapon != nil {
		passive := constant.SkillID(0)
		if model, ok := weapon.GetModel().(*wz.Weapon); ok {
			switch model.WeaponType {
			case constant.WeaponTypeBow, constant.WeaponTypeCrossbow:
				passive = constant.SkillCriticalShot
			case constant.WeaponTypeClaw:
				passive = constant.SkillCriticalThrow
			}
		}
		if data := ch.skills.LevelData(passive); passive != 0 && data != nil {
			prob, param = data.Prop, data.Damage
		}
	}
	if _, sharpEyes, ok := ch.buffs.GetBuffValue(constant.BuffFlagSharpEyes); ok {
		prob += int(min(max(sharpEyes>>8, 0), 100))
		if sharpEyes > 0 {
			if param > 0 {
				param += int(sharpEyes & 0xFF)
			} else {
				param = int(sharpEyes&0xFF) + 100
			}
		}
	}
	return prob, param
}

// A follow-up attack is a hit the skill's buff triggers under the skill's id, so it neither costs nor activates the
// skill again.
func followUpAttack(skill *SkillEntry) bool {
	return skill != nil && skill.Level() > 0 && skill.wz.FollowUpAttack
}

// strike validates the reported hits before skill costs and server attack hooks, and returns the item the others see
// thrown by a ranged attack.
func (ch *Character) strike(actx actor.Context, kind AttackKind, info dto.AttackInfo, skill *SkillEntry) (uint32, bool, error) {
	switch kind {
	case AttackMagic:
		if info.Skill == 0 {
			return 0, false, nil
		}
	case AttackEnergy:
		_, energy, active := ch.buffs.GetBuffValue(constant.BuffFlagEnergyCharge)
		if !active || energy < 10000 || info.Skill != 0 {
			return 0, false, nil
		}
	}
	switch constant.SkillID(info.Skill) {
	case constant.SkillPanicSword, constant.SkillPanicAxe, constant.SkillComaSword, constant.SkillComaAxe, constant.SkillPanic, constant.SkillComa:
		// The combo value counts one more than the orbs, and a finisher needs at least one orb.
		if ch.buffs.value(constant.BuffFlagCombo) <= 1 {
			return 0, false, nil
		}
	}

	var err error
	switch {
	case kind == AttackMagic:
		err = NewMagicValidator(ch).Validate(info)
	case constant.SkillID(info.Skill) == constant.SkillMesoExplosion:
		err = NewMesoExplosionValidator(ch).Validate(info)
	default:
		err = NewPhysicalValidator(ch).Validate(info)
	}
	if err != nil {
		ch.cheat.Rejected(actx, info, err)
		return 0, false, err
	}

	// The projectiles are counted before the skill costs, so a shot that cannot leave costs nothing.
	var ammo uint16
	var bullet uint32
	if info.Ranged {
		count, thrown, ok := ch.ammo(info)
		if !ok {
			return 0, false, nil
		}
		ammo, bullet = count, thrown
	}
	if info.Skill != 0 && !followUpAttack(skill) && !ch.skills.UseAttack(actx, skill) {
		return 0, false, nil
	}
	if ammo > 0 && !ch.inventory.useAmmo(int16(info.Slot), ammo) {
		return 0, false, nil
	}

	ch.callAttackHooks(actx, kind, skill, info)
	ch.damageTo(info.Damages)
	if len(info.Damages) > 0 {
		ch.lieDetector.huntedAt = clock.Now()
	}
	if constant.SkillID(info.Skill) == constant.SkillMesoExplosion {
		ch.ExplodeMesos(info.MesoOIDs)
	}
	return bullet, true, nil
}

// Scripts may change the reported damage, so the hits are read back from the table they edited.
func (ch *Character) callAttackHooks(actx actor.Context, kind AttackKind, skill *SkillEntry, info dto.AttackInfo) {
	root := ch.GetMap().GetLuaRoot()
	if root == nil {
		return
	}

	damages := info.Damages
	damagesTable := ch.damagesTable(root, damages)
	attackInfoTable := root.NewTable()
	attackInfoTable.RawSetString("magic", lua.LBool(kind == AttackMagic))
	attackInfoTable.RawSetString("ranged", lua.LBool(kind == AttackRanged))
	if info.Position.X != 0 || info.Position.Y != 0 {
		attackInfoTable.RawSetString("x", lua.LNumber(info.Position.X))
		attackInfoTable.RawSetString("y", lua.LNumber(info.Position.Y))
	}
	ch.skills.CallHook(actx, skill, "on_attack", damagesTable, attackInfoTable)

	oidToPairs := make(map[uint32][]dto.DamagePair)
	damagesTable.ForEach(func(key lua.LValue, value lua.LValue) {
		ud, ok := key.(*lua.LUserData)
		if !ok || ud.Value == nil {
			return
		}
		mob, ok := ud.Value.(*Mob)
		if !ok {
			return
		}
		hitsTbl, ok := value.(*lua.LTable)
		if !ok {
			return
		}
		n := hitsTbl.Len()
		pairs := make([]dto.DamagePair, 0, n)
		for i := 1; i <= n; i++ {
			if num, ok := hitsTbl.RawGetInt(i).(lua.LNumber); ok {
				pairs = append(pairs, dto.DamagePair{Damage: uint32(num), Unknown: false})
			}
		}
		oidToPairs[mob.GetOID()] = pairs
	})
	for i := range damages {
		if pairs, ok := oidToPairs[damages[i].OID]; ok {
			damages[i].DamagePairs = pairs
		}
	}
}

func (ch *Character) damagesTable(L *lua.LState, damages []dto.AttackPair) *lua.LTable {
	tbl := L.NewTable()
	m := ch.GetMap()
	if m == nil {
		return tbl
	}
	for _, ap := range damages {
		mob := m.GetMob(ap.OID)
		if mob == nil {
			continue
		}
		hits := L.NewTable()
		for i, dp := range ap.DamagePairs {
			hits.RawSetInt(i+1, lua.LNumber(dp.Damage))
		}
		tbl.RawSet(luax.NewLuable(L, mob), hits)
	}
	return tbl
}

// projectile is the ammo a ranged attack throws from info.Slot, nil when the attack throws none. ok is false when the
// slot does not hold the ammo the weapon needs.
func (ch *Character) projectile(info dto.AttackInfo) (item Item, projectile *wz.Consume, ok bool) {
	switch constant.SkillID(info.Skill) {
	case constant.SkillShadowMeso, constant.SkillBarrage, constant.SkillSoulBlade, constant.SkillSpark, constant.SkillSharkWave, constant.SkillShowdown4221003:
		return nil, nil, true
	}
	weapon := ch.inventory.equipped[constant.EquipmentPartsWeapon]
	if weapon == nil {
		return nil, nil, true
	}
	weaponModel, ok := weapon.GetModel().(*wz.Weapon)
	if !ok {
		return nil, nil, true
	}

	var ammo constant.ConsumeType
	switch weaponModel.WeaponType {
	case constant.WeaponTypeBow:
		ammo = constant.ConsumeTypeArrowBow
	case constant.WeaponTypeCrossbow:
		ammo = constant.ConsumeTypeArrowCrossBow
	case constant.WeaponTypeClaw:
		ammo = constant.ConsumeTypeShuriken
	case constant.WeaponTypeGun:
		ammo = constant.ConsumeTypeBullet
	default:
		return nil, nil, true
	}
	if ch.buffs.Has(constant.BuffFlagSoulArrow) && (ammo == constant.ConsumeTypeArrowBow || ammo == constant.ConsumeTypeArrowCrossBow) {
		return nil, nil, true
	}
	item = ch.inventory.GetItem(constant.InventoryTypeConsume, int16(info.Slot))
	if item == nil || item.GetCount() == 0 {
		return nil, nil, false
	}
	projectile, ok = item.GetModel().(*wz.Consume)
	if !ok || projectile.ConsumeType != ammo {
		return nil, nil, false
	}
	return item, projectile, true
}

// ammo is how many projectiles a ranged attack takes from info.Slot, 0 when it takes none, and the item the others
// see thrown: the cash star in info.CsStar when the client names one, otherwise the ammo itself.
func (ch *Character) ammo(info dto.AttackInfo) (uint16, uint32, bool) {
	item, projectile, ok := ch.projectile(info)
	if !ok || projectile == nil {
		return 0, 0, ok
	}
	bullet := item.GetModel().GetID()
	if info.CsStar != 0 {
		cash := ch.inventory.GetItem(constant.InventoryTypeCash, int16(info.CsStar))
		if cash == nil {
			return 0, 0, false
		}
		bullet = cash.GetModel().GetID()
	}
	if ch.buffs.Has(constant.BuffFlagSpiritClaw) && projectile.ConsumeType == constant.ConsumeTypeShuriken {
		return 0, bullet, true
	}

	count := 1
	if p := ch.skills.LevelData(constant.SkillID(info.Skill)); p != nil {
		count = max(1, p.BulletCount, p.AttackCount)
		if p.BulletConsume > 0 {
			count = p.BulletConsume
		}
	}
	if ch.buffs.Has(constant.BuffFlagShadowPartner) {
		count *= 2
	}
	if item.GetCount() < uint16(count) {
		return 0, 0, false
	}
	return uint16(count), bullet, true
}

func (ch *Character) HypnotizeDamage(sourceOID, targetOID, damage uint32) bool {
	if ch.Spectating() || !ch.IsAlive() || damage == 0 {
		return false
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return false
	}
	source := mapInstance.GetMob(sourceOID)
	target := mapInstance.GetMob(targetOID)
	if source == nil || target == nil || source == target {
		return false
	}
	causerID, hypnotized := source.buffs.Causer(constant.MobBuffHypnotize)
	if !hypnotized || causerID != ch.GetID() || target.wz.DamagedByMob {
		return false
	}
	if !source.canAttack(target, damage) {
		return false
	}
	if !target.TakeDamage(ch, damage) {
		target.listener.OnMobDamaged(target, int32(damage))
	}
	return true
}

func (ch *Character) damageTo(damages []dto.AttackPair) {
	if len(damages) == 0 {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	for _, damage := range damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil || !mob.IsAlive() {
			continue
		}
		total := uint64(0)
		for _, damagePair := range damage.DamagePairs {
			total += uint64(damagePair.Damage)
		}
		if total > 0 {
			ch.dojo.charge(mob, total)
		}
		for _, damagePair := range damage.DamagePairs {
			if damagePair.Damage == 0 {
				continue
			}
			mob.TakeDamage(ch, damagePair.Damage)
		}
	}
}

const mesoExplosionMargin = 50

func (ch *Character) ExplodeMesos(oids []uint32) {
	if len(oids) == 0 {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	area, ok := ch.mesoExplosionArea()
	if !ok {
		return
	}

	items := mapInstance.GetItems()
	for _, oid := range oids {
		meso, ok := items[oid].(*Meso)
		if !ok || !meso.inArea(area) {
			continue
		}
		if err := mapInstance.RemoveItem(oid, constant.RemoveItemTypeExplosion, ch.GetID()); err != nil {
			log.Printf("failed to remove meso oid %d for explosion: %v", oid, err)
		}
	}
}

func (ch *Character) mesoExplosionArea() (types.Rect[int32], bool) {
	skillID := uint32(constant.SkillMesoExplosion)
	skill := ch.gameWorld.GetResources().GetSkill(skillID)
	if skill == nil {
		return types.Rect[int32]{}, false
	}
	levelData := skill.GetLevelData(ch.skills.Level(skillID))
	if levelData == nil {
		return types.Rect[int32]{}, false
	}

	reach := max(levelData.LT.X, -levelData.LT.X, levelData.RB.X, -levelData.RB.X) + mesoExplosionMargin
	origin := types.Point[int32]{X: int32(ch.position.X), Y: int32(ch.position.Y)}
	return types.Rect[int32]{
		Left:   -reach,
		Top:    min(levelData.LT.Y, levelData.RB.Y) - mesoExplosionMargin,
		Right:  reach,
		Bottom: max(levelData.LT.Y, levelData.RB.Y) + mesoExplosionMargin,
	}.AtOrigin(origin), true
}

func (ch *Character) TakeHit(actx actor.Context, hit Hit) {
	if hit.Damage < -1 || hit.Damage > maxIncomingDamage {
		return
	}

	m := ch.GetMap()
	var mob *Mob
	if m != nil && hit.MobOID != 0 {
		mob = m.GetMob(hit.MobOID)
	}
	wasAlive := ch.IsAlive()

	var attack *wz.MobAttack
	if index, ok := hit.Type.MobAttackIndex(); ok && mob != nil && !mob.fake && hit.Damage > 0 {
		attack = mob.wz.Attack(index)
	}
	banish := hit.Type == constant.IncomingHitCollide && hit.Damage >= 0 && mob != nil && !mob.fake && mob.wz.Banish != nil && mob.wz.Banish.OnTouch
	switch {
	case attack != nil:
		if attack.DiseaseSkill != 0 {
			ch.debuffs.GiveMobSkill(attack.DiseaseSkill, attack.DiseaseLevel)
		}
		mp := mob.GetMp()
		mob.SetMp(mp-min(mp, uint32(max(attack.MpCon, 0))), false)
	case hit.Type == constant.IncomingHitMapDebuff && hit.SkillID != 0 && hit.SkillLevel != 0:
		ch.debuffs.GiveMobSkill(hit.SkillID, hit.SkillLevel)
	}

	damage := int32(0)
	if hit.Damage == -1 && hit.Type == constant.IncomingHitCollide {
		ch.block(actx, hit)
	} else {
		damage = ch.hitDamage(hit, mob, attack)
	}

	if ch.invincible {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{}, true)
	} else {
		ch.showHit(hit, mob, wasAlive)
		if mob != nil && damage > 0 {
			ch.dojo.AddEnergy(2)
		}
		ch.TakeDamage(damage)
	}

	// A banish to another actor's map hands the character over, so nothing may touch it afterwards.
	if banish {
		mob.Banish(actx, ch)
	}
}

// showHit shows the others the hit as the client reported it, before guards; a hidden or already dead character and a
// hit from a fake or missing mob show nothing. The mob id is the server's, not the client's.
func (ch *Character) showHit(hit Hit, mob *Mob, wasAlive bool) {
	if !wasAlive || ch.IsHidden() || (hit.MobOID != 0 && (mob == nil || mob.fake)) {
		return
	}
	mobID := uint32(0)
	if mob != nil {
		mobID = mob.wz.ID
	}
	fakeSkillID := uint32(0)
	if hit.Damage == -1 {
		fakeSkillID = uint32(constant.SkillShadowShifter)
		if ch.class.Of(constant.ClassBandit) {
			fakeSkillID = uint32(constant.SkillShadowShifter4220002)
		}
	}
	ch.listener.OnDamaged(ch, hit, mobID, fakeSkillID)
}

func (ch *Character) hitDamage(hit Hit, mob *Mob, attack *wz.MobAttack) int32 {
	m := ch.GetMap()
	if m == nil {
		return hit.Damage
	}
	if mob != nil && mob.fake {
		return 0
	}
	// Smoke Bomb keeps all damage off, so no guard spends MP or meso either.
	if m.inSmoke(ch) {
		return 0
	}
	root := m.GetLuaRoot()
	if root == nil {
		return hit.Damage
	}

	var attackerArg interface{} = lua.LNil
	if mob != nil {
		attackerArg = mob
	}
	var skillArg interface{} = lua.LNil
	if hit.SkillID != 0 {
		if model := ch.gameWorld.GetResources().GetSkill(hit.SkillID); model != nil {
			skillArg = NewSkillEntry(ch, model, int(hit.SkillLevel), 0)
		}
	}
	params := root.NewTable()
	params.RawSetString("hit_type", lua.LNumber(int32(hit.Type)))
	params.RawSetString("reflect_ratio", lua.LNumber(int32(hit.Reflect)))
	if attack != nil {
		params.RawSetString("deadly", lua.LBool(attack.DeadlyAttack))
		params.RawSetString("mp_burn", lua.LNumber(attack.MpBurn))
	}

	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_damaged: %v", err)
		return hit.Damage
	}
	ret, err := luax.Call(thread, "on_damaged", ch, attackerArg, skillArg, hit.Damage, params)
	if err != nil {
		log.Printf("Failed to call script on_damaged: %v", err)
		return hit.Damage
	}
	if ret == nil || ret.Type() != lua.LTNumber {
		return hit.Damage
	}
	return max(int32(lua.LVAsNumber(ret)), 0)
}

func (ch *Character) block(actx actor.Context, hit Hit) {
	m := ch.GetMap()
	if m == nil {
		return
	}
	root := m.GetLuaRoot()
	if root == nil {
		return
	}
	var attackerArg interface{} = lua.LNil
	if hit.MobOID != 0 {
		if mob := m.GetMob(hit.MobOID); mob != nil && mob.GetPosition().DistanceSq(ch.position) <= blockRange*blockRange {
			attackerArg = mob
		}
	}
	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_blocked: %v", err)
		return
	}
	luax.CallAsync(actx, root, thread, "on_blocked", ch, attackerArg).OnError(func(err error) {
		log.Printf("Failed to call script on_blocked: %v", err)
	})
}

func (ch *Character) TakeDamage(damage int32) {
	if damage > 0 {
		ch.mount.takeDamage(damage)
	}

	hp := int(ch.GetHp()) - int(damage)
	if hp < 0 {
		hp = 0
	}
	if maxHp := int(ch.GetMaxHp()); hp > maxHp {
		hp = maxHp
	}
	ch.SetHp(uint32(hp), false)
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
	}, true)
}
