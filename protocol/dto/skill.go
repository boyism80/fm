package dto

type Skill struct {
	ID          uint32
	SkillLevel  uint32
	MasterLevel uint32
}

// The client reads a master level only for fourth advancement skill ids.
func (s *Skill) WritesMasterLevel() bool {
	block := s.ID / 10000
	return block%100 > 0 && block%10 == 2
}
