package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
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
}

// Attack validates the reported hits before skill costs and server attack hooks.
func (ch *Character) Attack(info dto.AttackInfo, magic bool, activating func() bool, attacking func()) (bool, error) {
	if ch.GetMap() == nil {
		return false, fmt.Errorf("attacker has no map")
	}
	if ch.Spectating() || !ch.IsAlive() {
		return false, nil
	}

	var err error
	switch {
	case magic:
		err = NewMagicValidator(ch).Validate(info)
	case constant.SkillID(info.Skill) == constant.SkillMesoExplosion:
		err = NewMesoExplosionValidator(ch).Validate(info)
	default:
		err = NewPhysicalValidator(ch).Validate(info)
	}
	if err != nil {
		return false, err
	}

	if info.Skill != 0 && !ch.Skills.UseAttack(info.Skill, activating) {
		return false, nil
	}
	if info.Ranged && !ch.shoot(info) {
		return false, nil
	}

	attacking()
	ch.damageTo(info.Damages)
	if constant.SkillID(info.Skill) == constant.SkillMesoExplosion {
		ch.ExplodeMesos(info.MesoOIDs)
	}
	return true, nil
}

func (ch *Character) shoot(info dto.AttackInfo) bool {
	switch constant.SkillID(info.Skill) {
	case constant.SkillShadowMeso, constant.SkillBarrage, constant.SkillSoulBlade, constant.SkillSpark, constant.SkillSharkWave, constant.SkillShowdown4221003:
		return true
	}
	weapon := ch.Inventory.Equipped[constant.EquipmentPartsWeapon]
	if weapon == nil {
		return true
	}
	weaponModel, ok := weapon.GetModel().(*wz.Weapon)
	if !ok {
		return true
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
		return true
	}
	if ch.Buffs.Has(constant.BuffFlagSoulArrow) && (ammo == constant.ConsumeTypeArrowBow || ammo == constant.ConsumeTypeArrowCrossBow) {
		return true
	}
	item := ch.Inventory.GetItem(constant.InventoryTypeConsume, int16(info.Slot))
	if item == nil {
		return false
	}
	projectile, ok := item.GetModel().(*wz.Consume)
	if !ok || projectile.ConsumeType != ammo {
		return false
	}
	if ch.Buffs.Has(constant.BuffFlagSpiritClaw) && ammo == constant.ConsumeTypeShuriken {
		return true
	}

	count := 1
	if p := ch.Skills.LevelData(constant.SkillID(info.Skill)); p != nil {
		count = max(1, p.BulletCount, p.AttackCount)
		if p.BulletConsume > 0 {
			count = p.BulletConsume
		}
	}
	if ch.Buffs.Has(constant.BuffFlagShadowPartner) {
		count *= 2
	}
	return ch.Inventory.useAmmo(int16(info.Slot), uint16(count))
}

func (ch *Character) EnergyAttack(info dto.AttackInfo, attacking func()) (bool, error) {
	_, energy, active := ch.Buffs.GetBuffValue(constant.BuffFlagEnergyCharge)
	if !active || energy < 10000 || info.Skill != 0 {
		return false, nil
	}
	return ch.Attack(info, false, nil, attacking)
}

