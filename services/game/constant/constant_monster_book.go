package constant

const (
	MonsterBookCardMax   uint32 = 5
	MonsterBookCardFirst uint32 = 2380000
)

func IsMonsterCard(itemID uint32) bool {
	return GetConsumeType(itemID) == ConsumeTypeCard
}

func IsSpecialMonsterCard(itemID uint32) bool {
	return itemID/1000 == 2388
}
