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

func (s *Skill) IsBeginnerSkill() bool {
	if s == nil {
		return false
	}
	jobID := int(s.ID / 10000)
	return jobID == 0 || jobID == 1 || jobID == 1000 || jobID == 2000 || jobID == 2001 || jobID == 3000 || jobID == 3001 || jobID == 2002
}

func (s *Skill) IsFourthJob() bool {
	if s == nil {
		return false
	}
	id := s.ID
	if id/10000 == 2312 {
		return true
	}
	if (s.MaxLevel <= 15 && !s.Invisible && s.MasterLevel <= 0) ||
		id == 3220010 || id == 3120011 || id == 33120010 || id == 32120009 || id == 5321006 || id == 21120011 || id == 22181004 || id == 4340010 {
		return false
	}
	block := id / 10000
	if block >= 2212 && block < 3000 {
		return (block % 10) >= 7
	}
	if block >= 430 && block <= 434 {
		return (block%10) == 4 || s.MasterLevel > 0
	}
	return (block%10) == 2 && id < 90000000 && !s.IsBeginnerSkill()
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
