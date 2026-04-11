package dto

// Skill represents skill data for protocol
type Skill struct {
	ID          uint32
	SkillLevel  uint32
	MasterLevel uint32 // Only for certain skill types
}
