package constant

// DebuffFlag identifies a character debuff (disease) type for packets.
// Mask and Position must match MapleDisease (getValue, getPosition); position 3 or 4.
type DebuffFlag struct {
	Mask     uint32
	Position int
}

const MaxDebuffFlagPosition = 4

// Position 4 (Java MapleDisease first=4).
var (
	DebuffFlagStun     = DebuffFlag{0x20000, 4}
	DebuffFlagPoison   = DebuffFlag{0x40000, 4}
	DebuffFlagSeal     = DebuffFlag{0x80000, 4}
	DebuffFlagDarkness = DebuffFlag{0x100000, 4}
	DebuffFlagWeaken   = DebuffFlag{0x40000000, 4}
	DebuffFlagCurse    = DebuffFlag{0x80000000, 4}
)

// Position 3 (Java MapleDisease first=3).
var (
	DebuffFlagSlow             = DebuffFlag{0x1, 3}
	DebuffFlagSeduce           = DebuffFlag{0x80, 3}
	DebuffFlagZombify          = DebuffFlag{0x4000, 3}
	DebuffFlagReverseDirection = DebuffFlag{0x80000, 3}
)

// AllDebuffFlags returns all DebuffFlag constants for Lua injection.
func AllDebuffFlags() map[string]DebuffFlag {
	return map[string]DebuffFlag{
		"Stun":             DebuffFlagStun,
		"Poison":           DebuffFlagPoison,
		"Seal":             DebuffFlagSeal,
		"Darkness":         DebuffFlagDarkness,
		"Weaken":           DebuffFlagWeaken,
		"Curse":            DebuffFlagCurse,
		"Slow":             DebuffFlagSlow,
		"Seduce":           DebuffFlagSeduce,
		"Zombify":          DebuffFlagZombify,
		"ReverseDirection": DebuffFlagReverseDirection,
	}
}
