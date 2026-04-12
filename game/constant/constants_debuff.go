package constant

type DebuffFlag struct {
	Mask           uint32
	Position       int
	DiseaseSkillID uint16
}

const MaxDebuffFlagPosition = 4

var (
	DebuffFlagStun     = DebuffFlag{0x20000, 4, 123}
	DebuffFlagPoison   = DebuffFlag{0x40000, 4, 125}
	DebuffFlagSeal     = DebuffFlag{0x80000, 4, 120}
	DebuffFlagDarkness = DebuffFlag{0x100000, 4, 121}
	DebuffFlagWeaken   = DebuffFlag{0x40000000, 4, 122}
	DebuffFlagCurse    = DebuffFlag{0x80000000, 4, 124}
)

var (
	DebuffFlagSlow             = DebuffFlag{0x1, 3, 126}
	DebuffFlagSeduce           = DebuffFlag{0x80, 3, 128}
	DebuffFlagZombify          = DebuffFlag{0x4000, 3, 133}
	DebuffFlagReverseDirection = DebuffFlag{0x80000, 3, 132}
)

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
