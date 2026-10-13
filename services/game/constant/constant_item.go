package constant

const (
	ItemCategoryShuriken    uint32 = 207
	ItemCategoryPetFood     uint32 = 212
	ItemCategoryBullet      uint32 = 233
	ItemCategoryMarriage    uint32 = 421
	ItemCategoryPet         uint32 = 500
	ItemCategoryPetNameTag  uint32 = 517
	ItemCategoryPetSkill    uint32 = 519
	ItemCategoryPetCashFood uint32 = 524
)

const (
	ItemShurikenBase uint32 = 2070000
	ItemBulletBase   uint32 = 2330000
)

type ItemTrait uint32

const (
	ItemTraitRechargeable ItemTrait = 1 << iota
	ItemTraitPet
	ItemTraitMarriage
	ItemTraitEngagementRing
	ItemTraitWeddingRing
	ItemTraitCrushRing
	ItemTraitFriendshipRing
	ItemTraitSpecialMonsterCard
	ItemTraitShopScanner
	ItemTraitCashTeleportStone
	ItemTraitVipTeleportStone
)

const ChaosScrollRange = 5

const (
	ShopTransactionBuyOK    byte   = 0
	ShopTransactionUpdateOK byte   = 0x8
	ShopUnspecifiedQuantity uint16 = 0xFFFF
)
