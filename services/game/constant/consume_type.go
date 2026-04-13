package constant

type ConsumeType uint8

const (
	ConsumeTypeNone            ConsumeType = 0
	ConsumeTypeReturnScroll    ConsumeType = 1
	ConsumeTypeEnhanceScroll   ConsumeType = 2
	ConsumeTypeArrowBow        ConsumeType = 3
	ConsumeTypeArrowCrossBow   ConsumeType = 4
	ConsumeTypeShuriken        ConsumeType = 5
	ConsumeTypeMonsterSummon   ConsumeType = 6
	ConsumeTypePetFood         ConsumeType = 7
	ConsumeTypeWeddingRing     ConsumeType = 8
	ConsumeTypeRidingFood      ConsumeType = 9
	ConsumeTypeSkillBookGet    ConsumeType = 10
	ConsumeTypeSkillBookUnlock ConsumeType = 11
	ConsumeTypeMinervaOwl      ConsumeType = 12
	ConsumeTypeTeleportRock    ConsumeType = 13
	ConsumeTypeBullet          ConsumeType = 14
	ConsumeTypeCard            ConsumeType = 15
)

func AllConsumeTypes() map[string]ConsumeType {
	return map[string]ConsumeType{
		"None":            ConsumeTypeNone,
		"ReturnScroll":    ConsumeTypeReturnScroll,
		"EnhanceScroll":   ConsumeTypeEnhanceScroll,
		"ArrowBow":        ConsumeTypeArrowBow,
		"ArrowCrossBow":   ConsumeTypeArrowCrossBow,
		"Shuriken":        ConsumeTypeShuriken,
		"MonsterSummon":   ConsumeTypeMonsterSummon,
		"PetFood":         ConsumeTypePetFood,
		"WeddingRing":     ConsumeTypeWeddingRing,
		"RidingFood":      ConsumeTypeRidingFood,
		"SkillBookGet":    ConsumeTypeSkillBookGet,
		"SkillBookUnlock": ConsumeTypeSkillBookUnlock,
		"MinervaOwl":      ConsumeTypeMinervaOwl,
		"TeleportRock":    ConsumeTypeTeleportRock,
		"Bullet":          ConsumeTypeBullet,
		"Card":            ConsumeTypeCard,
	}
}
