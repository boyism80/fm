package wz

// MobHpDisplay is who sees a mob's HP while it is hit, decided once at load as the original did.
type MobHpDisplay uint8

const (
	// MobHpDisplayNone shows no one: a boss without an HP tag, unless it is a party play mob, or a mob only mobs hurt.
	MobHpDisplayNone MobHpDisplay = iota
	// MobHpDisplayBossBar shows the boss HP bar near a boss with an HP tag.
	MobHpDisplayBossBar
	// MobHpDisplayAttacker shows the attacker the mob's small HP bar.
	MobHpDisplayAttacker
)

type Mob struct {
	ID                    uint32
	BodyAttack            int
	Level                 uint8
	MaxHP                 int
	MaxMP                 int
	MPRecovery            int
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
	PartyBonus            bool
	NoDoom                bool
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
	Buff                  uint32
	RemoveAfter           int
	HasSelfDestruction    bool
	SelfDestructionAction int8
	SelfDestructionHP     int
	HpTagColor            uint8
	HpTagBgColor          uint8
	HpDisplay             MobHpDisplay
	DropItemPeriod        int
	DamagedByMob          bool
	Invincible            bool
	FirstAttack           bool
	FixedDamage           int
	OnlyNormalAttack      bool
	DamagedBySkills       []uint32
	CP                    int
	Point                 int
}
