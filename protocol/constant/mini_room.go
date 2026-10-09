package constant

type MiniRoomMode uint8

const (
	MiniRoomCreate             MiniRoomMode = 0x00
	MiniRoomInvite             MiniRoomMode = 0x02
	MiniRoomDecline            MiniRoomMode = 0x03
	MiniRoomVisit              MiniRoomMode = 0x04
	MiniRoomChat               MiniRoomMode = 0x06
	MiniRoomExit               MiniRoomMode = 0x0A
	MiniRoomOpen               MiniRoomMode = 0x0B
	MiniRoomTradePutItem       MiniRoomMode = 0x0D
	MiniRoomTradePutMeso       MiniRoomMode = 0x0E
	MiniRoomTradeConfirm       MiniRoomMode = 0x0F
	MiniRoomPersonalAddItem    MiniRoomMode = 0x12
	MiniRoomPersonalBuy        MiniRoomMode = 0x13
	MiniRoomPersonalRemoveItem MiniRoomMode = 0x17
	MiniRoomKick               MiniRoomMode = 0x18
	MiniRoomKickTimeout        MiniRoomMode = 0x19
	MiniRoomBlacklist          MiniRoomMode = 0x1A
	MiniRoomAddItem            MiniRoomMode = 0x1D
	MiniRoomBuy                MiniRoomMode = 0x1E
	MiniRoomRemoveItem         MiniRoomMode = 0x22
	MiniRoomMaintenanceOff     MiniRoomMode = 0x23
	MiniRoomArrange            MiniRoomMode = 0x24
	MiniRoomClose              MiniRoomMode = 0x25
	MiniRoomWithdrawMeso       MiniRoomMode = 0x27
	MiniRoomRequestTie         MiniRoomMode = 0x2A
	MiniRoomAnswerTie          MiniRoomMode = 0x2B
	MiniRoomGiveUp             MiniRoomMode = 0x2C
	MiniRoomRequestRetreat     MiniRoomMode = 0x2E
	MiniRoomAnswerRetreat      MiniRoomMode = 0x2F
	MiniRoomExitAfterGame      MiniRoomMode = 0x30
	MiniRoomCancelExit         MiniRoomMode = 0x31
	MiniRoomReady              MiniRoomMode = 0x32
	MiniRoomUnready            MiniRoomMode = 0x33
	MiniRoomExpel              MiniRoomMode = 0x34
	MiniRoomStart              MiniRoomMode = 0x35
	MiniRoomSkip               MiniRoomMode = 0x37
	MiniRoomMoveOmok           MiniRoomMode = 0x38
	MiniRoomSelectCard         MiniRoomMode = 0x3C
)

type MiniRoomResult uint8

const (
	MiniRoomResultInvite        MiniRoomResult = 0x02
	MiniRoomResultInviteResult  MiniRoomResult = 0x03
	MiniRoomResultVisitor       MiniRoomResult = 0x04
	MiniRoomResultEnter         MiniRoomResult = 0x05
	MiniRoomResultChat          MiniRoomResult = 0x06
	MiniRoomResultLeave         MiniRoomResult = 0x0A
	MiniRoomResultTradeItem     MiniRoomResult = 0x0D
	MiniRoomResultTradeMeso     MiniRoomResult = 0x0E
	MiniRoomResultTradeConfirm  MiniRoomResult = 0x0F
	MiniRoomResultBuy           MiniRoomResult = 0x14
	MiniRoomResultItems         MiniRoomResult = 0x15
	MiniRoomResultSold          MiniRoomResult = 0x16
	MiniRoomResultRemoved       MiniRoomResult = 0x17
	MiniRoomResultArranged      MiniRoomResult = 0x24
	MiniRoomResultClosed        MiniRoomResult = 0x26
	MiniRoomResultMesoWithdrawn MiniRoomResult = 0x28
	MiniRoomResultRequestTie    MiniRoomResult = 0x2A
	MiniRoomResultDenyTie       MiniRoomResult = 0x2B
	MiniRoomResultExitAfterGame MiniRoomResult = 0x30
	MiniRoomResultCancelExit    MiniRoomResult = 0x31
	MiniRoomResultReady         MiniRoomResult = 0x32
	MiniRoomResultUnready       MiniRoomResult = 0x33
	MiniRoomResultStart         MiniRoomResult = 0x35
	MiniRoomResultGameOver      MiniRoomResult = 0x36
	MiniRoomResultSkip          MiniRoomResult = 0x37
	MiniRoomResultMoveOmok      MiniRoomResult = 0x38
	MiniRoomResultSelectCard    MiniRoomResult = 0x3C
)

