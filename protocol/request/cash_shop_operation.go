package request

import (
	"github.com/boyism80/fm/stream"
)

type CashShopAction uint8

const (
	CashShopActionBuy      CashShopAction = 3
	CashShopActionWishlist CashShopAction = 5
	CashShopActionTakeOut  CashShopAction = 0x0C
	CashShopActionPutIn    CashShopAction = 0x0D
)

type CashShopOperation struct {
	Action        CashShopAction
	Currency      uint8
	CommoditySN   uint32
	Wishlist      []uint32
	Serial        uint64
	InventoryType uint8
	Slot          int16
}

func (*CashShopOperation) Opcode() byte {
	return 0xBC
}

func (a *CashShopOperation) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(a.Action))
	switch a.Action {
	case CashShopActionBuy:
		writer.WriteU8(a.Currency)
		writer.WriteU32(a.CommoditySN)
	case CashShopActionWishlist:
		for i := range 10 {
			if i < len(a.Wishlist) {
				writer.WriteU32(a.Wishlist[i])
			} else {
				writer.WriteU32(0)
			}
		}
	case CashShopActionTakeOut:
		writer.WriteU64(a.Serial)
		writer.WriteU8(a.InventoryType)
		writer.WriteU16(uint16(a.Slot))
	case CashShopActionPutIn:
		writer.WriteU64(a.Serial)
		writer.WriteU8(a.InventoryType)
	}
	return nil
}

func (a *CashShopOperation) Deserialize(reader *stream.StreamReader) {
	a.Action = CashShopAction(reader.ReadU8())
	switch a.Action {
	case CashShopActionBuy:
		a.Currency = reader.ReadU8()
		a.CommoditySN = reader.ReadU32()
	case CashShopActionWishlist:
		a.Wishlist = make([]uint32, 10)
		for i := range a.Wishlist {
			a.Wishlist[i] = reader.ReadU32()
		}
	case CashShopActionTakeOut:
		a.Serial = reader.ReadU64()
		a.InventoryType = reader.ReadU8()
		a.Slot = int16(reader.ReadU16())
	case CashShopActionPutIn:
		a.Serial = reader.ReadU64()
		a.InventoryType = reader.ReadU8()
	}
}
