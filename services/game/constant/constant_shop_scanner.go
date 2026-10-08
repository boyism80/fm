package constant

func IsShopScanner(itemID uint32) bool {
	return itemID/10000 == 523 || itemID/10000 == 231
}
