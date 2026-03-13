package constant

type MobStatus uint32

const (
	MobStatusWatk                MobStatus = 0x1
	MobStatusWdef                MobStatus = 0x2
	MobStatusMatk                MobStatus = 0x4
	MobStatusMdef                MobStatus = 0x8
	MobStatusAcc                 MobStatus = 0x10
	MobStatusAvoid               MobStatus = 0x20
	MobStatusSpeed               MobStatus = 0x40
	MobStatusStun                MobStatus = 0x80
	MobStatusFreeze              MobStatus = 0x100
	MobStatusPoison              MobStatus = 0x200
	MobStatusSeal                MobStatus = 0x400
	MobStatusShowdown            MobStatus = 0x800
	MobStatusWeaponAttackUp      MobStatus = 0x1000
	MobStatusWeaponDefenseUp     MobStatus = 0x2000
	MobStatusMagicAttackUp       MobStatus = 0x4000
	MobStatusMagicDefenseUp      MobStatus = 0x8000
	MobStatusDoom                MobStatus = 0x10000
	MobStatusShadowWeb           MobStatus = 0x20000
	MobStatusWeaponImmunity      MobStatus = 0x40000
	MobStatusMagicImmunity       MobStatus = 0x80000
	MobStatusDamageImmunity      MobStatus = 0x200000
	MobStatusNinjaAmbush         MobStatus = 0x400000
	MobStatusVenom               MobStatus = 0x1000000
	MobStatusBlind               MobStatus = 0x2000000
	MobStatusSealSkill           MobStatus = 0x4000000
	MobStatusHypnotize           MobStatus = 0x10000000
	MobStatusWeaponDamageReflect MobStatus = 0x20000000
	MobStatusMagicDamageReflect  MobStatus = 0x40000000
)

// Has returns true if this MobStatus bitmask includes the given filter (bitwise AND).
func (t MobStatus) Has(other MobStatus) bool {
	return (t & other) != 0
}

func AllMobStatuses() map[string]MobStatus {
	return map[string]MobStatus{
		"Watk": MobStatusWatk, "Wdef": MobStatusWdef, "Matk": MobStatusMatk, "Mdef": MobStatusMdef,
		"Acc": MobStatusAcc, "Avoid": MobStatusAvoid, "Speed": MobStatusSpeed,
		"Stun": MobStatusStun, "Freeze": MobStatusFreeze, "Poison": MobStatusPoison, "Seal": MobStatusSeal,
		"Showdown":       MobStatusShowdown,
		"WeaponAttackUp": MobStatusWeaponAttackUp, "WeaponDefenseUp": MobStatusWeaponDefenseUp,
		"MagicAttackUp": MobStatusMagicAttackUp, "MagicDefenseUp": MobStatusMagicDefenseUp,
		"Doom": MobStatusDoom, "ShadowWeb": MobStatusShadowWeb,
		"WeaponImmunity": MobStatusWeaponImmunity, "MagicImmunity": MobStatusMagicImmunity,
		"DamageImmunity": MobStatusDamageImmunity, "NinjaAmbush": MobStatusNinjaAmbush,
		"Venom": MobStatusVenom, "Blind": MobStatusBlind, "SealSkill": MobStatusSealSkill,
		"Hypnotize":           MobStatusHypnotize,
		"WeaponDamageReflect": MobStatusWeaponDamageReflect, "MagicDamageReflect": MobStatusMagicDamageReflect,
	}
}
