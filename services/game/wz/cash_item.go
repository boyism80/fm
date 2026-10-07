package wz

import "github.com/boyism80/fm/services/game/constant"

type CashItem struct {
	*ItemCore
	Life        int
	PetSkill    constant.PetSkill
	PetSkillAdd bool
	PetFood     *PetFood
}
