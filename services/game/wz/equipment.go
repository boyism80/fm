package wz

type Equipment interface {
	Item
	GetEnhanceChance() uint8
	GetRequired() RequiredStats
	GetAbility() AbilityStats
}
