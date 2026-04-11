package constant

type MobBuffFlag uint32

const (
	MobBuffWatk                MobBuffFlag = 0x1
	MobBuffWdef                MobBuffFlag = 0x2
	MobBuffMatk                MobBuffFlag = 0x4
	MobBuffMdef                MobBuffFlag = 0x8
	MobBuffAcc                 MobBuffFlag = 0x10
	MobBuffAvoid               MobBuffFlag = 0x20
	MobBuffSpeed               MobBuffFlag = 0x40
	MobBuffStun                MobBuffFlag = 0x80
	MobBuffFreeze              MobBuffFlag = 0x100
	MobBuffPoison              MobBuffFlag = 0x200
	MobBuffSeal                MobBuffFlag = 0x400
	MobBuffShowdown            MobBuffFlag = 0x800
	MobBuffWeaponAttackUp      MobBuffFlag = 0x1000
	MobBuffWeaponDefenseUp     MobBuffFlag = 0x2000
	MobBuffMagicAttackUp       MobBuffFlag = 0x4000
	MobBuffMagicDefenseUp      MobBuffFlag = 0x8000
	MobBuffDoom                MobBuffFlag = 0x10000
	MobBuffShadowWeb           MobBuffFlag = 0x20000
	MobBuffWeaponImmunity      MobBuffFlag = 0x40000
	MobBuffMagicImmunity       MobBuffFlag = 0x80000
	MobBuffDamageImmunity      MobBuffFlag = 0x200000
	MobBuffNinjaAmbush         MobBuffFlag = 0x400000
	MobBuffVenom               MobBuffFlag = 0x1000000
	MobBuffBlind               MobBuffFlag = 0x2000000
	MobBuffSealSkill           MobBuffFlag = 0x4000000
	MobBuffHypnotize           MobBuffFlag = 0x10000000
	MobBuffWeaponDamageReflect MobBuffFlag = 0x20000000
	MobBuffMagicDamageReflect  MobBuffFlag = 0x40000000
)

// Has returns true if this MobBuff bitmask includes the given filter (bitwise AND).
func (t MobBuffFlag) Has(other MobBuffFlag) bool {
	return (t & other) != 0
}

func AllMobBuffs() map[string]MobBuffFlag {
	return map[string]MobBuffFlag{
		"Watk":                MobBuffWatk,
		"Wdef":                MobBuffWdef,
		"Matk":                MobBuffMatk,
		"Mdef":                MobBuffMdef,
		"Acc":                 MobBuffAcc,
		"Avoid":               MobBuffAvoid,
		"Speed":               MobBuffSpeed,
		"Stun":                MobBuffStun,
		"Freeze":              MobBuffFreeze,
		"Poison":              MobBuffPoison,
		"Seal":                MobBuffSeal,
		"Showdown":            MobBuffShowdown,
		"WeaponAttackUp":      MobBuffWeaponAttackUp,
		"WeaponDefenseUp":     MobBuffWeaponDefenseUp,
		"MagicAttackUp":       MobBuffMagicAttackUp,
		"MagicDefenseUp":      MobBuffMagicDefenseUp,
		"Doom":                MobBuffDoom,
		"ShadowWeb":           MobBuffShadowWeb,
		"WeaponImmunity":      MobBuffWeaponImmunity,
		"MagicImmunity":       MobBuffMagicImmunity,
		"DamageImmunity":      MobBuffDamageImmunity,
		"NinjaAmbush":         MobBuffNinjaAmbush,
		"Venom":               MobBuffVenom,
		"Blind":               MobBuffBlind,
		"SealSkill":           MobBuffSealSkill,
		"Hypnotize":           MobBuffHypnotize,
		"WeaponDamageReflect": MobBuffWeaponDamageReflect,
		"MagicDamageReflect":  MobBuffMagicDamageReflect,
	}
}
