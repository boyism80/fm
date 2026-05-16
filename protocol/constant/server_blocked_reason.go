package constant

type ServerBlockedReason uint8

const (
	ServerBlockedChannelMoveUnavailable ServerBlockedReason = 1
	ServerBlockedCashShopUnavailable    ServerBlockedReason = 2
	ServerBlockedTradingShopUnavailable ServerBlockedReason = 3
	ServerBlockedTradeShopUserLimit     ServerBlockedReason = 4
)
