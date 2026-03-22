package wz

import (
	"time"

	"github.com/boyism80/fm/types"
)

// SkillLevelData contains level-specific skill data
type SkillLevelData struct {
	// Resource consumption
	MPCon    int // MP consumption amount
	HPCon    int // HP consumption amount
	MoneyCon int // Meso consumption amount

	// Item consumption
	ItemCon       int // Item consumption ID
	ItemConNo     int // Item consumption count
	ItemConsume   int // Item consume (code uses itemCon, itemConNo instead)
	BulletConsume int // Bullet consumption amount
	BulletCount   int // Bullet count (default: 1)

	// Damage and combat stats
	Damage         int // Damage (default: 100, percentage) (can be int or string, string is parsed as int)
	DamagePC       int // Damage percentage
	FixDamage      int // Fixed damage
	CriticalDamage int // Critical damage (stored in cr field)
	AttackCount    int // Attack count (default: 1) (can be int or string, string is parsed as int)
	MobCount       int // Target mob count (default: 1)

	// Defense and evasion
	PAD int // Physical attack (watk)
	MAD int // Magic attack (matk)
	PDD int // Physical defense (wdef)
	MDD int // Magic defense (mdef)
	EVA int // Evasion (avoid)
	ACC int // Accuracy (can be int or string, string is parsed as int)

	// Stat bonuses
	STR   int // STR increase (commented out in code)
	HP    int // HP absolute recovery/consumption (heal for skills, direct recovery for items)
	MP    int // MP absolute recovery/consumption
	Jump  int // Jump power
	Speed int // Movement speed

	// Skill properties
	Mastery  int           // Mastery
	Prop     int           // Probability (>= 100: always success, < 100: probability applied)
	Range    int           // Range
	Time     time.Duration // Buff duration (WZ value in seconds; .Milliseconds() for packet)
	Cooldown time.Duration // Cooldown in milliseconds

	// Special properties
	Morph int // Morph ID
	X     int // Multi-purpose variable (varies by skill: buff value, calculation value, etc.)
	Y     int // Multi-purpose variable (varies by skill: calculation value, etc.)
	Z     int // Multi-purpose variable (varies by skill)

	// Attack range (vectors)
	LT types.Vector2[int32] // Left-top coordinate (vector, parsed with getData())
	RB types.Vector2[int32] // Right-bottom coordinate (vector, parsed with getData())

	// String properties (animation/effect related)
	HS     string // Hit sound information (string)
	Action string // Action information (string, animation related)
}

// Skill represents skill data loaded from WZ files
type Skill struct {
	ID           uint32
	MaxLevel     int
	MasterLevel  int
	Invisible    bool
	TimeLimited  bool
	CombatOrders bool
	ElemAttr     string
	LevelData    map[int]*SkillLevelData // Level -> LevelData mapping
}

// GetLevelData returns the level data for a specific skill level
// If the level doesn't exist, returns the highest available level data
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

	// Try to get exact level
	if data, exists := s.LevelData[level]; exists {
		return data
	}

	// If exact level doesn't exist, find the highest available level <= requested level
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
