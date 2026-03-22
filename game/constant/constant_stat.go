package constant

// Stat represents character statistics as bit flags.
type Stat uint32

const (
	STAT_SKIN         Stat = 0x1     // Skin color
	STAT_FACE         Stat = 0x2     // Face ID
	STAT_HAIR         Stat = 0x4     // Hair ID
	STAT_PET          Stat = 0x8     // Pet data
	STAT_LEVEL        Stat = 0x10    // Character level
	STAT_CLASS        Stat = 0x20    // Job class
	STAT_STR          Stat = 0x40    // Strength
	STAT_DEX          Stat = 0x80    // Dexterity
	STAT_INT          Stat = 0x100   // Intelligence
	STAT_LUK          Stat = 0x200   // Luck
	STAT_HP           Stat = 0x400   // Current HP
	STAT_MAX_HP       Stat = 0x800   // Maximum HP
	STAT_MP           Stat = 0x1000  // Current MP
	STAT_MAX_MP       Stat = 0x2000  // Maximum MP
	STAT_AVAILABLE_AP Stat = 0x4000  // Available ability points
	STAT_AVAILABLE_SP Stat = 0x8000  // Available skill points
	STAT_EXP          Stat = 0x10000 // Experience points
	STAT_FAME         Stat = 0x20000 // Fame points
	STAT_MESO         Stat = 0x40000 // Meso currency
)

// StatType represents the stat type for AP distribution
type StatType uint32

const (
	STAT_TYPE_STR StatType = 64   // Strength
	STAT_TYPE_DEX StatType = 128  // Dexterity
	STAT_TYPE_INT StatType = 256  // Intelligence
	STAT_TYPE_LUK StatType = 512  // Luck
	STAT_TYPE_HP  StatType = 2048 // Maximum HP
	STAT_TYPE_MP  StatType = 8192 // Maximum MP
)

// Stat limits
const (
	STAT_MAX_STR_DEX_INT_LUK uint16 = 999   // Maximum value for STR, DEX, INT, LUK
	STAT_MAX_HP_MP           uint32 = 30000
	HP_AP_USED_MAX           uint16 = 10000 // Maximum HP/MP AP usage count
)
