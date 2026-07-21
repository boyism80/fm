package wz

type Mob struct {
	ID                    uint32
	BodyAttack            int
	Level                 uint8
	MaxHP                 int
	MaxMP                 int
	Speed                 int16
	PADamage              int
	PDDamage              int
	MADamage              int
	MDDamage              int
	ACC                   int
	EVA                   int
	EXP                   uint32
	Undead                bool
	Pushed                bool
	Boss                  bool
	FfaLoot               bool
	ExplosiveReward       bool
	FS                    float32
	SummonType            uint8
	MobType               uint8
	Link                  string
	ElemResist            map[string]int
	Skills                []MobSkillSlot
	Attacks               []MobAttack
	Banish                *MobBanish
	Revives               []uint32
	RemoveAfter           int
	SelfDestructionAction int8
	HpTagColor            uint8
	HpTagBgColor          uint8
	DropItemPeriod        int
	DamagedByMob          bool
}
