package wz

type Mob struct {
	ID         uint32
	BodyAttack int
	Level      uint8
	MaxHP      int
	MaxMP      int
	Speed      int16
	PADamage   int
	PDDamage   int
	MADamage   int
	MDDamage   int
	ACC        int
	EVA        int
	EXP        uint32
	Undead     bool
	Pushed     bool
	Boss       bool
	FS         float32
	SummonType uint8
	MobType    uint8
	Link       string
	ElemResist map[string]int
	Skills     []MobSkillSlot
	Attacks    []MobAttack
	Banish     *MobBanish
}
