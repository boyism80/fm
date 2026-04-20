package wz

type Equipment interface {
	Item
	GetEnchantChance() uint8
	GetRequired() RequiredStats
	GetAbility() AbilityStats
}
