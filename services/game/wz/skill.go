package wz

import (
	"time"

	"github.com/boyism80/fm/types"
)

type SkillLevelData struct {
	MPCon          int
	HPCon          int
	MoneyCon       int
	ItemCon        int
	ItemConNo      int
	ItemConsume    int
	BulletConsume  int
	BulletCount    int
	Damage         int
	DamagePC       int
	FixDamage      int
	CriticalDamage int
	AttackCount    int
	MobCount       int
	PAD            int
	MAD            int
	PDD            int
	MDD            int
	EVA            int
	ACC            int
	STR            int
	HP             int
	MP             int
	Jump           int
	Speed          int
	Mastery        int
	Prop           int
	Range          int
	Time           time.Duration
	Cooldown       time.Duration
	Morph          int
	X              int
	Y              int
	Z              int
	LT             types.Vector2[int32]
	RB             types.Vector2[int32]
	HS             string
	Action         string
}

type Skill struct {
	ID           uint32
	MaxLevel     int
	MasterLevel  int
	Invisible    bool
	TimeLimited  bool
	CombatOrders bool
	ElemAttr     string
	LevelData    map[int]*SkillLevelData
}

func (s *Skill) GetLevelData(level int) *SkillLevelData {
	if s.LevelData == nil {
		return nil
	}
	if level <= 0 {
		level = 1
	}
	if level > s.MaxLevel {
		level = s.MaxLevel
	}
	if data, exists := s.LevelData[level]; exists {
		return data
	}
	maxFound := 0
	for l := range s.LevelData {
		if l <= level && l > maxFound {
			maxFound = l
		}
	}
	if maxFound > 0 {
		return s.LevelData[maxFound]
	}
	return nil
}
