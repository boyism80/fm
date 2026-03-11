package constant

type Debuff struct {
	Mask uint32
}

var (
	DebuffWatk                = Debuff{0x1}
	DebuffWdef                = Debuff{0x2}
	DebuffMatk                = Debuff{0x4}
	DebuffMdef                = Debuff{0x8}
	DebuffAcc                 = Debuff{0x10}
	DebuffAvoid               = Debuff{0x20}
	DebuffSpeed               = Debuff{0x40}
	DebuffStun                = Debuff{0x80}
	DebuffFreeze              = Debuff{0x100}
	DebuffPoison              = Debuff{0x200}
	DebuffSeal                = Debuff{0x400}
	DebuffShowdown            = Debuff{0x800}
	DebuffWeaponAttackUp      = Debuff{0x1000}
	DebuffWeaponDefenseUp     = Debuff{0x2000}
	DebuffMagicAttackUp       = Debuff{0x4000}
	DebuffMagicDefenseUp      = Debuff{0x8000}
	DebuffDoom                = Debuff{0x10000}
	DebuffShadowWeb           = Debuff{0x20000}
	DebuffWeaponImmunity      = Debuff{0x40000}
	DebuffMagicImmunity       = Debuff{0x80000}
	DebuffDamageImmunity      = Debuff{0x200000}
	DebuffNinjaAmbush         = Debuff{0x400000}
	DebuffVenom               = Debuff{0x1000000}
	DebuffBlind               = Debuff{0x2000000}
	DebuffSealSkill           = Debuff{0x4000000}
	DebuffHypnotize           = Debuff{0x10000000}
	DebuffWeaponDamageReflect = Debuff{0x20000000}
	DebuffMagicDamageReflect  = Debuff{0x40000000}
)

func AllDebuffs() map[string]Debuff {
	return map[string]Debuff{
		"Watk": DebuffWatk, "Wdef": DebuffWdef, "Matk": DebuffMatk, "Mdef": DebuffMdef,
		"Acc": DebuffAcc, "Avoid": DebuffAvoid, "Speed": DebuffSpeed,
		"Stun": DebuffStun, "Freeze": DebuffFreeze, "Poison": DebuffPoison, "Seal": DebuffSeal,
		"Showdown":       DebuffShowdown,
		"WeaponAttackUp": DebuffWeaponAttackUp, "WeaponDefenseUp": DebuffWeaponDefenseUp,
		"MagicAttackUp": DebuffMagicAttackUp, "MagicDefenseUp": DebuffMagicDefenseUp,
		"Doom": DebuffDoom, "ShadowWeb": DebuffShadowWeb,
		"WeaponImmunity": DebuffWeaponImmunity, "MagicImmunity": DebuffMagicImmunity,
		"DamageImmunity": DebuffDamageImmunity, "NinjaAmbush": DebuffNinjaAmbush,
		"Venom": DebuffVenom, "Blind": DebuffBlind, "SealSkill": DebuffSealSkill,
		"Hypnotize":           DebuffHypnotize,
		"WeaponDamageReflect": DebuffWeaponDamageReflect, "MagicDamageReflect": DebuffMagicDamageReflect,
	}
}