type MiniGameOutcome uint8

const (
	MiniGameGiveUp MiniGameOutcome = 0
	MiniGameTie    MiniGameOutcome = 1
	MiniGameWin    MiniGameOutcome = 2
)

const (
	MiniRoomTypeOmok          uint8 = 1
	MiniRoomTypeMatchCard     uint8 = 2
	MiniRoomTypeTrade         uint8 = 3
	MiniRoomTypePersonalShop  uint8 = 4
	MiniRoomTypeEntrustedShop uint8 = 5
	MiniRoomTradeUsers        uint8 = 2
	MiniRoomShopUsers         uint8 = 4
	MiniRoomGameUsers         uint8 = 2
	MiniRoomChatShop          uint8 = 8
)

type MiniRoomInviteResult uint8

const (
	MiniRoomInviteNotFound MiniRoomInviteResult = 1
	MiniRoomInviteBusy     MiniRoomInviteResult = 2
	MiniRoomInviteDeclined MiniRoomInviteResult = 3
	MiniRoomInviteBlocked  MiniRoomInviteResult = 4
)

type MiniRoomEnterError uint8

const (
	MiniRoomEnterClosed     MiniRoomEnterError = 1
	MiniRoomEnterFull       MiniRoomEnterError = 2
	MiniRoomEnterNearPortal MiniRoomEnterError = 10
	MiniRoomEnterCannotOpen MiniRoomEnterError = 11
	MiniRoomEnterFreeMarket MiniRoomEnterError = 13
	MiniRoomEnterBlocked    MiniRoomEnterError = 15
	MiniRoomEnterOrganizing MiniRoomEnterError = 16
)

type MiniRoomLeaveReason uint8

const (
	MiniRoomLeaveExit        MiniRoomLeaveReason = 0
	MiniRoomLeaveTradeCancel MiniRoomLeaveReason = 2
	MiniRoomLeaveShopClosed  MiniRoomLeaveReason = 3
	MiniRoomLeaveKicked      MiniRoomLeaveReason = 5
	MiniRoomLeaveTradeDone   MiniRoomLeaveReason = 6
	MiniRoomLeaveTradeFail   MiniRoomLeaveReason = 7
	MiniRoomLeaveTradeOnly   MiniRoomLeaveReason = 8
	MiniRoomLeaveSoldOut     MiniRoomLeaveReason = 10
	MiniRoomLeaveStayTimeout MiniRoomLeaveReason = 11
	MiniRoomLeaveOrganizing  MiniRoomLeaveReason = 13
	MiniRoomLeaveTimeUp      MiniRoomLeaveReason = 14
	MiniRoomLeaveMapMoved    MiniRoomLeaveReason = 15
	MiniRoomLeaveClosed      MiniRoomLeaveReason = 16
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
	EntrustedShopTitle          EntrustedShopCheck = 7
	EntrustedShopAlreadyOpen    EntrustedShopCheck = 8
	EntrustedShopStoreBankFull  EntrustedShopCheck = 9
	EntrustedShopAccountBusy    EntrustedShopCheck = 10
	EntrustedShopCannotOpen     EntrustedShopCheck = 11
	EntrustedShopRemoteLocation EntrustedShopCheck = 16
	EntrustedShopRemoteVisit    EntrustedShopCheck = 17
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

type ShopScannerResult uint8

const (
	ShopScannerResultSearch  ShopScannerResult = 6
	ShopScannerResultPopular ShopScannerResult = 7
)
