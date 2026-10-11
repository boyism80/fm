package entity

import (
	"math"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) damageBuff(flag constant.BuffFlag) float64 {
	_, v, ok := ch.Buffs.GetBuffValue(flag)
	if !ok {
		return 0
	}
	return float64(v)
}

func (ch *Character) skillLevelData(id constant.SkillID) *wz.SkillLevelData {
	entry := ch.Skills.Get(uint32(id))
	if entry == nil || entry.Level() <= 0 {
		return nil
	}
	return entry.Wz.GetLevelData(entry.Level())
}

// comboPercent implements the client's stored-count minus one convention.
func (ch *Character) comboPercent(finisher bool) float64 {
	baseID, advancedID, capLevel := constant.SkillComboAttack, constant.SkillAdvancedCombo, 30
	if ch.Class >= 1000 {
		baseID, advancedID, capLevel = constant.SkillComboAttackCygnus, constant.SkillAdvancedComboCygnus, 20
	}
	skill := ch.skillLevelData(baseID)
	if skill == nil {
		return 100
	}
	level := ch.GetTotalSkillLevel(uint32(baseID))
	if level >= capLevel {
		if advanced := ch.skillLevelData(advancedID); advanced != nil {
			skill = advanced
		}
	}

	c := int(ch.damageBuff(constant.BuffFlagCombo)) - 1
	d := float64(skill.Damage)
	l := float64(level)
	if finisher {
		if c >= 2 && c <= 5 {
			return d + math.Trunc([]float64{24, 64, 120, 184}[c-2]*(l-1)/29) + []float64{6, 16, 30, 46}[c-2]
		}
		if c >= 6 && c <= 10 {
			return d + []float64{250, 270, 290, 305, 320}[c-6]
		}
	}
	if c == 1 {
		return d
	}
	if c >= 2 && c <= 5 {
		return d + math.Trunc(5*float64(c-1)*(l-1)/29)
	}
	if c >= 6 && c <= 10 {
		return d + []float64{24, 28, 32, 36, 40}[c-6]
	}
	if c < 0 {
		return 0
	}
	return 100
}

