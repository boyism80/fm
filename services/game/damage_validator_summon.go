package game

import (
	"fmt"

	"github.com/boyism80/fm/protocol/dto"
)

type SummonValidator struct {
	attackValidator
	summon *Summon
}

func NewSummonValidator(ch *Character, summon *Summon) *SummonValidator {
	return &SummonValidator{attackValidator: attackValidator{ch: ch}, summon: summon}
}

func (v *SummonValidator) Validate(damages []dto.AttackPair) error {
	if v.ch.gm.instantKill {
		return nil
	}
	bounds, err := v.maximumDamages(damages)
	if err != nil {
		return err
	}
	return v.checkDamage(damages, bounds)
}

func (v *SummonValidator) maximumDamages(damages []dto.AttackPair) ([]dto.AttackPair, error) {
	if v.ch.gameWorld == nil {
		return nil, fmt.Errorf("no game world")
	}
	skill := v.ch.gameWorld.GetResources().GetSkill(uint32(v.summon.skillID))
	if skill == nil {
		return nil, fmt.Errorf("no summon skill %d", v.summon.skillID)
	}
	a, err := v.start(attackSummonPhysical, skill, int(v.summon.skillLevel))
	if err != nil {
		return nil, err
	}
	if a.data.MAD > 0 {
		a.kind = attackSummonMagic
	}
	maxTargets, _ := a.limits()
	info := dto.AttackInfo{AttackHeader: dto.AttackHeader{Targets: uint8(len(damages)), Hits: 1}, Damages: damages}
	if err := v.checkShape(info, maxTargets); err != nil {
		return nil, err
	}

	return v.maximums(damages, 1, a, attack.maximum)
}
