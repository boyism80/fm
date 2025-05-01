package entity

type Skill struct {
	Id uint32
	// 필요한 만큼 필드 추가 가능
}

type SkillEntry struct {
	SkillLevel  int
	MasterLevel int
	Expiration  int64
}
