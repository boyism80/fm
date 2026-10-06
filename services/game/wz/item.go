package wz

type Item interface {
	GetID() uint32
	GetName() string
	GetPrice() int
	IsCash() bool
	IsQuest() bool
	GetCapacity() uint16
	IsTradeAvailable() int
	IsTradeBlock() bool
	IsOnly() bool
	IsAccountSharable() bool

	IsConsumeOnPickup() bool
}
