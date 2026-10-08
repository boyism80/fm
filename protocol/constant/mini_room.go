package constant

type MiniRoomMode uint8

const (
	MiniRoomCreate         MiniRoomMode = 0x00
	MiniRoomVisit          MiniRoomMode = 0x04
	MiniRoomChat           MiniRoomMode = 0x06
	MiniRoomExit           MiniRoomMode = 0x0A
	MiniRoomOpen           MiniRoomMode = 0x0B
	MiniRoomAddItem        MiniRoomMode = 0x1D
	MiniRoomBuy            MiniRoomMode = 0x1E
	MiniRoomRemoveItem     MiniRoomMode = 0x22
	MiniRoomMaintenanceOff MiniRoomMode = 0x23
	MiniRoomArrange        MiniRoomMode = 0x24
	MiniRoomClose          MiniRoomMode = 0x25
	MiniRoomWithdrawMeso   MiniRoomMode = 0x27
)

type MiniRoomResult uint8

const (
	MiniRoomResultVisitor       MiniRoomResult = 0x04
	MiniRoomResultEnter         MiniRoomResult = 0x05
	MiniRoomResultChat          MiniRoomResult = 0x06
	MiniRoomResultLeave         MiniRoomResult = 0x0A
	MiniRoomResultBuy           MiniRoomResult = 0x14
	MiniRoomResultItems         MiniRoomResult = 0x15
	MiniRoomResultArranged      MiniRoomResult = 0x24
	MiniRoomResultClosed        MiniRoomResult = 0x26
	MiniRoomResultMesoWithdrawn MiniRoomResult = 0x28
)

const (
	MiniRoomTypeHiredMerchant  uint8 = 5
	MiniRoomHiredMerchantUsers uint8 = 4
	MiniRoomChatShop           uint8 = 8
)

type MiniRoomEnterError uint8

const (
	MiniRoomEnterClosed     MiniRoomEnterError = 1
	MiniRoomEnterFull       MiniRoomEnterError = 2
	MiniRoomEnterNearPortal MiniRoomEnterError = 10
	MiniRoomEnterCannotOpen MiniRoomEnterError = 11
	MiniRoomEnterOrganizing MiniRoomEnterError = 16
)

type MiniRoomLeaveReason uint8

const (
	MiniRoomLeaveExit       MiniRoomLeaveReason = 0
	MiniRoomLeaveOrganizing MiniRoomLeaveReason = 13
	MiniRoomLeaveTimeUp     MiniRoomLeaveReason = 14
	MiniRoomLeaveClosed     MiniRoomLeaveReason = 16
)

type MiniRoomBuyResult uint8

const (
	MiniRoomBuyNotEnoughItem MiniRoomBuyResult = 1
	MiniRoomBuyNotEnoughMeso MiniRoomBuyResult = 2
	MiniRoomBuySellerLimit   MiniRoomBuyResult = 4
	MiniRoomBuyInventoryFull MiniRoomBuyResult = 5
	MiniRoomBuyOnlyOne       MiniRoomBuyResult = 6
	MiniRoomBuyUnknown       MiniRoomBuyResult = 7
)

type MiniRoomCloseResult uint8

const (
	MiniRoomCloseAll           MiniRoomCloseResult = 0
	MiniRoomCloseMesoOver      MiniRoomCloseResult = 1
	MiniRoomCloseOnlyOne       MiniRoomCloseResult = 2
	MiniRoomCloseInventoryFull MiniRoomCloseResult = 3
)

type EntrustedShopCheck uint8

const (
	EntrustedShopTitle         EntrustedShopCheck = 7
	EntrustedShopAlreadyOpen   EntrustedShopCheck = 8
	EntrustedShopStoreBankFull EntrustedShopCheck = 9
	EntrustedShopAccountBusy   EntrustedShopCheck = 10
	EntrustedShopCannotOpen    EntrustedShopCheck = 11
)

type StoreBankMode uint8

const (
	StoreBankWithdraw StoreBankMode = 0x19
	StoreBankConfirm  StoreBankMode = 0x1A
	StoreBankClose    StoreBankMode = 0x1B
)

type StoreBankResult uint8

const (
	StoreBankResultOpen          StoreBankResult = 0x23
	StoreBankResultFee           StoreBankResult = 0x24
	StoreBankResultLocation      StoreBankResult = 0x25
	StoreBankResultClaimed       StoreBankResult = 0x1D
	StoreBankResultMesoOver      StoreBankResult = 0x1E
	StoreBankResultOnlyOne       StoreBankResult = 0x1F
	StoreBankResultNotEnoughMeso StoreBankResult = 0x20
	StoreBankResultInventoryFull StoreBankResult = 0x21
)
