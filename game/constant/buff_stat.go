// Package constant defines buff types for temporary buffs.
// Mask is the client bit flag; Position is the mask group index (1-based, 2–4).
// MaxBuffFlag = 4: packet uses 4 int32 mask words.
package constant

const MaxBuffFlag = 4

// BuffFlag identifies a buff type for packets (mask value + position).
type BuffFlag struct {
	Mask     uint32
	Position int
}

// Position 4.
var (
	BuffFlagWeaponAtk     = BuffFlag{0x1, 4}
	BuffFlagWeaponDef     = BuffFlag{0x2, 4}
	BuffFlagMagicAtk      = BuffFlag{0x4, 4}
	BuffFlagMagicDef      = BuffFlag{0x8, 4}
	BuffFlagAcc           = BuffFlag{0x10, 4}
	BuffFlagAvoid         = BuffFlag{0x20, 4}
	BuffFlagHands         = BuffFlag{0x40, 4}
	BuffFlagSpeed         = BuffFlag{0x80, 4}
	BuffFlagJump          = BuffFlag{0x100, 4}
	BuffFlagMagicGuard    = BuffFlag{0x200, 4}
	BuffFlagDarksight     = BuffFlag{0x400, 4}
	BuffFlagBooster       = BuffFlag{0x800, 4}
	BuffFlagPowerguard    = BuffFlag{0x1000, 4}
	BuffFlagMaxHp         = BuffFlag{0x2000, 4}
	BuffFlagMaxMp         = BuffFlag{0x4000, 4}
	BuffFlagInvincible    = BuffFlag{0x8000, 4}
	BuffFlagSoulArrow     = BuffFlag{0x10000, 4}
	BuffFlagCombo         = BuffFlag{0x200000, 4}
	BuffFlagSummon        = BuffFlag{0x200000, 4}
	BuffFlagWkCharge      = BuffFlag{0x400000, 4}
	BuffFlagDragonBlood   = BuffFlag{0x800000, 4}
	BuffFlagHolySymbol    = BuffFlag{0x1000000, 4}
	BuffFlagMesoUp        = BuffFlag{0x2000000, 4}
	BuffFlagShadowPartner = BuffFlag{0x4000000, 4}
	BuffFlagPickpocket    = BuffFlag{0x8000000, 4}
	BuffFlagPuppet        = BuffFlag{0x8000000, 4}
	BuffFlagMesoGuard     = BuffFlag{0x10000000, 4}
	BuffFlagHpLossGuard   = BuffFlag{0x20000000, 4}
)

// Position 3.
var (
	BuffFlagMorph          = BuffFlag{0x2, 3}
	BuffFlagRecovery       = BuffFlag{0x4, 3}
	BuffFlagMapleWarrior   = BuffFlag{0x8, 3}
	BuffFlagStance         = BuffFlag{0x10, 3}
	BuffFlagSharpEyes      = BuffFlag{0x20, 3}
	BuffFlagManaReflection = BuffFlag{0x40, 3}
	BuffFlagSpiritClaw     = BuffFlag{0x100, 3}
	BuffFlagInfinity       = BuffFlag{0x200, 3}
	BuffFlagHolyShield     = BuffFlag{0x400, 3}
	BuffFlagHamstring      = BuffFlag{0x800, 3}
	BuffFlagBlind          = BuffFlag{0x1000, 3}
	BuffFlagConcentrate    = BuffFlag{0x2000, 3}
	BuffFlagEchoOfHero     = BuffFlag{0x8000, 3}
	BuffFlagUnknown3       = BuffFlag{0x10000, 3}
	BuffFlagGhostMorph     = BuffFlag{0x20000, 3}
	BuffFlagAriantCossImu  = BuffFlag{0x40000, 3}
	BuffFlagReverseDir     = BuffFlag{0x80000, 3}
	BuffFlagUnknown6       = BuffFlag{0x400000, 3}
	BuffFlagUnknown7       = BuffFlag{0x800000, 3}
	BuffFlagUnknown8       = BuffFlag{0x1000000, 3}
	BuffFlagBerserkFury    = BuffFlag{0x2000000, 3}
	BuffFlagDivineBody     = BuffFlag{0x4000000, 3}
	BuffFlagSpark          = BuffFlag{0x8000000, 3}
	BuffFlagAriantCossImu2 = BuffFlag{0x10000000, 3}
	BuffFlagFinalAttack    = BuffFlag{0x20000000, 3}
	BuffFlagElementReset   = BuffFlag{0x80000000, 3}
)

