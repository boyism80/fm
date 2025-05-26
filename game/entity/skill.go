package entity

type Skill struct {
	ID uint32
}

type SkillEntry struct {
	SkillLevel  int
	MasterLevel int
	Expiration  int64
}
