package entity

type CooldownEntry struct {
	SkillId   uint32
	StartTime int64 // milliseconds
	Length    int64 // milliseconds
}
