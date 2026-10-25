package wz

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
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

// WZ lt/rb describe the area while facing left; facing right mirrors it. Range widens the front.
func (ld *SkillLevelData) Area(x, y int32, facingLeft bool) types.Rect[int32] {
	r := int32(ld.Range)
	if ld.LT == (types.Vector2[int32]{}) && ld.RB == (types.Vector2[int32]{}) {
		if facingLeft {
			return types.Rect[int32]{Left: x - 200 - r, Top: y - 100 - r, Right: x, Bottom: y + 100}
		}
		return types.Rect[int32]{Left: x, Top: y - 100 - r, Right: x + 200 + r, Bottom: y + 100}
	}
	if facingLeft {
		return types.Rect[int32]{Left: x + ld.LT.X - r, Top: y + ld.LT.Y, Right: x + ld.RB.X, Bottom: y + ld.RB.Y}
	}
	return types.Rect[int32]{Left: x - ld.RB.X, Top: y + ld.LT.Y, Right: x - ld.LT.X + r, Bottom: y + ld.RB.Y}
}

type Skill struct {
	ID           uint32
	Class        constant.ClassType
	MaxLevel     int
	MasterLevel  int
	Invisible    bool
	TimeLimited  bool
	CombatOrders bool
	Charge       bool
	ElemAttr     string
	MobSkill     bool
	LevelData    map[int]*SkillLevelData

	ManualCooldown bool
	ExtraElement   string
	Dojo           bool
	FullHolySymbol bool
	FollowUpAttack bool

	Beginner      bool
	FourthClass   bool
	Distributable bool
	Attack        bool
}

func (s *Skill) DefaultMasterLevel() int {
	if s.MasterLevel > 0 {
		return s.MasterLevel
	}
	if s.FourthClass {
		return 0
	}
	return s.MaxLevel
}

func (s *Skill) LearnableBy(class constant.ClassType) bool {
	if s.Class == constant.ClassBeginner {
		return class < constant.ClassNoblesse
	}
	return class.Of(s.Class)
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
