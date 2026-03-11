package constant

type ConsumeType uint8

const (
	ConsumeTypeNone            ConsumeType = 0  // unknown or not consume
	ConsumeTypeReturnScroll    ConsumeType = 1  // 203 귀환주문서
	ConsumeTypeEnhanceScroll   ConsumeType = 2  // 204 장비강화주문서
	ConsumeTypeArrowBow        ConsumeType = 3  // 206 활전용화살
	ConsumeTypeArrowCrossBow   ConsumeType = 4  // 206 석궁전용화살
	ConsumeTypeShuriken        ConsumeType = 5  // 207 표창
	ConsumeTypeMonsterSummon   ConsumeType = 6  // 210 몬스터소환아이템
	ConsumeTypePetFood         ConsumeType = 7  // 212 펫먹이
	ConsumeTypeWeddingRing     ConsumeType = 8  // 224 웨딩링
	ConsumeTypeRidingFood      ConsumeType = 9  // 226 라이딩 먹이
	ConsumeTypeSkillBookGet    ConsumeType = 10 // 228 스킬마스터리북(스킬획득)
	ConsumeTypeSkillBookUnlock ConsumeType = 11 // 229 스킬마스터리북(마스터레벨 언락)
	ConsumeTypeMinervaOwl      ConsumeType = 12 // 231 미네르바부엉이
	ConsumeTypeTeleportRock    ConsumeType = 13 // 232 순간이동의돌
	ConsumeTypeBullet          ConsumeType = 14 // 233 총알
	ConsumeTypeCard            ConsumeType = 15 // 238 카드
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