// Position 2.
var (
	BuffFlagEnergyCharge  = BuffFlag{0x2, 2}
	BuffFlagDashSpeed     = BuffFlag{0x4, 2}
	BuffFlagDashJump      = BuffFlag{0x8, 2}
	BuffFlagMonsterRiding = BuffFlag{0x10, 2}
	BuffFlagWindBooster   = BuffFlag{0x20, 2}
	BuffFlagHomingBeacon  = BuffFlag{0x40, 2}
	BuffFlagExpRate       = BuffFlag{0x400, 2}
	BuffFlagDropRate      = BuffFlag{0x800, 2}
	BuffFlagMesoRate      = BuffFlag{0x1000, 2}
	BuffFlagSoaring       = BuffFlag{0x40000, 2}
)

// IsRemoteStatFlag reports whether a buff flag is visible to other players in spawn packets.
// This mirrors Java's TemporaryStatsPacket.isForRemoteStat.
func IsRemoteStatFlag(flag BuffFlag) bool {
	return flag == BuffFlagSpeed ||
		flag == BuffFlagCombo ||
		flag == BuffFlagWkCharge ||
		flag == BuffFlagShadowPartner ||
		flag == BuffFlagDarksight ||
		flag == BuffFlagSoulArrow ||
		flag == BuffFlagMorph ||
		flag == BuffFlagSpiritClaw ||
		flag == BuffFlagBerserkFury ||
		flag == BuffFlagDivineBody
}

// AllBuffFlags returns all BuffFlag constants by name for injection into Lua (e.g. BuffFlag table).
func AllBuffFlags() map[string]BuffFlag {
	return map[string]BuffFlag{
		"WeaponAtk":      BuffFlagWeaponAtk,
		"WeaponDef":      BuffFlagWeaponDef,
		"MagicAtk":       BuffFlagMagicAtk,
		"MagicDef":       BuffFlagMagicDef,
		"Acc":            BuffFlagAcc,
		"Avoid":          BuffFlagAvoid,
		"Hands":          BuffFlagHands,
		"Speed":          BuffFlagSpeed,
		"Jump":           BuffFlagJump,
		"MagicGuard":     BuffFlagMagicGuard,
		"Darksight":      BuffFlagDarksight,
		"Booster":        BuffFlagBooster,
		"Powerguard":     BuffFlagPowerguard,
		"MaxHp":          BuffFlagMaxHp,
		"MaxMp":          BuffFlagMaxMp,
		"Invincible":     BuffFlagInvincible,
		"SoulArrow":      BuffFlagSoulArrow,
		"Combo":          BuffFlagCombo,
		"Summon":         BuffFlagSummon,
		"WkCharge":       BuffFlagWkCharge,
		"DragonBlood":    BuffFlagDragonBlood,
		"HolySymbol":     BuffFlagHolySymbol,
		"MesoUp":         BuffFlagMesoUp,
		"ShadowPartner":  BuffFlagShadowPartner,
		"Pickpocket":     BuffFlagPickpocket,
		"Puppet":         BuffFlagPuppet,
		"MesoGuard":      BuffFlagMesoGuard,
		"HpLossGuard":    BuffFlagHpLossGuard,
		"Morph":          BuffFlagMorph,
		"Recovery":       BuffFlagRecovery,
		"MapleWarrior":   BuffFlagMapleWarrior,
		"Stance":         BuffFlagStance,
		"SharpEyes":      BuffFlagSharpEyes,
		"ManaReflection": BuffFlagManaReflection,
		"SpiritClaw":     BuffFlagSpiritClaw,
		"Infinity":       BuffFlagInfinity,
		"HolyShield":     BuffFlagHolyShield,
		"Hamstring":      BuffFlagHamstring,
		"Blind":          BuffFlagBlind,
		"Concentrate":    BuffFlagConcentrate,
		"EchoOfHero":     BuffFlagEchoOfHero,
		"Unknown3":       BuffFlagUnknown3,
		"GhostMorph":     BuffFlagGhostMorph,
		"AriantCossImu":  BuffFlagAriantCossImu,
		"ReverseDir":     BuffFlagReverseDir,
		"Unknown6":       BuffFlagUnknown6,
		"Unknown7":       BuffFlagUnknown7,
		"Unknown8":       BuffFlagUnknown8,
		"BerserkFury":    BuffFlagBerserkFury,
		"DivineBody":     BuffFlagDivineBody,
		"Spark":          BuffFlagSpark,
		"AriantCossImu2": BuffFlagAriantCossImu2,
		"FinalAttack":    BuffFlagFinalAttack,
		"ElementReset":   BuffFlagElementReset,
		"EnergyCharge":   BuffFlagEnergyCharge,
		"DashSpeed":      BuffFlagDashSpeed,
		"DashJump":       BuffFlagDashJump,
		"MonsterRiding":  BuffFlagMonsterRiding,
		"Ridding":        BuffFlagMonsterRiding,
		"WindBooster":    BuffFlagWindBooster,
		"HomingBeacon":   BuffFlagHomingBeacon,
		"ExpRate":        BuffFlagExpRate,
		"DropRate":       BuffFlagDropRate,
		"MesoRate":       BuffFlagMesoRate,
		"Soaring":        BuffFlagSoaring,
	}
}
