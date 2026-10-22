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
	// Picking up the Longinus Spear skips the client-reported range check, as in the original server.
	ItemLonginusSpear uint32 = 4001025
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
	ItemTraitCash
	ItemTraitQuest
	ItemTraitTradeBlock
	ItemTraitUnique
	ItemTraitAccountSharable
	ItemTraitEquipTradeBlock
	ItemTraitConsumeOnPickup
	ItemTraitSummonBag
	ItemTraitCatchItem
	ItemTraitSkillBook
	ItemTraitItemEffect
	ItemTraitEntrustedShopPermit
	ItemTraitPersonalShopPermit
	ItemTraitNotSale
	ItemTraitPartyQuest
	ItemTraitScroll
)

const ChaosScrollRange = 5

const (
	ShopTransactionBuyOK    byte   = 0
	ShopTransactionUpdateOK byte   = 0x8
	ShopUnspecifiedQuantity uint16 = 0xFFFF
)
