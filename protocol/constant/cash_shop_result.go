package constant

type CashShopResultKind uint8

const (
	CashShopResultLocker               CashShopResultKind = 44
	CashShopResultGifts                CashShopResultKind = 46
	CashShopResultWishlist             CashShopResultKind = 48
	CashShopResultWishlistUpdated      CashShopResultKind = 50
	CashShopResultWishlistFailed       CashShopResultKind = 51
	CashShopResultBought               CashShopResultKind = 52
	CashShopResultBuyFailed            CashShopResultKind = 53
	CashShopResultCouponRedeemed       CashShopResultKind = 54
	CashShopResultCouponFailed         CashShopResultKind = 57
	CashShopResultGiftSent             CashShopResultKind = 59
	CashShopResultGiftFailed           CashShopResultKind = 60
	CashShopResultInventorySlots       CashShopResultKind = 61
	CashShopResultInventorySlotsFailed CashShopResultKind = 62
	CashShopResultStorageSlots         CashShopResultKind = 63
	CashShopResultStorageSlotsFailed   CashShopResultKind = 64
	CashShopResultCharacterSlots       CashShopResultKind = 65
	CashShopResultCharacterSlotsFailed CashShopResultKind = 66
	CashShopResultTakenOut             CashShopResultKind = 67
	CashShopResultTakeOutFailed        CashShopResultKind = 68
	CashShopResultPutIn                CashShopResultKind = 69
	CashShopResultPutInFailed          CashShopResultKind = 70
	CashShopResultExpired              CashShopResultKind = 73
	CashShopResultPaidBack             CashShopResultKind = 96
	CashShopResultPayBackFailed        CashShopResultKind = 97
	CashShopResultPackageBought        CashShopResultKind = 100
	CashShopResultPackageFailed        CashShopResultKind = 101
	CashShopResultQuestItemBought      CashShopResultKind = 104
	CashShopResultQuestItemFailed      CashShopResultKind = 105
)

type CashShopFailure uint8

const (
	CashShopFailureUnknown           CashShopFailure = 0
	CashShopFailureNotEnoughCash     CashShopFailure = 122
	CashShopFailureSameAccount       CashShopFailure = 125
	CashShopFailureWrongName         CashShopFailure = 126
	CashShopFailureGender            CashShopFailure = 127
	CashShopFailureRecipientFull     CashShopFailure = 128
	CashShopFailureLockerFull        CashShopFailure = 129
	CashShopFailureCouponWrong       CashShopFailure = 131
	CashShopFailureCouponUsed        CashShopFailure = 133
	CashShopFailureCouponNoGift      CashShopFailure = 139
	CashShopFailureNotEnoughSlots    CashShopFailure = 141
	CashShopFailureRing              CashShopFailure = 143
	CashShopFailureNotPurchasableNow CashShopFailure = 145
	CashShopFailureNotEnoughMeso     CashShopFailure = 148
)
