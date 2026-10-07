package constant

type CashShopResultKind uint8

const (
	CashShopResultLocker          CashShopResultKind = 44
	CashShopResultGifts           CashShopResultKind = 46
	CashShopResultWishlist        CashShopResultKind = 48
	CashShopResultWishlistUpdated CashShopResultKind = 50
	CashShopResultBuyFailed       CashShopResultKind = 51
	CashShopResultBought          CashShopResultKind = 52
	CashShopResultTakenOut        CashShopResultKind = 67
	CashShopResultTakeOutFailed   CashShopResultKind = 68
	CashShopResultPutIn           CashShopResultKind = 69
	CashShopResultPutInFailed     CashShopResultKind = 70
)

type CashShopFailure uint8

const (
	CashShopFailureUnknown           CashShopFailure = 0
	CashShopFailureNotEnoughCash     CashShopFailure = 122
	CashShopFailureGender            CashShopFailure = 127
	CashShopFailureLockerFull        CashShopFailure = 129
	CashShopFailureNotEnoughSlots    CashShopFailure = 141
	CashShopFailureNotPurchasableNow CashShopFailure = 145
)
