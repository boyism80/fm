package constant

type Stat uint32

const (
	STAT_SKIN         Stat = 0x1
	STAT_FACE         Stat = 0x2
	STAT_HAIR         Stat = 0x4
	STAT_PET          Stat = 0x8
	STAT_LEVEL        Stat = 0x10
	STAT_CLASS        Stat = 0x20
	STAT_STR          Stat = 0x40
	STAT_DEX          Stat = 0x80
	STAT_INT          Stat = 0x100
	STAT_LUK          Stat = 0x200
	STAT_HP           Stat = 0x400
	STAT_MAX_HP       Stat = 0x800
	STAT_MP           Stat = 0x1000
	STAT_MAX_MP       Stat = 0x2000
	STAT_AVAILABLE_AP Stat = 0x4000
	STAT_AVAILABLE_SP Stat = 0x8000
	STAT_EXP          Stat = 0x10000
	STAT_FAME         Stat = 0x20000
	STAT_MESO         Stat = 0x40000
)

type StatType uint32

const (
	STAT_TYPE_STR StatType = 0x40
	STAT_TYPE_DEX StatType = 0x80
	STAT_TYPE_INT StatType = 0x100
	STAT_TYPE_LUK StatType = 0x200
	STAT_TYPE_HP  StatType = 0x800
	STAT_TYPE_MP  StatType = 0x2000
)

const (
	STAT_MAX_STR_DEX_INT_LUK uint16 = 999
	STAT_MAX_HP_MP           uint32 = 30000
	HP_AP_USED_MAX           uint16 = 10000
)