// damageStats snapshots what the client damage formula reads: the basic stats
// actually serialized to this client, equipment, passive skills and buffs.
func (ch *Character) damageStats() (damageStats, damageModifiers) {
	s := damageStats{
		STR:   float64(ch.Stats.TotalStr()),
		DEX:   float64(ch.Stats.TotalDex()),
		INT:   float64(ch.Stats.TotalInt()),
		LUK:   float64(ch.Stats.TotalLuk()),
		Level: int(ch.GetLevel()),
		Class: constant.ClassType(ch.Class),
	}
	base := s

	// Include all equipped bonuses as an envelope when client eligibility is
	// unavailable; ineligible pieces can only increase this nonnegative bound.
	for slot, item := range ch.Inventory.Equipped {
		// Cash appearance items have no stats.
		if slot <= -100 || item == nil {
			continue
		}
		e := ToEquipmentDTOFromCore(item.GetEquipmentCore())
		if e == nil {
			continue
		}
		s.STR += float64(e.Str)
		s.DEX += float64(e.Dex)
		s.INT += float64(e.Int)
		s.LUK += float64(e.Luk)
		s.PAD += float64(e.PAD)
		s.MAD += float64(e.MAD)
		s.ACC += float64(e.ACC)
		if slot == constant.EquipmentPartsWeapon {
			s.Weapon = constant.GetWeaponType(e.ItemId)
		}
	}

	mapleWarrior := ch.damageBuff(constant.BuffFlagMapleWarrior)
	s.STR += math.Trunc(base.STR * mapleWarrior / 100)
	s.DEX += math.Trunc(base.DEX * mapleWarrior / 100)
	s.INT += math.Trunc(base.INT * mapleWarrior / 100)
	s.LUK += math.Trunc(base.LUK * mapleWarrior / 100)
	s.MAD += s.INT
	if s.Weapon == constant.WeaponTypeNone && s.Class%1000/100 == 5 {
		s.PAD += min(math.Trunc(.7*float64(s.Level)+10), 31)
	}

	switch {
	case s.Class%1000/10 == 51:
		s.ACC += math.Trunc(.9*s.DEX + .3*s.LUK)
	case s.Class%1000/100 >= 3 && s.Class%1000/100 <= 5:
		s.ACC += math.Trunc(.6*s.DEX + .3*s.LUK)
	default:
		s.ACC += math.Trunc(.8*s.DEX + .5*s.LUK)
	}
	accuracySkills := []constant.SkillID{
		constant.SkillTheBlessingOfAmazon,
		constant.SkillNimbleBody,
		constant.SkillNimbleBodyCygnus,
		constant.SkillBulletTime,
		constant.SkillBulletTimeCygnus,
	}
	for _, id := range accuracySkills {
		if p := ch.skillLevelData(id); p != nil {
			s.ACC += float64(p.X)
		}
	}
	for _, id := range []constant.SkillID{constant.SkillBlessingOfTheFairy, constant.SkillBlessingOfTheFairyCygnus} {
		if p := ch.skillLevelData(id); p != nil {
			s.PAD += float64(p.X)
			s.MAD += float64(p.Y)
			s.ACC += float64(p.Z)
		}
	}

	b := damageModifiers{
		PAD:          ch.damageBuff(constant.BuffFlagWeaponAtk),
		MAD:          ch.damageBuff(constant.BuffFlagMagicAtk),
		ACC:          ch.damageBuff(constant.BuffFlagAcc),
		Echo:         ch.damageBuff(constant.BuffFlagEchoOfHero),
		SharpEyes:    ch.damageBuff(constant.BuffFlagSharpEyes),
		ElementReset: ch.damageBuff(constant.BuffFlagElementReset),
		Double:       ch.damageBuff(constant.BuffFlagBerserkFury) > 0,
		GM:           s.Class%1000/100 == 9 || s.Class >= 28 && s.Class <= 30,
	}
	if ch.Debuffs != nil {
		b.Sealed = ch.Debuffs.Has(constant.DebuffFlagSeal)
	}

	for _, id := range ch.masterySkills(s.Weapon) {
		if p := ch.skillLevelData(id); p != nil {
			b.Mastery = max(b.Mastery, float64(p.Mastery))
			s.PAD += float64(p.X)
			s.ACC += float64(p.X)
		}
	}
	if s.Weapon == constant.WeaponTypeSpear || s.Weapon == constant.WeaponTypePoleArm {
		if p := ch.skillLevelData(constant.SkillBeholder); p != nil {
			b.Mastery += float64(p.Mastery)
		}
	}

	b.CriticalPercent = 100
	if critical := ch.criticalSkill(s.Weapon); critical != nil {
		b.CriticalChance = float64(critical.Prop)
		b.CriticalPercent = float64(critical.Damage)
	}

	// Energy charge's PAD is active only in the charged state. The server buff
	// value is authoritative; include its source PAD for an upper envelope.
	if ch.damageBuff(constant.BuffFlagEnergyCharge) > 0 {
		energyCharge := constant.SkillEnergyCharge
		if s.Class >= 1000 {
			energyCharge = constant.SkillEnergyChargeCygnus
		}
		if p := ch.skillLevelData(energyCharge); p != nil {
			b.EnergyPAD = float64(p.PAD)
		}
	}

	amplificationSkills := []constant.SkillID{
		constant.SkillElementAmplification,
		constant.SkillElementAmplification2210001,
		constant.SkillElementAmplificationCygnus,
	}
	for _, id := range amplificationSkills {
		if p := ch.skillLevelData(id); p != nil {
			b.Amplification = max(b.Amplification, float64(p.Y))
		}
	}

	if ch.damageBuff(constant.BuffFlagCombo) > 0 {
		b.ComboPercent = ch.comboPercent(false)
	}

	if charge, ok := ch.Buffs.GetEntity(constant.BuffFlagWkCharge).(*SkillBuff); ok {
		if p := charge.Wz.GetLevelData(int(charge.SkillLevel)); p != nil {
			b.ChargeDamage = float64(p.Damage)
			b.ChargeStrength = float64(p.Z)
		}
		switch constant.SkillID(charge.Wz.ID) {
		case constant.SkillFireChargeSword, constant.SkillFireChargeBw:
			b.ChargeElement = "f"
		case constant.SkillIceChargeSword, constant.SkillBlizzardChargeBw:
			b.ChargeElement = "i"
		case constant.SkillThunderChargeSword, constant.SkillLightningChargeBw, constant.SkillRadiantCharge, constant.SkillLightningCharge:
			b.ChargeElement = "l"
		case constant.SkillHolyChargeSword, constant.SkillDivineChargeMace:
			b.ChargeElement = "h"
		}
	}

	if berserk := ch.skillLevelData(constant.SkillBerserk); s.Class == constant.ClassDarkKnight && berserk != nil {
		if float64(ch.GetHp()) <= math.Trunc(float64(ch.GetMaxHp())*float64(berserk.X)/100) {
			b.Berserk = float64(berserk.Damage)
		}
	}

	return s, b
}