func (ch *Character) SelfDestructMob(oid uint32) bool {
	if ch.Spectating() || ch.IsHidden() || !ch.IsAlive() {
		return false
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return false
	}
	mob := mapInstance.GetMob(oid)
	if mob == nil {
		return false
	}
	return mob.SelfDestruct(ch)
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
	causerID, hypnotized := source.Buffs.Causer(constant.MobBuffHypnotize)
	if !hypnotized || causerID != ch.GetID() || target.Wz.DamagedByMob {
		return false
	}
	if !source.canAttack(target, damage) {
		return false
	}
	if target.TakeDamage(ch, damage) {
		ch.OnKill([]*Mob{target})
	} else {
		target.Listener.OnMobDamaged(target, int32(damage))
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

	var killed []*Mob
	seen := make(map[uint32]struct{})
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
			ch.chargeDojoEnergy(mob, total)
		}
		for _, damagePair := range damage.DamagePairs {
			if damagePair.Damage == 0 {
				continue
			}
			if !mob.TakeDamage(ch, damagePair.Damage) {
				continue
			}
			if _, ok := seen[mob.OID]; ok {
				continue
			}
			seen[mob.OID] = struct{}{}
			killed = append(killed, mob)
		}
	}
	if len(killed) > 0 {
		ch.OnKill(killed)
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
	skill := ch.GameWorld.GetResources().GetSkill(skillID)
	if skill == nil {
		return types.Rect[int32]{}, false
	}
	levelData := skill.GetLevelData(ch.Skills.Level(skillID))
	if levelData == nil {
		return types.Rect[int32]{}, false
	}

	reach := max(levelData.LT.X, -levelData.LT.X, levelData.RB.X, -levelData.RB.X) + mesoExplosionMargin
	origin := types.Point[int32]{X: int32(ch.Position.X), Y: int32(ch.Position.Y)}
	return types.Rect[int32]{
		Left:   -reach,
		Top:    min(levelData.LT.Y, levelData.RB.Y) - mesoExplosionMargin,
		Right:  reach,
		Bottom: max(levelData.LT.Y, levelData.RB.Y) + mesoExplosionMargin,
	}.AtOrigin(origin), true
}

func (ch *Character) TakeHit(actx actor.Context, hit Hit) {
	damage := int32(0)
	if hit.Damage == -1 && hit.Type == constant.IncomingHitCollide {
		ch.block(actx, hit)
	} else {
		damage = ch.hitDamage(hit)
	}

	if ch.Invincible {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{}, true)
		return
	}
	ch.TakeDamage(damage)
}

func (ch *Character) hitDamage(hit Hit) int32 {
	m := ch.GetMap()
	if m == nil {
		return hit.Damage
	}
	var mob *Mob
	if hit.MobOID != 0 {
		mob = m.GetMob(hit.MobOID)
	}
	if mob != nil && mob.Fake {
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
		if model := ch.GameWorld.GetResources().GetSkill(hit.SkillID); model != nil {
			skillArg = NewSkillEntry(ch, model, int(hit.SkillLevel), 0)
		}
	}
	params := root.NewTable()
	params.RawSetString("hit_type", lua.LNumber(int32(hit.Type)))
	params.RawSetString("reflect_ratio", lua.LNumber(int32(hit.Reflect)))

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
		if mob := m.GetMob(hit.MobOID); mob != nil {
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
		ch.AddDojoEnergy(2)
		ch.Mount.takeDamage(damage)
	}

	hp := int(ch.GetHp()) - int(damage)
	if hp < 0 {
		hp = 0
	}
	if maxHp := int(ch.GetMaxHp()); hp > maxHp {
		hp = maxHp
	}
	ch.SetHp(uint32(hp), false)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
	}, true)
}

func (ch *Character) OnKill(mobs []*Mob) {
	if ch == nil || len(mobs) == 0 {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	byID := make(map[uint32][]luax.Luable)
	order := make([]uint32, 0)
	all := make([]luax.Luable, 0, len(mobs))
	for _, mob := range mobs {
		if mob == nil || mob.Wz == nil {
			continue
		}
		id := mob.Wz.ID
		if _, ok := byID[id]; !ok {
			order = append(order, id)
		}
		byID[id] = append(byID[id], mob)
		all = append(all, mob)
	}
	if len(all) == 0 {
		return
	}

	for _, id := range order {
		group := byID[id]
		scriptPath := fmt.Sprintf("script/mob/%d.lua", id)
		mapInstance.runMobLuaHook(root, scriptPath, "on_mob_kill", ch, group)
	}
	mapInstance.runMobLuaHook(root, constant.CharacterHookScriptPath, "on_mob_kill", ch, all)
	if sm := ch.StateMachine(); sm != nil {
		sm.CallHook("on_mob_kill", ch, all)
	}
}
