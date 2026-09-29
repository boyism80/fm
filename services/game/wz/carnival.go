package wz

import "github.com/boyism80/fm/types"

type CarnivalMobEntry struct {
	ID      uint32
	SpendCP int
}

type CarnivalSkillEntry struct {
	ID      uint32
	SpendCP int
}

type CarnivalGenPos struct {
	Pos  types.Point[int16]
	Team int
}

type CarnivalField struct {
	MobGenPos      []CarnivalGenPos
	GuardianGenPos []CarnivalGenPos
	Mobs           []CarnivalMobEntry
	Skills         []CarnivalSkillEntry
}

type CarnivalSkill struct {
	ID         uint32
	SpendCP    int
	MobSkillID uint32
	Level      uint8
	TargetsAll bool
	HitChance  int
}

type CarnivalGuardian struct {
	ID         uint32
	SpendCP    int
	MobSkillID uint32
	Level      uint8
}