func (ch *Character) masterySkills(weapon constant.WeaponType) []constant.SkillID {
	switch weapon {
	case constant.WeaponTypeSword1H, constant.WeaponTypeSword2H:
		return []constant.SkillID{constant.SkillSwordMastery, constant.SkillSwordMastery1200000, constant.SkillSwordMasteryCygnus}
	case constant.WeaponTypeAxe1H, constant.WeaponTypeAxe2H:
		return []constant.SkillID{constant.SkillAxeMastery}
	case constant.WeaponTypeBlunt1H, constant.WeaponTypeBlunt2H:
		return []constant.SkillID{constant.SkillBwMastery}
	case constant.WeaponTypeDagger:
		return []constant.SkillID{constant.SkillDaggerMastery}
	case constant.WeaponTypeSpear:
		return []constant.SkillID{constant.SkillSpearMastery}
	case constant.WeaponTypePoleArm:
		return []constant.SkillID{constant.SkillPoleArmMastery}
	case constant.WeaponTypeBow:
		if ch.Class >= 1000 {
			return []constant.SkillID{constant.SkillBowMasteryCygnus, constant.SkillBowExpertCygnus}
		}
		return []constant.SkillID{constant.SkillBowMastery, constant.SkillBowExpert}
	case constant.WeaponTypeCrossbow:
		return []constant.SkillID{constant.SkillCrossbowMastery, constant.SkillCrossbowExpert}
	case constant.WeaponTypeClaw:
		if ch.Class >= 1000 {
			return []constant.SkillID{constant.SkillClawMasteryCygnus}
		}
		return []constant.SkillID{constant.SkillClawMastery}
	case constant.WeaponTypeKnuckle:
		if ch.Class >= 1000 {
			return []constant.SkillID{constant.SkillKnuckleMasteryCygnus}
		}
		return []constant.SkillID{constant.SkillKnuckleMastery}
	case constant.WeaponTypeGun:
		return []constant.SkillID{constant.SkillGunMastery}
	}
	return nil
}

func (ch *Character) criticalSkill(weapon constant.WeaponType) *wz.SkillLevelData {
	switch weapon {
	case constant.WeaponTypeBow, constant.WeaponTypeCrossbow:
		if ch.Class >= 1000 {
			return ch.skillLevelData(constant.SkillCriticalShotCygnus)
		}
		return ch.skillLevelData(constant.SkillCriticalShot)
	case constant.WeaponTypeClaw:
		if ch.Class >= 1000 {
			return ch.skillLevelData(constant.SkillCriticalThrowCygnus)
		}
		return ch.skillLevelData(constant.SkillCriticalThrow)
	case constant.WeaponTypeKnuckle:
		return ch.skillLevelData(constant.SkillCriticalPunch)
	}
	return nil
}
