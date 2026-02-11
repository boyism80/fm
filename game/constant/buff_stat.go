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

// Common buff flags (position 4 = first mask group).
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
	BuffFlagWkCharge      = BuffFlag{0x400000, 4}
	BuffFlagDragonBlood   = BuffFlag{0x800000, 4}
	BuffFlagHolySymbol    = BuffFlag{0x1000000, 4}
	BuffFlagMesoUp        = BuffFlag{0x2000000, 4}
	BuffFlagShadowPartner = BuffFlag{0x4000000, 4}
	BuffFlagPickpocket    = BuffFlag{0x8000000, 4}
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
	BuffFlagBerserkFury    = BuffFlag{0x2000000, 3}
	BuffFlagDivineBody     = BuffFlag{0x4000000, 3}
	BuffFlagFinalAttack    = BuffFlag{0x20000000, 3}
)

// Position 2.
var (
	BuffFlagEnergyCharge  = BuffFlag{0x2, 2}
	BuffFlagDashSpeed     = BuffFlag{0x4, 2}
	BuffFlagDashJump      = BuffFlag{0x8, 2}
	BuffFlagMonsterRiding = BuffFlag{0x10, 2}
	BuffFlagSpeedInfusion = BuffFlag{0x20, 2}
	BuffFlagHomingBeacon  = BuffFlag{0x40, 2}
	BuffFlagExpRate       = BuffFlag{0x400, 2}
	BuffFlagDropRate      = BuffFlag{0x800, 2}
	BuffFlagMesoRate      = BuffFlag{0x1000, 2}
)

// AllBuffFlags returns all BuffFlag constants by name for injection into Lua (e.g. BuffFlag table).
func AllBuffFlags() map[string]BuffFlag {
	return map[string]BuffFlag{
		"WeaponAtk": BuffFlagWeaponAtk, "WeaponDef": BuffFlagWeaponDef, "MagicAtk": BuffFlagMagicAtk, "MagicDef": BuffFlagMagicDef,
		"Acc": BuffFlagAcc, "Avoid": BuffFlagAvoid, "Hands": BuffFlagHands, "Speed": BuffFlagSpeed, "Jump": BuffFlagJump,
		"MagicGuard": BuffFlagMagicGuard, "Darksight": BuffFlagDarksight, "Booster": BuffFlagBooster, "Powerguard": BuffFlagPowerguard,
		"MaxHp": BuffFlagMaxHp, "MaxMp": BuffFlagMaxMp, "Invincible": BuffFlagInvincible, "SoulArrow": BuffFlagSoulArrow,
		"Combo": BuffFlagCombo, "WkCharge": BuffFlagWkCharge, "DragonBlood": BuffFlagDragonBlood, "HolySymbol": BuffFlagHolySymbol,
		"MesoUp": BuffFlagMesoUp, "ShadowPartner": BuffFlagShadowPartner, "Pickpocket": BuffFlagPickpocket, "MesoGuard": BuffFlagMesoGuard, "HpLossGuard": BuffFlagHpLossGuard,
		"Morph": BuffFlagMorph, "Recovery": BuffFlagRecovery, "MapleWarrior": BuffFlagMapleWarrior, "Stance": BuffFlagStance,
		"SharpEyes": BuffFlagSharpEyes, "ManaReflection": BuffFlagManaReflection, "SpiritClaw": BuffFlagSpiritClaw, "Infinity": BuffFlagInfinity,
		"BerserkFury": BuffFlagBerserkFury, "DivineBody": BuffFlagDivineBody, "FinalAttack": BuffFlagFinalAttack,
		"EnergyCharge": BuffFlagEnergyCharge, "DashSpeed": BuffFlagDashSpeed, "DashJump": BuffFlagDashJump,
		"MonsterRiding": BuffFlagMonsterRiding, "SpeedInfusion": BuffFlagSpeedInfusion, "HomingBeacon": BuffFlagHomingBeacon,
		"ExpRate": BuffFlagExpRate, "DropRate": BuffFlagDropRate, "MesoRate": BuffFlagMesoRate,
	}
}
