package constant

type BuffType uint8

const (
	BuffTypeAll   BuffType = 0
	BuffTypeItem  BuffType = 1
	BuffTypeSkill BuffType = 2
)

func AllBuffTypes() map[string]BuffType {
	return map[string]BuffType{
		"All":   BuffTypeAll,
		"Item":  BuffTypeItem,
		"Skill": BuffTypeSkill,
	}
}
