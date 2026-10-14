package domain

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
)

// MesoExplosionValidator checks meso explosion, whose hits are the exploded
// mesos rather than the skill's attack count.
type MesoExplosionValidator struct {
	attackValidator
}

func NewMesoExplosionValidator(ch *Character) *MesoExplosionValidator {
	return &MesoExplosionValidator{attackValidator{ch: ch}}
}

func (v *MesoExplosionValidator) Validate(info dto.AttackInfo) error {
	if v.ch.GM.InstantKill {
		return nil
	}
	bounds, err := v.maximumDamages(info)
	if err != nil {
		return err
	}

	// The client does not say which meso made which hit: pair the largest
	// damage with the largest meso.
	damages := make([]dto.AttackPair, len(info.Damages))
	for i, ap := range info.Damages {
		pairs := slices.Clone(ap.DamagePairs)
		slices.SortFunc(pairs, func(a, b dto.DamagePair) int {
			return cmp.Compare(b.Damage, a.Damage)
		})
		damages[i] = dto.AttackPair{OID: ap.OID, DamagePairs: pairs}
	}
	return v.checkDamage(damages, bounds)
}

// maximumDamages returns bounds ordered by descending meso amount for each target.
// Validate matches these with the reported hits in descending damage order.
func (v *MesoExplosionValidator) maximumDamages(info dto.AttackInfo) ([]dto.AttackPair, error) {
	if constant.SkillID(info.Skill) != constant.SkillMesoExplosion {
		return nil, fmt.Errorf("not a meso explosion skill")
	}
	if len(info.MesoOIDs) > 15 || len(info.Damages) > 15 {
		return nil, fmt.Errorf("invalid meso explosion shape")
	}
	skill, level, err := v.skill(info.Skill)
	if err != nil {
		return nil, err
	}
	a, err := v.start(attackMesoExplosion, skill, level)
	if err != nil {
		return nil, err
	}
	amounts, err := v.amounts(info.MesoOIDs)
	if err != nil {
		return nil, err
	}

	maxTargets, _ := a.limits()
	if int(info.Targets) != len(info.Damages) {
		return nil, fmt.Errorf("target count mismatch")
	}
	if len(info.Damages) > maxTargets {
		return nil, fmt.Errorf("too many attack targets: %d > %d", len(info.Damages), maxTargets)
	}

	return v.maximums(info.Damages, len(amounts), a, func(a attack) (uint32, error) {
		a.meso = amounts[a.hit]
		return a.maximum()
	})
}

// amounts returns the exploded meso amounts, largest first.
func (v *MesoExplosionValidator) amounts(oids []uint32) ([]int32, error) {
	area, ok := v.ch.mesoExplosionArea()
	if !ok {
		return nil, fmt.Errorf("missing meso explosion level")
	}

	items := v.ch.GetMap().GetItems()
	seen := map[uint32]bool{}
	amounts := make([]int32, 0, len(oids))
	for _, oid := range oids {
		meso, ok := items[oid].(*Meso)
		if !ok || meso.Count <= 0 || seen[oid] || !meso.inArea(area) {
			return nil, fmt.Errorf("invalid explosion meso %d", oid)
		}
		seen[oid] = true
		amounts = append(amounts, meso.Count)
	}
	slices.Sort(amounts)
	slices.Reverse(amounts)
	return amounts, nil
}
