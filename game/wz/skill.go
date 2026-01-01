package wz

// SkillLevelData contains level-specific skill data
type SkillLevelData struct {
	MPCon    int // MP consumption
	Cooldown int // Cooldown in seconds
	Damage   int // Damage percentage
	HPCon    int // HP consumption
}

// Skill represents skill data loaded from WZ files
type Skill struct {
	ID           uint32
	MaxLevel     int
	MasterLevel  int
	Invisible    bool
	TimeLimited  bool
	CombatOrders bool
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
