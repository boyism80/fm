package constant

type Stat uint32

const (
	StatSkin        Stat = 0x1
	StatFace        Stat = 0x2
	StatHair        Stat = 0x4
	StatPet         Stat = 0x8
	StatLevel       Stat = 0x10
	StatClass       Stat = 0x20
	StatStr         Stat = 0x40
	StatDex         Stat = 0x80
	StatInt         Stat = 0x100
	StatLuk         Stat = 0x200
	StatHP          Stat = 0x400
	StatMaxHP       Stat = 0x800
	StatMP          Stat = 0x1000
	StatMaxMP       Stat = 0x2000
	StatAvailableAP Stat = 0x4000
	StatAvailableSP Stat = 0x8000
	StatEXP         Stat = 0x10000
	StatPopulation  Stat = 0x20000
	StatMeso        Stat = 0x40000
)

type StatType uint32

const (
	StatTypeStr StatType = 0x40
	StatTypeDex StatType = 0x80
	StatTypeInt StatType = 0x100
	StatTypeLuk StatType = 0x200
	StatTypeHP  StatType = 0x800
	StatTypeMP  StatType = 0x2000
)

const (
	StatMaxStrDexIntLuk uint16 = 999
	StatMaxHPMP         uint32 = 30000
	HpAPUsedMax         uint16 = 10000
)
