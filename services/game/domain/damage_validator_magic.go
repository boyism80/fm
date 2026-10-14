package domain

import (
	"fmt"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
)

type MagicValidator struct {
	attackValidator
}

func NewMagicValidator(ch *Character) *MagicValidator {
	return &MagicValidator{attackValidator{ch: ch}}
}

func (v *MagicValidator) Validate(info dto.AttackInfo) error {
	if v.ch.GM.InstantKill {
		return nil
	}
	bounds, err := v.maximumDamages(info)
	if err != nil {
		return err
	}
	return v.checkDamage(info.Damages, bounds)
}

func (v *MagicValidator) maximumDamages(info dto.AttackInfo) ([]dto.AttackPair, error) {
	if info.Skill == 0 {
		return nil, fmt.Errorf("magic attack requires a skill")
	}
	if len(info.MesoOIDs) > 0 {
		return nil, fmt.Errorf("mesos on magic attack")
	}
	skill, level, err := v.skill(info.Skill)
	if err != nil {
		return nil, err
	}
	a, err := v.start(attackMagic, skill, level)
	if err != nil {
		return nil, err
	}
	maxTargets, maxHits := a.limits()
	if err := v.checkShape(info, maxTargets); err != nil {
		return nil, err
	}

	if constant.SkillID(info.Skill) == constant.SkillHeal {
		for _, ap := range info.Damages {
			mob := v.ch.GetMap().GetMob(ap.OID)
			if mob != nil && mob.IsAlive() && !mob.Wz.Undead {
				return nil, fmt.Errorf("heal hit non-undead mob %d", ap.OID)
			}
		}
	}

	a.charge = info.Charge
	return v.maximums(info.Damages, maxHits, a, attack.maximum)
}
