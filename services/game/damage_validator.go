package game

import (
	"fmt"
	"math"
	"strings"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
)

// attackValidator bounds client-reported damage by a maximum computed only from
// server-owned state. Validate before activation hooks consume HP, combo
// counters, ammunition or buffs: the maximum reads them.
type attackValidator struct {
	ch *Character
}

// skill returns the attacker's learned skill; id 0 is the basic attack.
func (v attackValidator) skill(id uint32) (*wz.Skill, int, error) {
	if id == 0 {
		return nil, 0, nil
	}
	entry := v.ch.Skills.Get(id)
	if entry == nil || entry.Level() <= 0 {
		return nil, 0, fmt.Errorf("skill %d not learned", id)
	}
	return entry.Wz, entry.Level(), nil
}

// start snapshots the attacker for one attack. A nil skill is the basic attack.
func (v attackValidator) start(kind attackKind, skill *wz.Skill, level int) (attack, error) {
	stats, modifiers := v.ch.damageStats()
	a := attack{kind: kind, stats: stats, modifiers: modifiers, data: &wz.SkillLevelData{}}
	if skill == nil {
		return a, nil
	}

	data := skill.GetLevelData(level)
	if data == nil {
		return attack{}, fmt.Errorf("no skill %d level %d", skill.ID, level)
	}
	a.skill = constant.SkillID(skill.ID)
	a.elemAttr = strings.ToLower(skill.ElemAttr)
	a.data = data
	return a, nil
}

func (v attackValidator) checkShape(info dto.AttackInfo, maxTargets int) error {
	if int(info.Targets) != len(info.Damages) {
		return fmt.Errorf("target count mismatch")
	}
	if len(info.Damages) > 15 || info.Hits > 15 {
		return fmt.Errorf("invalid attack shape")
	}
	if len(info.Damages) > maxTargets {
		return fmt.Errorf("too many attack targets: %d > %d", len(info.Damages), maxTargets)
	}
	for _, ap := range info.Damages {
		if len(ap.DamagePairs) != int(info.Hits) {
			return fmt.Errorf("hit count mismatch")
		}
	}
	return nil
}

func (v attackValidator) maximums(damages []dto.AttackPair, maxHits int, a attack, maximum func(attack) (uint32, error)) ([]dto.AttackPair, error) {
	bounds := make([]dto.AttackPair, len(damages))
	modifiers := a.modifiers
	a.targets = len(damages)
	seen := map[uint32]bool{}
	for ti, ap := range damages {
		if seen[ap.OID] {
			return nil, fmt.Errorf("duplicate target %d", ap.OID)
		}
		seen[ap.OID] = true
		if len(ap.DamagePairs) < 1 || len(ap.DamagePairs) > maxHits {
			return nil, fmt.Errorf("invalid hit count for %d", ap.OID)
		}

		bounds[ti] = dto.AttackPair{OID: ap.OID, DamagePairs: make([]dto.DamagePair, len(ap.DamagePairs))}
		// The client can still target a mob that died before this packet
		// arrived. damageTo skips it, so any reported damage passes.
		mob := v.ch.GetMap().GetMob(ap.OID)
		if mob == nil || !mob.IsAlive() {
			for hi := range bounds[ti].DamagePairs {
				bounds[ti].DamagePairs[hi].Damage = math.MaxUint32
			}
			continue
		}

		a.mob = mob
		a.modifiers = v.targetModifiers(modifiers, mob)
		a.target = ti
		a.hits = len(ap.DamagePairs)
		for hi := range ap.DamagePairs {
			a.hit = hi
			upper, err := maximum(a)
			if err != nil {
				return nil, err
			}
			bounds[ti].DamagePairs[hi].Damage = upper
		}
	}
	return bounds, nil
}

func (v attackValidator) checkDamage(damages, bounds []dto.AttackPair) error {
	for i, ap := range damages {
		for hi, pair := range ap.DamagePairs {
			upper := bounds[i].DamagePairs[hi].Damage
			if pair.Damage > upper {
				return fmt.Errorf("damage exceeds maximum: target=%d hit=%d actual=%d maximum=%d", ap.OID, hi, pair.Damage, upper)
			}
		}
	}
	return nil
}

func (v attackValidator) targetModifiers(modifiers damageModifiers, mob *Mob) damageModifiers {
	if mob.Buffs.GetValue(constant.MobBuffStun) > 0 {
		if p := v.ch.Skills.LevelData(constant.SkillStunMastery); p != nil {
			modifiers.CriticalChance += float64(p.Prop)
			modifiers.CriticalPercent += float64(p.Damage)
		}
	}
	if v.ch.Skills.HomingTargetOID != nil && *v.ch.Skills.HomingTargetOID == mob.GetOID() {
		if p := v.ch.Skills.LevelData(constant.SkillAdvancedHoming); p != nil {
			modifiers.HomingPercent = 100 + float64(p.X)
		}
	}
	return modifiers
}
