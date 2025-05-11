package entity

type Skill struct {
	Id uint32
}

type SkillEntry struct {
	SkillLevel  int
	MasterLevel int
	Expiration  int64
}
