package entity

import (
	"fmt"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type PhysicalValidator struct {
	attackValidator
}

func NewPhysicalValidator(ch *Character) *PhysicalValidator {
	return &PhysicalValidator{attackValidator{ch: ch}}
}

func (v *PhysicalValidator) Validate(info dto.AttackInfo) error {
	if v.ch.GM.InstantKill {
		return nil
	}
	bounds, err := v.maximumDamages(info)
	if err != nil {
		return err
	}
	return v.checkDamage(info.Damages, bounds)
}

func (v *PhysicalValidator) maximumDamages(info dto.AttackInfo) ([]dto.AttackPair, error) {
	if len(info.MesoOIDs) > 0 {
		return nil, fmt.Errorf("mesos on non-explosion attack")
	}
	skill, level, err := v.skill(info.Skill)
	if err != nil {
		return nil, err
	}
	a, err := v.start(attackPhysical, skill, level)
	if err != nil {
		return nil, err
	}
	maxTargets, maxHits := a.limits()

	if partner, ok := v.ch.Buffs.GetEntity(constant.BuffFlagShadowPartner).(*SkillBuff); ok && a.skill != constant.SkillShowdown {
		p := partner.Wz.GetLevelData(int(partner.SkillLevel))
		if p == nil {
			return nil, fmt.Errorf("no shadow partner level %d", partner.SkillLevel)
		}
		if a.skill == 0 {
			a.modifiers.ShadowPercent = float64(p.X)
		} else {
			a.modifiers.ShadowPercent = float64(p.Y)
		}
		maxHits *= 2
		if info.Hits%2 != 0 {
			return nil, fmt.Errorf("odd shadow partner hit count")
		}
	}

	if info.Slot > 0 {
		item := v.ch.Inventory.Tabs[constant.InventoryTypeConsume].Items[int16(info.Slot)]
		if item == nil || item.GetCount() == 0 {
			return nil, fmt.Errorf("missing projectile")
		}
		projectile, ok := item.GetModel().(*wz.Consume)
		if !ok {
			return nil, fmt.Errorf("invalid projectile %d", item.GetModel().GetID())
		}
		switch projectile.ConsumeType {
		case constant.ConsumeTypeShuriken, constant.ConsumeTypeBullet, constant.ConsumeTypeArrowBow, constant.ConsumeTypeArrowCrossBow:
			a.projectile = projectile.ID
			a.modifiers.ProjectilePAD = float64(projectile.ScrollIncPAD)
		default:
			return nil, fmt.Errorf("invalid projectile %d", item.GetModel().GetID())
		}
	}

	if a.stats.Class.ID == constant.ClassPaladin {
		if p := v.ch.skillLevelData(constant.SkillAdvancedCharge); p != nil {
			a.modifiers.SkillDamageOverride = float64(p.Damage)
		}
	}
	switch a.skill {
	case constant.SkillPanicSword, constant.SkillPanicAxe, constant.SkillComaSword, constant.SkillComaAxe, constant.SkillPanic, constant.SkillComa:
		a.modifiers.ComboPercent = v.ch.comboPercent(true)
	case constant.SkillAssassinate:
		// The server keeps no dark sight timer, so take the longest the skill allows.
		a.modifiers.AssassinateSeconds = a.data.Time.Seconds()
	}
	if v.ch.Buffs.value(constant.BuffFlagDarksight) > 0 && a.stats.Class.Family == constant.ClassFamilyCygnus {
		if p := v.ch.skillLevelData(constant.SkillVanish); p != nil {
			a.modifiers.ExtraPercent = float64(p.Damage)
		}
		if buff, ok := v.ch.Buffs.GetEntity(constant.BuffFlagDarksight).(*SkillBuff); ok && constant.SkillID(buff.Wz.ID) == constant.SkillWindWalk {
			p := buff.Wz.GetLevelData(int(buff.SkillLevel))
			if p == nil {
				return nil, fmt.Errorf("no wind walk level %d", buff.SkillLevel)
			}
			a.modifiers.ExtraPercent = float64(p.Damage)
		}
	}

	if err := v.checkShape(info, maxTargets); err != nil {
		return nil, err
	}
	a.charge = info.Charge
	return v.maximums(info.Damages, maxHits, a, v.maximum)
}

func (v *PhysicalValidator) maximum(a attack) (uint32, error) {
	upper, err := a.maximum()
	if err != nil {
		return 0, err
	}
	if a.data.Prop <= 0 || a.data.CriticalDamage <= 0 {
		return upper, nil
	}
	a.skillCritical = true
	critical, err := a.maximum()
	if err != nil {
		return 0, err
	}
	return max(upper, critical), nil
}
