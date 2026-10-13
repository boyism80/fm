package entity

import (
	"math"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

// comboPercent implements the client's stored-count minus one convention.
func (ch *Character) comboPercent(finisher bool) float64 {
	baseID, advancedID, capLevel := constant.SkillComboAttack, constant.SkillAdvancedCombo, 30
	if ch.Class >= constant.ClassNoblesse {
		baseID, advancedID, capLevel = constant.SkillComboAttackCygnus, constant.SkillAdvancedComboCygnus, 20
	}
	skill := ch.Skills.LevelData(baseID)
	if skill == nil {
		return 100
	}
	level := ch.Skills.Level(uint32(baseID))
	if level >= capLevel {
		if advanced := ch.Skills.LevelData(advancedID); advanced != nil {
			skill = advanced
		}
	}

	c := int(ch.Buffs.value(constant.BuffFlagCombo)) - 1
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
		Class: ch.Class,
	}
	base := s

	// Include all equipped bonuses as an envelope when client eligibility is
	// unavailable; ineligible pieces can only increase this nonnegative bound.
	for slot, item := range ch.Inventory.Equipped {
		// Cash appearance items have no stats.
		if slot <= -100 || item == nil {
			continue
		}
		e := item.GetEquipmentCore().BuildEquipmentS2C()
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
		if weapon, ok := item.GetModel().(*wz.Weapon); ok && slot == constant.EquipmentPartsWeapon {
			s.Weapon = weapon.WeaponType
		}
	}

	mapleWarrior := ch.Buffs.value(constant.BuffFlagMapleWarrior)
	s.STR += math.Trunc(base.STR * mapleWarrior / 100)
	s.DEX += math.Trunc(base.DEX * mapleWarrior / 100)
	s.INT += math.Trunc(base.INT * mapleWarrior / 100)
	s.LUK += math.Trunc(base.LUK * mapleWarrior / 100)
	s.MAD += s.INT
	if s.Weapon == constant.WeaponTypeNone && (s.Class.Of(constant.ClassPirate) || s.Class.Of(constant.ClassThunderBreaker1)) {
		s.PAD += min(math.Trunc(.7*float64(s.Level)+10), 31)
	}

	switch {
	case s.Class.Of(constant.ClassBrawler) || s.Class.Of(constant.ClassThunderBreaker2):
		s.ACC += math.Trunc(.9*s.DEX + .3*s.LUK)
	case s.Class.Of(constant.ClassBowman) || s.Class.Of(constant.ClassWindArcher1),
		s.Class.Of(constant.ClassThief) || s.Class.Of(constant.ClassNightWalker1),
		s.Class.Of(constant.ClassPirate) || s.Class.Of(constant.ClassThunderBreaker1):
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
		if p := ch.Skills.LevelData(id); p != nil {
			s.ACC += float64(p.X)
		}
	}
	for _, id := range []constant.SkillID{constant.SkillBlessingOfTheFairy, constant.SkillBlessingOfTheFairyCygnus} {
		if p := ch.Skills.LevelData(id); p != nil {
			s.PAD += float64(p.X)
			s.MAD += float64(p.Y)
			s.ACC += float64(p.Z)
		}
	}

	b := damageModifiers{
		PAD:          ch.Buffs.value(constant.BuffFlagWeaponAtk),
		MAD:          ch.Buffs.value(constant.BuffFlagMagicAtk),
		ACC:          ch.Buffs.value(constant.BuffFlagAcc),
		Echo:         ch.Buffs.value(constant.BuffFlagEchoOfHero),
		SharpEyes:    ch.Buffs.value(constant.BuffFlagSharpEyes),
		ElementReset: ch.Buffs.value(constant.BuffFlagElementReset),
		Double:       ch.Buffs.value(constant.BuffFlagBerserkFury) > 0,
		GM:           s.Class.Of(constant.ClassGM) || s.Class >= 28 && s.Class <= 30,
	}
	if ch.Debuffs != nil {
		b.Sealed = ch.Debuffs.Has(constant.DebuffFlagSeal)
	}

	for _, id := range ch.masterySkills(s.Weapon) {
		if p := ch.Skills.LevelData(id); p != nil {
			b.Mastery = max(b.Mastery, float64(p.Mastery))
			s.PAD += float64(p.X)
			s.ACC += float64(p.X)
		}
	}
	if s.Weapon == constant.WeaponTypeSpear || s.Weapon == constant.WeaponTypePoleArm {
		if p := ch.Skills.LevelData(constant.SkillBeholder); p != nil {
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
	if ch.Buffs.value(constant.BuffFlagEnergyCharge) > 0 {
		energyCharge := constant.SkillEnergyCharge
		if s.Class >= constant.ClassNoblesse {
			energyCharge = constant.SkillEnergyChargeCygnus
		}
		if p := ch.Skills.LevelData(energyCharge); p != nil {
			b.EnergyPAD = float64(p.PAD)
		}
	}

	amplificationSkills := []constant.SkillID{
		constant.SkillElementAmplification,
		constant.SkillElementAmplification2210001,
		constant.SkillElementAmplificationCygnus,
	}
	for _, id := range amplificationSkills {
		if p := ch.Skills.LevelData(id); p != nil {
			b.Amplification = max(b.Amplification, float64(p.Y))
		}
	}

	if ch.Buffs.value(constant.BuffFlagCombo) > 0 {
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

	if berserk := ch.Skills.LevelData(constant.SkillBerserk); s.Class == constant.ClassDarkKnight && berserk != nil {
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
		if ch.Class >= constant.ClassNoblesse {
			return []constant.SkillID{constant.SkillBowMasteryCygnus, constant.SkillBowExpertCygnus}
		}
		return []constant.SkillID{constant.SkillBowMastery, constant.SkillBowExpert}
	case constant.WeaponTypeCrossbow:
		return []constant.SkillID{constant.SkillCrossbowMastery, constant.SkillCrossbowExpert}
	case constant.WeaponTypeClaw:
		if ch.Class >= constant.ClassNoblesse {
			return []constant.SkillID{constant.SkillClawMasteryCygnus}
		}
		return []constant.SkillID{constant.SkillClawMastery}
	case constant.WeaponTypeKnuckle:
		if ch.Class >= constant.ClassNoblesse {
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
		if ch.Class >= constant.ClassNoblesse {
			return ch.Skills.LevelData(constant.SkillCriticalShotCygnus)
		}
		return ch.Skills.LevelData(constant.SkillCriticalShot)
	case constant.WeaponTypeClaw:
		if ch.Class >= constant.ClassNoblesse {
			return ch.Skills.LevelData(constant.SkillCriticalThrowCygnus)
		}
		return ch.Skills.LevelData(constant.SkillCriticalThrow)
	case constant.WeaponTypeKnuckle:
		return ch.Skills.LevelData(constant.SkillCriticalPunch)
	}
	return nil
}
