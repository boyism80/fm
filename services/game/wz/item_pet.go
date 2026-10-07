package wz

import "github.com/boyism80/fm/services/game/constant"

type Pet struct {
	*ItemCore
	Hungry      int
	Life        int
	LimitedLife int
	NoRevive    bool
	Skills      constant.PetSkill
	Commands    []PetCommand
}

type PetCommand struct {
	Prob     int
	Inc      int
	MinLevel int
	MaxLevel int
}

type PetFood struct {
	Fullness int
	Pets     []uint32
}

func (f *PetFood) Feeds(petID uint32) bool {
	if len(f.Pets) == 0 {
		return true
	}
	for _, id := range f.Pets {
		if id == petID {
			return true
		}
	}
	return false
}
