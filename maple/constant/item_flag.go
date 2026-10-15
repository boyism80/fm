package constant

type ItemFlag uint16

const (
	ItemFlagLock        ItemFlag = 0x01
	ItemFlagSpikes      ItemFlag = 0x02
	ItemFlagCold        ItemFlag = 0x04
	ItemFlagUntradeable ItemFlag = 0x08
	ItemFlagEquipKarma  ItemFlag = 0x10
	ItemFlagBundleKarma ItemFlag = 0x02
)
