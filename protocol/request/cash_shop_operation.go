package request

import (
	"github.com/boyism80/fm/stream"
)

type CashShopAction uint8

const (
	CashShopActionBuy            CashShopAction = 3
	CashShopActionGift           CashShopAction = 4
	CashShopActionWishlist       CashShopAction = 5
	CashShopActionInventorySlots CashShopAction = 6
	CashShopActionStorageSlots   CashShopAction = 7
	CashShopActionCharacterSlots CashShopAction = 8
	CashShopActionTakeOut        CashShopAction = 0x0C
	CashShopActionPutIn          CashShopAction = 0x0D
	CashShopActionPayBack        CashShopAction = 0x19
	CashShopActionCoupleRing     CashShopAction = 0x1C
	CashShopActionPackage        CashShopAction = 0x1D
	CashShopActionGiftPackage    CashShopAction = 0x1E
	CashShopActionQuestItem      CashShopAction = 0x1F
	CashShopActionFriendshipRing CashShopAction = 0x22
)

const (
	cashShopSlotExpansionByDefault uint8 = 0
	cashShopSlotExpansionByCoupon  uint8 = 1
)

type CashShopOperation struct {
	Action        CashShopAction
	Currency      uint8
	CommoditySN   uint32
	Wishlist      []uint32
	Serial        uint64
	InventoryType uint8
	Slot          int16
	Birthday      uint32
	Recipient     string
	Message       string
	ByCoupon      bool
}

func (*CashShopOperation) Opcode() byte {
	return 0xBC
}

func (a *CashShopOperation) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(a.Action))
	switch a.Action {
	case CashShopActionBuy, CashShopActionPackage, CashShopActionCharacterSlots:
		writer.WriteU8(a.Currency)
		writer.WriteU32(a.CommoditySN)
	case CashShopActionGift, CashShopActionGiftPackage, CashShopActionCoupleRing, CashShopActionFriendshipRing:
		writer.WriteU32(a.Birthday)
		writer.WriteU32(a.CommoditySN)
		writer.WriteStr16(a.Recipient)
		writer.WriteStr16(a.Message)
	case CashShopActionWishlist:
		for i := range 10 {
			if i < len(a.Wishlist) {
				writer.WriteU32(a.Wishlist[i])
			} else {
				writer.WriteU32(0)
			}
		}
	case CashShopActionInventorySlots, CashShopActionStorageSlots:
		writer.WriteU8(a.Currency)
		if a.ByCoupon {
			writer.WriteU8(cashShopSlotExpansionByCoupon)
			writer.WriteU32(a.CommoditySN)
			break
		}
		writer.WriteU8(cashShopSlotExpansionByDefault)
		if a.Action == CashShopActionInventorySlots {
			writer.WriteU8(a.InventoryType)
		}
	case CashShopActionTakeOut:
		writer.WriteU64(a.Serial)
		writer.WriteU8(a.InventoryType)
		writer.WriteU16(uint16(a.Slot))
	case CashShopActionPutIn:
		writer.WriteU64(a.Serial)
		writer.WriteU8(a.InventoryType)
	case CashShopActionPayBack:
		writer.WriteU32(a.Birthday)
		writer.WriteU64(a.Serial)
	case CashShopActionQuestItem:
		writer.WriteU32(a.CommoditySN)
	}
	return nil
}

func (a *CashShopOperation) Deserialize(reader *stream.StreamReader) {
	a.Action = CashShopAction(reader.ReadU8())
	switch a.Action {
	case CashShopActionBuy, CashShopActionPackage, CashShopActionCharacterSlots:
		a.Currency = reader.ReadU8()
		a.CommoditySN = reader.ReadU32()
	case CashShopActionGift, CashShopActionGiftPackage, CashShopActionCoupleRing, CashShopActionFriendshipRing:
		a.Birthday = reader.ReadU32()
		a.CommoditySN = reader.ReadU32()
		a.Recipient = reader.ReadStr16()
		a.Message = reader.ReadStr16()
	case CashShopActionWishlist:
		a.Wishlist = make([]uint32, 10)
		for i := range a.Wishlist {
			a.Wishlist[i] = reader.ReadU32()
		}
	case CashShopActionInventorySlots, CashShopActionStorageSlots:
		a.Currency = reader.ReadU8()
		a.ByCoupon = reader.ReadU8() == cashShopSlotExpansionByCoupon
		switch {
		case a.ByCoupon:
			a.CommoditySN = reader.ReadU32()
		case a.Action == CashShopActionInventorySlots:
			a.InventoryType = reader.ReadU8()
		}
	case CashShopActionTakeOut:
		a.Serial = reader.ReadU64()
		a.InventoryType = reader.ReadU8()
		a.Slot = int16(reader.ReadU16())
	case CashShopActionPutIn:
		a.Serial = reader.ReadU64()
		a.InventoryType = reader.ReadU8()
	case CashShopActionPayBack:
		a.Birthday = reader.ReadU32()
		a.Serial = reader.ReadU64()
	case CashShopActionQuestItem:
		a.CommoditySN = reader.ReadU32()
	}
}
