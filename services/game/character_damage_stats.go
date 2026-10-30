package game

import (
	"math"

	"github.com/boyism80/fm/maple/constant"
)

// comboPercent implements the client's stored-count minus one convention.
func (ch *Character) comboPercent(finisher bool) float64 {
	baseID, advancedID, capLevel := constant.SkillComboAttack, constant.SkillAdvancedCombo, 30
	if ch.class >= constant.ClassNoblesse {
		baseID, advancedID, capLevel = constant.SkillComboAttackCygnus, constant.SkillAdvancedComboCygnus, 20
	}
	skill := ch.skills.LevelData(baseID)
	if skill == nil {
		return 100
	}
	level := ch.skills.Level(uint32(baseID))
	if level >= capLevel {
		if advanced := ch.skills.LevelData(advancedID); advanced != nil {
			skill = advanced
		}
	}

	c := int(ch.buffs.value(constant.BuffFlagCombo)) - 1
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
	// Include all equipped bonuses as an envelope when client eligibility is
	// unavailable; ineligible pieces can only increase this nonnegative bound.
	s := ch.totalStats()

	mapleWarrior := ch.buffs.value(constant.BuffFlagMapleWarrior)
	s.STR += math.Trunc(float64(ch.stats.base.str) * mapleWarrior / 100)
	s.DEX += math.Trunc(float64(ch.stats.base.dex) * mapleWarrior / 100)
	s.INT += math.Trunc(float64(ch.stats.base.int) * mapleWarrior / 100)
	s.LUK += math.Trunc(float64(ch.stats.base.luk) * mapleWarrior / 100)
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
		if p := ch.skills.LevelData(id); p != nil {
			s.ACC += float64(p.X)
		}
	}
	for _, id := range []constant.SkillID{constant.SkillBlessingOfTheFairy, constant.SkillBlessingOfTheFairyCygnus} {
		if p := ch.skills.LevelData(id); p != nil {
			s.PAD += float64(p.X)
			s.MAD += float64(p.Y)
			s.ACC += float64(p.Z)
		}
	}

	b := damageModifiers{
		PAD:          ch.buffs.value(constant.BuffFlagWeaponAtk),
		MAD:          ch.buffs.value(constant.BuffFlagMagicAtk),
		ACC:          ch.buffs.value(constant.BuffFlagAcc),
		Echo:         ch.buffs.value(constant.BuffFlagEchoOfHero),
		SharpEyes:    ch.buffs.value(constant.BuffFlagSharpEyes),
		ElementReset: ch.buffs.value(constant.BuffFlagElementReset),
		Double:       ch.buffs.value(constant.BuffFlagBerserkFury) > 0,
		GM:           s.Class.Of(constant.ClassGM) || s.Class >= 28 && s.Class <= 30,
	}
	if ch.debuffs != nil {
		b.Sealed = ch.debuffs.Has(constant.DebuffFlagSeal)
	}

	for _, id := range ch.skills.masteries(s.Weapon) {
		if p := ch.skills.LevelData(id); p != nil {
			b.Mastery = max(b.Mastery, float64(p.Mastery))
			s.PAD += float64(p.X)
			s.ACC += float64(p.X)
		}
	}
	if s.Weapon == constant.WeaponTypeSpear || s.Weapon == constant.WeaponTypePoleArm {
		if p := ch.skills.LevelData(constant.SkillBeholder); p != nil {
			b.Mastery += float64(p.Mastery)
		}
	}

	b.CriticalPercent = 100
	if critical := ch.skills.criticalLevelData(s.Weapon); critical != nil {
		b.CriticalChance = float64(critical.Prop)
		b.CriticalPercent = float64(critical.Damage)
	}

	// Energy charge's PAD is active only in the charged state. The server buff
	// value is authoritative; include its source PAD for an upper envelope.
	if ch.buffs.value(constant.BuffFlagEnergyCharge) > 0 {
		energyCharge := constant.SkillEnergyCharge
		if s.Class >= constant.ClassNoblesse {
			energyCharge = constant.SkillEnergyChargeCygnus
		}
		if p := ch.skills.LevelData(energyCharge); p != nil {
			b.EnergyPAD = float64(p.PAD)
		}
	}

	amplificationSkills := []constant.SkillID{
		constant.SkillElementAmplification,
		constant.SkillElementAmplification2210001,
		constant.SkillElementAmplificationCygnus,
	}
	for _, id := range amplificationSkills {
		if p := ch.skills.LevelData(id); p != nil {
			b.Amplification = max(b.Amplification, float64(p.Y))
		}
	}

	if ch.buffs.value(constant.BuffFlagCombo) > 0 {
		b.ComboPercent = ch.comboPercent(false)
	}

	if charge, ok := ch.buffs.GetEntity(constant.BuffFlagWkCharge).(*SkillBuff); ok {
		if p := charge.wz.GetLevelData(int(charge.skillLevel)); p != nil {
			b.ChargeDamage = float64(p.Damage)
			b.ChargeStrength = float64(p.Z)
		}
		switch constant.SkillID(charge.wz.ID) {
		case constant.SkillFireChargeSword, constant.SkillFireChargeBw:
			b.ChargeElement = "f"
		case constant.SkillIceChargeSword, constant.SkillBlizzardChargeBw:
			b.ChargeElement = "i"
		case constant.SkillThunderChargeSword, constant.SkillLightningChargeBw, constant.SkillLightningCharge:
			b.ChargeElement = "l"
		case constant.SkillHolyChargeSword, constant.SkillDivineChargeMace, constant.SkillRadiantCharge:
			b.ChargeElement = "h"
		}
	}

	if berserk := ch.skills.LevelData(constant.SkillBerserk); s.Class == constant.ClassDarkKnight && berserk != nil {
		if float64(ch.GetHp()) <= math.Trunc(float64(ch.GetMaxHp())*float64(berserk.X)/100) {
			b.Berserk = float64(berserk.Damage)
		}
	}

	return s, b
}
