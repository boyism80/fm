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

type MonsterCarnival struct {
	MobGenPos      []CarnivalGenPos
	GuardianGenPos []CarnivalGenPos
	Mobs           []CarnivalMobEntry
	Skills         []CarnivalSkillEntry
}

type MCSkillLevel struct {
	CP            int
	MobSkillID    uint32
	MobSkillLevel uint8
}

type MCSkill struct {
	ID     uint32
	Levels map[uint8]*MCSkillLevel
}

func (s *MCSkill) LevelData(level uint8) *MCSkillLevel {
	if s == nil || s.Levels == nil {
		return nil
	}
	return s.Levels[level]
}

type MCGuardianLevel struct {
	CP        int
	ReactorID uint32
	Duration  int
}

type MCGuardian struct {
	ID     uint32
	Levels map[uint8]*MCGuardianLevel
}

func (g *MCGuardian) LevelData(level uint8) *MCGuardianLevel {
	if g == nil || g.Levels == nil {
		return nil
	}
	return g.Levels[level]
}
