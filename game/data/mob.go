package data

type MobSpec struct {
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
	EXP        int
	Undead     bool
	Pushed     bool
	FS         float32
	SummonType uint8
	MobType    uint8
}
