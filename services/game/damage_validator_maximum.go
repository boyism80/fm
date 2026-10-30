package game

import (
	"fmt"
	"math"
	"slices"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

type attackKind uint8

const (
	attackPhysical attackKind = iota
	attackMagic
	attackSummonPhysical
	attackSummonMagic
	attackMesoExplosion
)

// The client scales charge and element percentages by float32(0.01).
const clientPercent = .009999999776482582

type damageStats struct {
	STR    float64
	DEX    float64
	INT    float64
	LUK    float64
	PAD    float64
	MAD    float64
	ACC    float64
	Level  int
	Class  constant.ClassType
	Weapon constant.WeaponType
}

// damageModifiers come from the server's buff and skill state, never from the packet.
// A zero percent means the client's identity value of 100.
type damageModifiers struct {
	Mastery             float64
	PAD                 float64
	MAD                 float64
	ACC                 float64
	Echo                float64
	EnergyPAD           float64
	ProjectilePAD       float64
	CriticalChance      float64
	CriticalPercent     float64
	SharpEyes           float64
	Amplification       float64
	ComboPercent        float64
	ElementReset        float64
	ChargeDamage        float64
	ChargeStrength      float64
	ChargeElement       string
	Berserk             float64
	HomingPercent       float64
	ShadowPercent       float64
	ExtraPercent        float64
	AssassinateSeconds  float64
	SkillDamageOverride float64
	Double              bool
	Sealed              bool
	GM                  bool
}

// attack is one hit bounded by a conservative per-hit maximum from
// server-owned state. RNG streams are not synchronized: independent endpoint
// bounds may include combinations the client's shared stream cannot produce.
// The damage reported in a packet is never an input.
type attack struct {
	kind      attackKind
	stats     damageStats
	modifiers damageModifiers
	skill     constant.SkillID
	elemAttr  string
	data      *wz.SkillLevelData
	mob       *Mob
	target    int
	targets   int
	hit       int
	hits      int
	// The attack action is not decoded, so the maximum covers every weapon
	// action. skillCritical takes the skill's own critical action instead of
	// the passive critical.
	skillCritical bool
	charge        uint32
	projectile    uint32
	meso          int32
}

type damageRange struct {
	lo float64
	hi float64
}

func newDamageRange(a, b float64) damageRange {
	return damageRange{min(a, b), max(a, b)}
}

func (r damageRange) add(o damageRange) damageRange {
	return damageRange{r.lo + o.lo, r.hi + o.hi}
}

func (r damageRange) mul(o damageRange) damageRange {
	lo := min(r.lo*o.lo, r.lo*o.hi, r.hi*o.lo, r.hi*o.hi)
	hi := max(r.lo*o.lo, r.lo*o.hi, r.hi*o.lo, r.hi*o.hi)
	return damageRange{lo, hi}
}

func (r damageRange) scale(x float64) damageRange {
	return r.mul(damageRange{x, x})
}

func (r damageRange) trunc() damageRange {
	return newDamageRange(math.Trunc(r.lo), math.Trunc(r.hi))
}

func (r damageRange) floorAt(x float64) damageRange {
	return damageRange{max(x, r.lo), max(x, r.hi)}
}

func (r damageRange) union(o damageRange) damageRange {
	return damageRange{min(r.lo, o.lo), max(r.hi, o.hi)}
}

// upper rounds the top outward to allow for wider x87 intermediates without
// a percentage damage tolerance.
func (r damageRange) upper() float64 {
	x := max(0, r.hi)
	if x == 0 {
		return 0
	}

	return math.Floor(x + max(1e-9, math.Abs(x)*1e-12))
}

func (a attack) limits() (targets, hits int) {
	if a.skill == 0 {
		return 1, 1
	}

	return max(1, a.data.MobCount), max(1, a.data.AttackCount, a.data.BulletCount)
}

func (a attack) element(d damageRange, attr string, strength float64) damageRange {
	switch a.mob.wz.ElemResist[attr] {
	case 1:
		return damageRange{}
	case 2:
		return d.scale(1 - .5*strength)
	case 3:
		return d.union(d.scale(1 + .5*strength))
	}
	return d
}

func (a attack) skillElement(d damageRange, strength float64) damageRange {
	switch a.skill {
	case constant.SkillElementComposition:
		return a.element(d, "f", 1).scale(.5).add(a.element(d, "p", 1).scale(.5))
	case constant.SkillElementComposition2211006:
		return a.element(d, "i", 1).scale(.5).add(a.element(d, "l", 1).scale(.5))
	case constant.SkillInferno, constant.SkillBlizzard3211003:
		strength = float64(a.data.X) * clientPercent
	}
	return a.element(d, a.elemAttr, strength)
}

func (a attack) physicalBase(w float64) (damageRange, error) {
	m := a.modifiers.Mastery
	str := newDamageRange(a.stats.STR*(5*m+10)*.009, a.stats.STR)
	dex := newDamageRange(a.stats.DEX*(5*m+10)*.009, a.stats.DEX)
	luk := newDamageRange(a.stats.LUK*(5*m+10)*.009, a.stats.LUK)
	formula := func(p damageRange, c, other, div float64) damageRange {
		return p.scale(c).add(damageRange{other, other}).scale(w / div)
	}

	switch a.stats.Weapon {
	case constant.WeaponTypeSword1H:
		return formula(str, 4, a.stats.DEX, 100), nil
	case constant.WeaponTypeDagger:
		if a.stats.Class.Of(constant.ClassThief) || a.stats.Class.Of(constant.ClassNightWalker1) {
			return formula(luk, 3.6, a.stats.STR+a.stats.DEX, 100), nil
		}
		return formula(str, 4, a.stats.DEX, 100), nil
	case constant.WeaponTypeAxe1H, constant.WeaponTypeBlunt1H, constant.WeaponTypeWand, constant.WeaponTypeStaff:
		return formula(str, 4.4, a.stats.DEX, 100).union(formula(str, 3.2, a.stats.DEX, 100)), nil
	case constant.WeaponTypeNone:
		c := 4.2
		if a.stats.Class == constant.ClassPirate {
			c = 3
		}
		return formula(str, c, a.stats.DEX, 100), nil
	case constant.WeaponTypeSword2H:
		return formula(str, 4.6, a.stats.DEX, 100), nil
	case constant.WeaponTypeAxe2H, constant.WeaponTypeBlunt2H:
		return formula(str, 4.8, a.stats.DEX, 100).union(formula(str, 3.4, a.stats.DEX, 100)), nil
	case constant.WeaponTypeSpear, constant.WeaponTypePoleArm:
		if a.skill == constant.SkillDragonRoar {
			return formula(str, 4, a.stats.DEX, 100), nil
		}
		return formula(str, 5, a.stats.DEX, 100).union(formula(str, 3, a.stats.DEX, 100)), nil
	case constant.WeaponTypeBow, constant.WeaponTypeCrossbow:
		if a.skill == constant.SkillPowerKnockBack || a.skill == constant.SkillPowerKnockBack3201003 {
			return formula(dex, 3.4, a.stats.STR, 150), nil
		}
		c := 3.4
		if a.stats.Weapon == constant.WeaponTypeCrossbow {
			c = 3.6
		}
		return formula(dex, c, a.stats.STR, 100).union(formula(dex, 1, a.stats.STR, 200)), nil
	case constant.WeaponTypeClaw:
		switch a.skill {
		case constant.SkillLuckySeven, constant.SkillTripleThrow, constant.SkillLuckySevenCygnus, constant.SkillTripleThrowCygnus:
			return newDamageRange(.5*a.stats.LUK, a.stats.LUK).scale(5 * w / 100), nil
		}
		other := a.stats.STR + a.stats.DEX
		return formula(luk, 1, other, 150).union(formula(luk, 3.6, other, 100)), nil
	case constant.WeaponTypeKnuckle:
		return formula(str, 4.8, a.stats.DEX, 100), nil
	case constant.WeaponTypeGun:
		return formula(dex, 3.6, a.stats.STR, 100).union(formula(dex, 1, a.stats.STR, 200)), nil
	}
	return damageRange{}, fmt.Errorf("unsupported weapon type %d", a.stats.Weapon)
}

// maximum returns a rounded-outward upper bound for one hit. Unsupported or
// malformed inputs return an error, never an unchecked uint32 maximum.
func (a attack) maximum() (uint32, error) {
	if a.targets < 1 || a.targets > 15 || a.target < 0 || a.target >= a.targets ||
		a.hits < 1 || a.hits > 15 || a.hit < 0 || a.hit >= a.hits {
		return 0, fmt.Errorf("invalid attack shape")
	}

	if a.mob.wz.OnlyNormalAttack && a.skill != 0 {
		return 0, nil
	}

	if len(a.mob.wz.DamagedBySkills) > 0 && !slices.Contains(a.mob.wz.DamagedBySkills, uint32(a.skill)) {
		return 0, nil
	}

	// Dojo HP damage precedes ordinary hit and invincibility checks.
	if a.skill == constant.SkillBambooRain || a.skill == constant.SkillBambooRainCygnus {
		return a.mob.GetMaxHp(), nil
	}

	if a.mob.wz.Invincible && !a.modifiers.GM {
		return 0, nil
	}

	// These are early returns in the client, before ordinary accuracy and buffs.
	if a.skill == constant.SkillSanctuary {
		return a.limit(1)
	}

	// Snipe's fixed damage depends on the client calculator mode, which the
	// packet does not carry.
	if a.skill == constant.SkillSnipe {
		ordinary, err := a.ordinary()
		if err != nil {
			return 0, err
		}

		return a.limit(max(600000, ordinary))
	}

	ordinary, err := a.ordinary()
	if err != nil {
		return 0, err
	}
	return a.limit(ordinary)
}

// ordinary returns the core result before target-order and template overrides.
// Shadow Partner copies this stored result, then runs the tail a second time.
func (a attack) ordinary() (float64, error) {
	magical := a.kind == attackMagic || a.kind == attackSummonMagic
	if magical && a.modifiers.Sealed {
		return 0, nil
	}
	immune := a.mob.buffs.GetValue(constant.MobBuffWeaponImmunity) > 0
	if magical {
		immune = a.mob.buffs.GetValue(constant.MobBuffMagicImmunity) > 0
	}
	immunityDamage := 0.0
	if immune {
		immunityDamage = 1
	}

	// r%100 > penetration is the immunity branch. Even penetration zero can
	// reach the ordinary branch (r%100 == 0), so immunity cannot cap it at one.
	gap := max(0, int(a.mob.wz.Level)-a.stats.Level)
	d := damageRange{}
	hit := true
	var err error
	switch {
	case a.data.FixDamage > 0:
		d = damageRange{float64(a.data.FixDamage), float64(a.data.FixDamage)}
	case a.kind == attackPhysical:
		d, hit, err = a.physical(gap)
	case magical:
		d, hit = a.magic()
	case a.kind == attackSummonPhysical:
		d = a.summon()
	case a.kind == attackMesoExplosion:
		d, err = a.mesoExplosion()
	}
	if err != nil {
		return 0, err
	}

	if !hit {
		return immunityDamage, nil
	}

	// Hit feasibility; fixed damage uses its own level-based accuracy.
	if a.kind != attackMesoExplosion && !a.modifiers.GM {
		maxHit := 1.3 * 100 * min(999, max(0, a.stats.ACC+a.modifiers.ACC)) / float64(255+10*gap)
		if magical {
			acc := 5 * (math.Trunc(a.stats.INT/10) + math.Trunc(a.stats.LUK/10))
			maxHit = 1.2 * 100 * acc / float64(255+10*gap)
		}
		if a.data.FixDamage > 0 {
			maxHit = 1.3 * 100 * float64(a.stats.Level) / float64(255+150*gap)
		}

		eva := min(999, max(0, float64(a.mob.wz.EVA)+float64(a.mob.buffs.GetValue(constant.MobBuffAvoid))))
		if maxHit < eva {
			return immunityDamage, nil
		}
	}

	if a.data.FixDamage > 0 {
		return max(immunityDamage, d.upper()), nil
	}

	if a.kind == attackPhysical && a.modifiers.ShadowPercent != 0 && a.skill != constant.SkillShowdown && a.hit >= a.hits/2 {
		original := a
		original.hit -= a.hits / 2
		value, err := original.ordinary()
		if err != nil {
			return 0, err
		}
		value = math.Trunc(min(math.MaxUint32, value) * a.modifiers.ShadowPercent / 100)
		d = damageRange{0, value}
	}

	if a.kind == attackPhysical || a.kind == attackMagic {
		if a.kind == attackPhysical && a.skill == constant.SkillAssassinate {
			d = d.scale(math.Trunc(min(a.data.Time.Seconds(), max(0, a.modifiers.AssassinateSeconds))/3) + 1)
		}
		if a.modifiers.ExtraPercent != 0 {
			d = d.scale(a.modifiers.ExtraPercent / 100)
		}
		if a.charge > 0 {
			denom := 2000.0
			if a.skill == constant.SkillCorkscrewBlow || a.skill == constant.SkillCorkscrewBlowCygnus {
				denom = 1000
			}
			d = d.scale((10 + math.Trunc(90*min(float64(a.charge), denom)/denom)) / 100)
		}

		d = d.scale((100 + a.modifiers.Berserk) / 100)
	}

	if a.modifiers.HomingPercent != 0 {
		d = d.scale(a.modifiers.HomingPercent / 100)
	}

	if a.modifiers.Double && (a.kind == attackPhysical || a.kind == attackMagic) {
		d = d.scale(2)
	}

	if a.kind == attackPhysical && (a.skill == constant.SkillSnatch || a.skill == constant.SkillSnatchCygnus) {
		switch a.hit {
		case 4:
			d = d.scale(2)
		case 5:
			d = d.scale(4)
		}
	}

	if a.kind == attackPhysical || a.kind == attackMagic {
		d = d.floorAt(1)
	}

	if immune {
		d = d.union(damageRange{1, 1})
	}

	return d.upper(), nil
}

// physical returns false when damage immunity leaves no critical to land.
func (a attack) physical(gap int) (damageRange, bool, error) {
	w := min(5000, max(0, a.stats.PAD+max(a.modifiers.PAD, a.modifiers.EnergyPAD)+a.modifiers.ProjectilePAD))
	w = min(5000, max(0, w+math.Trunc(w*a.modifiers.Echo/100)))
	d, err := a.physicalBase(w)
	if err != nil {
		return damageRange{}, false, err
	}

	if a.skill == constant.SkillShadowMeso {
		moneyCon := float64(a.data.MoneyCon)
		d = newDamageRange(.5*moneyCon, 1.5*moneyCon).trunc().scale(10)
	}

	if (a.skill == constant.SkillFlamethrower && a.projectile/1000 != 2331) ||
		(a.skill == constant.SkillIceSplitter && a.projectile/1000 != 2332) {
		d = d.scale(.5)
	}

	if !a.modifiers.GM {
		d = d.scale(float64(100-gap) / 100)
	}

	d = a.skillElement(d, 1)
	if a.modifiers.ChargeElement != "" {
		charged := d.scale(a.modifiers.ChargeDamage * clientPercent)
		d = a.element(charged, a.modifiers.ChargeElement, a.modifiers.ChargeStrength*clientPercent)
	}

	// The client skips the mob's defense for these skills, sacrifice on a boss included, whatever its description says.
	switch a.skill {
	case constant.SkillSacrifice, constant.SkillShadowMeso, constant.SkillAssaulter, constant.SkillSuperTransformation:
	default:
		pdd := min(1999, max(0, float64(a.mob.wz.PDDamage)+float64(a.mob.buffs.GetValue(constant.MobBuffWdef))))
		d = d.add(newDamageRange(-.6*pdd, -.5*pdd))
	}
	captured := d.trunc()

	skillDamage := float64(a.data.Damage)
	if a.skill == constant.SkillChargedBlow && a.modifiers.SkillDamageOverride > 0 {
		skillDamage = a.modifiers.SkillDamageOverride
	}
	if skillDamage > 0 {
		d = d.scale(skillDamage / 100)
	}

	combo := a.modifiers.ComboPercent
	if combo == 0 {
		combo = 100
	}
	d = d.scale(combo / 100)

	sharpChance := min(100, max(0, math.Floor(a.modifiers.SharpEyes/256)))
	sharpPercent := float64(uint32(a.modifiers.SharpEyes) & 255)
	chance := a.modifiers.CriticalChance + sharpChance
	percent := a.modifiers.CriticalPercent
	if percent == 0 {
		percent = 100
	}
	percent += sharpPercent
	if a.skillCritical {
		chance = float64(a.data.Prop) + sharpChance
		percent = float64(a.data.CriticalDamage) + sharpPercent
	}
	if a.skill == constant.SkillShadowMeso {
		chance = float64(a.data.Prop)
		percent = 100 + float64(a.data.X)
	}
	if a.skill == constant.SkillAssassinate {
		chance = 0
	}

	crit := d.add(captured.scale((percent - 100) / 100))
	if a.mob.buffs.GetValue(constant.MobBuffDamageImmunity) > 0 {
		if chance <= 0 {
			return damageRange{}, false, nil
		}
		d = crit
	} else if chance > 0 {
		d = d.union(crit)
	}

	defense := float64(a.mob.buffs.GetValue(constant.MobBuffWeaponDefenseUp))
	if a.skill != constant.SkillSacrifice && defense != 0 {
		d = d.scale(defense / 100)
	}

	return d, true, nil
}

// magic returns false when damage immunity leaves no critical to land.
func (a attack) magic() (damageRange, bool) {
	capMAD := 96830.0
	if a.kind == attackSummonMagic {
		capMAD = 1999
	}

	mad := min(capMAD, max(0, a.stats.MAD+a.modifiers.MAD))
	mad = min(capMAD, max(0, mad+math.Trunc(mad*a.modifiers.Echo/100)))
	mastery := float64(a.data.Mastery)
	base := math.Pow(.058*mad, 2) + .5*a.stats.INT
	d := damageRange{base, base}.add(newDamageRange(mad*(5*mastery+10)*.009, mad).scale(3.3))
	if a.kind == attackSummonMagic {
		d = d.scale(float64(a.data.MAD) / 100)
		return a.skillElement(d, 1-a.modifiers.ElementReset*clientPercent), true
	}

	const magicConstant = .011672973633297692
	amp := a.modifiers.Amplification
	if amp == 0 {
		amp = 100
	}
	if a.skill == constant.SkillHeal {
		// Heal's share shrinks as it reaches more targets, and the
		// packet does not carry how many it healed; one is the largest.
		const n = 1
		stat := newDamageRange(.2*a.stats.INT, .8*a.stats.INT).scale(1.5).add(damageRange{a.stats.LUK, a.stats.LUK})
		d = stat.scale((1 + .3*n) * float64(a.data.HP) / 100 * mad * .005 / n * amp * magicConstant)
	} else {
		d = d.scale(float64(a.data.MAD) * amp * magicConstant * magicConstant)
	}
	d = a.skillElement(d, 1-a.modifiers.ElementReset*clientPercent)

	mdd := min(1999, max(0, float64(a.mob.wz.MDDamage)+float64(a.mob.buffs.GetValue(constant.MobBuffMdef))))
	d = d.add(newDamageRange(-.6*mdd, -.5*mdd))
	if defense := float64(a.mob.buffs.GetValue(constant.MobBuffMagicDefenseUp)); defense != 0 {
		d = d.scale(defense / 100)
	}

	chance := min(100, max(0, math.Floor(a.modifiers.SharpEyes/256)))
	crit := d.scale(float64(uint32(a.modifiers.SharpEyes)&255) / 100)
	if a.mob.buffs.GetValue(constant.MobBuffDamageImmunity) > 0 {
		if chance <= 0 {
			return damageRange{}, false
		}
		d = crit
	} else if chance > 0 {
		d = d.union(crit)
	}

	return d, true
}

func (a attack) summon() damageRange {
	p, secondary := a.stats.DEX, a.stats.STR
	switch {
	case a.stats.Class.Of(constant.ClassWarrior) || a.stats.Class.Of(constant.ClassDawnWarrior1) || a.stats.Class.Of(constant.ClassAran2),
		a.stats.Class.Of(constant.ClassPirate) || a.stats.Class.Of(constant.ClassThunderBreaker1):
		p, secondary = a.stats.STR, a.stats.DEX
	case a.stats.Class.Of(constant.ClassThief) || a.stats.Class.Of(constant.ClassNightWalker1):
		p, secondary = a.stats.LUK, a.stats.DEX
	}

	d := newDamageRange(.7*p, p).scale(2.5).add(damageRange{secondary, secondary})
	d = d.scale(float64(a.data.PAD) / 100)
	return a.skillElement(d, 1)
}

func (a attack) mesoExplosion() (damageRange, error) {
	if a.meso <= 0 {
		return damageRange{}, fmt.Errorf("invalid meso amount")
	}

	amount := float64(a.meso)
	f := amount / (amount + 5250)
	if amount <= 1000 {
		f = (.82*amount + 28) / 5300
	}
	return newDamageRange(.5, 1).scale(50 * float64(a.data.X) * f), nil
}

// limit spreads one hit's maximum over its target position and clamps it to
// the packet domain only if a valid computation actually reaches it. No
// wraparound intermediate can exceed this conservative bound.
func (a attack) limit(value float64) (uint32, error) {
	value = damageRange{0, value * a.targetFactor()}.upper()
	if value > 0 && a.mob.wz.FixedDamage > 0 {
		value = float64(a.mob.wz.FixedDamage)
	}

	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("non-finite damage")
	}

	return uint32(min(math.MaxUint32, value)), nil
}

func (a attack) targetFactor() float64 {
	i := a.target
	switch a.skill {
	case constant.SkillIronArrowCrossbow:
		return []float64{
			1, .9, .81, .729, .6561, .59049, .531441, .478296,
			.430467, .387410, .348678, .313810, .282429, .254186, .228767,
		}[i]
	case constant.SkillChainLightning, constant.SkillBarrage:
		if i >= 6 {
			return 0
		}
		return []float64{1, .5, .25, .125, .0625, .03125}[i]
	case constant.SkillPiercing:
		if i >= 6 {
			return 0
		}
		return []float64{1, 1.2, 1.44, 1.728, 2.0736, 2.48832}[i]
	case constant.SkillArrowBombBow:
		if i == 0 {
			return .5
		}
		return float64(a.data.X) / 100
	}
	return 1
}
