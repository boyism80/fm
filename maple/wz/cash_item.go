package wz

import "github.com/boyism80/fm/maple/constant"

type CashItem struct {
	*ItemCore
	Life        int
	PetSkill    constant.PetSkill
	PetSkillAdd bool
	PetFood     *PetFood
	SoldInform  bool
	// StateChangeItem is the consume item a weather item gives everyone on the map.
	StateChangeItem uint32
}
