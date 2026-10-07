package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type CashShopResult struct {
	Kind           constant.CashShopResultKind
	Locker         []*dto.CashShopItem
	StorageSlots   uint16
	CharacterSlots uint16
	Wishlist       []uint32
	Failure        constant.CashShopFailure
	CommoditySN    uint32
	Item           *dto.CashShopItem
	Items          []*dto.CashShopItem
	Gifts          []*dto.CashShopGift
	Slot           int16
	TakenOut       dto.Item
	Recipient      string
	ItemID         uint32
	Count          uint16
	InventoryType  uint8
	Slots          uint16
	Serial         uint64
	MaplePoint     uint32
	Meso           uint32
	Granted        []*dto.CashShopGrantedItem
}

func (p *CashShopResult) Opcode() uint16 {
	return 0xFA
}

func (p *CashShopResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Kind))
	switch p.Kind {
	case constant.CashShopResultLocker:
		writer.WriteU16(uint16(len(p.Locker)))
		for _, item := range p.Locker {
			item.Serialize(writer)
		}
		writer.WriteU16(p.StorageSlots)
		writer.WriteU16(p.CharacterSlots)
	case constant.CashShopResultGifts:
		writer.WriteU16(uint16(len(p.Gifts)))
		for _, gift := range p.Gifts {
			gift.Serialize(writer)
		}
	case constant.CashShopResultWishlist, constant.CashShopResultWishlistUpdated:
		for i := range 10 {
			if i < len(p.Wishlist) {
				writer.WriteU32(p.Wishlist[i])
			} else {
				writer.WriteU32(0)
			}
		}
	case constant.CashShopResultBuyFailed, constant.CashShopResultGiftFailed:
		writer.WriteU8(uint8(p.Failure))
		if p.Failure == constant.CashShopFailureNotPurchasableNow {
			writer.WriteU32(p.CommoditySN)
		}
	case constant.CashShopResultWishlistFailed, constant.CashShopResultCouponFailed,
		constant.CashShopResultInventorySlotsFailed, constant.CashShopResultStorageSlotsFailed, constant.CashShopResultCharacterSlotsFailed,
		constant.CashShopResultTakeOutFailed, constant.CashShopResultPutInFailed, constant.CashShopResultPayBackFailed,
		constant.CashShopResultPackageFailed, constant.CashShopResultQuestItemFailed:
		writer.WriteU8(uint8(p.Failure))
	case constant.CashShopResultBought, constant.CashShopResultPutIn:
		p.Item.Serialize(writer)
	case constant.CashShopResultCouponRedeemed:
		writer.WriteU8(uint8(len(p.Items)))
		for _, item := range p.Items {
			item.Serialize(writer)
		}
		writer.WriteU32(p.MaplePoint)
		writer.WriteU32(uint32(len(p.Granted)))
		for _, granted := range p.Granted {
			granted.Serialize(writer)
		}
		writer.WriteU32(p.Meso)
	case constant.CashShopResultGiftSent:
		writer.WriteStr16(p.Recipient)
		writer.WriteU32(p.ItemID)
		writer.WriteU16(p.Count)
	case constant.CashShopResultInventorySlots:
		writer.WriteU8(p.InventoryType)
		writer.WriteU16(p.Slots)
	case constant.CashShopResultStorageSlots, constant.CashShopResultCharacterSlots:
		writer.WriteU16(p.Slots)
	case constant.CashShopResultTakenOut:
		writer.WriteU16(uint16(p.Slot))
		p.TakenOut.Serialize(writer, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
	case constant.CashShopResultExpired:
		writer.WriteU64(p.Serial)
	case constant.CashShopResultPaidBack:
		writer.WriteU64(p.Serial)
		writer.WriteU32(p.MaplePoint)
	case constant.CashShopResultPackageBought:
		writer.WriteU8(uint8(len(p.Items)))
		for _, item := range p.Items {
			item.Serialize(writer)
		}
		writer.WriteU16(0)
	case constant.CashShopResultQuestItemBought:
		writer.WriteU32(uint32(len(p.Granted)))
		for _, granted := range p.Granted {
			granted.Serialize(writer)
		}
	}
	return nil
}

func (p *CashShopResult) Deserialize(reader *stream.StreamReader) {
	p.Kind = constant.CashShopResultKind(reader.ReadU8())
	switch p.Kind {
	case constant.CashShopResultLocker:
		count := reader.ReadU16()
		p.Locker = make([]*dto.CashShopItem, 0, count)
		for range count {
			item := &dto.CashShopItem{}
			item.Deserialize(reader)
			p.Locker = append(p.Locker, item)
		}
		p.StorageSlots = reader.ReadU16()
		p.CharacterSlots = reader.ReadU16()
	case constant.CashShopResultGifts:
		count := reader.ReadU16()
		p.Gifts = make([]*dto.CashShopGift, 0, count)
		for range count {
			gift := &dto.CashShopGift{}
			gift.Deserialize(reader)
			p.Gifts = append(p.Gifts, gift)
		}
	case constant.CashShopResultWishlist, constant.CashShopResultWishlistUpdated:
		p.Wishlist = make([]uint32, 10)
		for i := range p.Wishlist {
			p.Wishlist[i] = reader.ReadU32()
		}
	case constant.CashShopResultBuyFailed, constant.CashShopResultGiftFailed:
		p.Failure = constant.CashShopFailure(reader.ReadU8())
		if p.Failure == constant.CashShopFailureNotPurchasableNow {
			p.CommoditySN = reader.ReadU32()
		}
	case constant.CashShopResultWishlistFailed, constant.CashShopResultCouponFailed,
		constant.CashShopResultInventorySlotsFailed, constant.CashShopResultStorageSlotsFailed, constant.CashShopResultCharacterSlotsFailed,
		constant.CashShopResultTakeOutFailed, constant.CashShopResultPutInFailed, constant.CashShopResultPayBackFailed,
		constant.CashShopResultPackageFailed, constant.CashShopResultQuestItemFailed:
		p.Failure = constant.CashShopFailure(reader.ReadU8())
	case constant.CashShopResultBought, constant.CashShopResultPutIn:
		p.Item = &dto.CashShopItem{}
		p.Item.Deserialize(reader)
	case constant.CashShopResultCouponRedeemed:
		p.Items = p.deserializeItems(reader, int(reader.ReadU8()))
		p.MaplePoint = reader.ReadU32()
		count := reader.ReadU32()
		for range count {
			granted := &dto.CashShopGrantedItem{}
			granted.Deserialize(reader)
			p.Granted = append(p.Granted, granted)
		}
		p.Meso = reader.ReadU32()
	case constant.CashShopResultGiftSent:
		p.Recipient = reader.ReadStr16()
		p.ItemID = reader.ReadU32()
		p.Count = reader.ReadU16()
	case constant.CashShopResultInventorySlots:
		p.InventoryType = reader.ReadU8()
		p.Slots = reader.ReadU16()
	case constant.CashShopResultStorageSlots, constant.CashShopResultCharacterSlots:
		p.Slots = reader.ReadU16()
	case constant.CashShopResultTakenOut:
		p.Slot = int16(reader.ReadU16())
		p.TakenOut = dto.NewItemFromStream(reader)
	case constant.CashShopResultExpired:
		p.Serial = reader.ReadU64()
	case constant.CashShopResultPaidBack:
		p.Serial = reader.ReadU64()
		p.MaplePoint = reader.ReadU32()
	case constant.CashShopResultPackageBought:
		p.Items = p.deserializeItems(reader, int(reader.ReadU8()))
		reader.ReadU16()
	case constant.CashShopResultQuestItemBought:
		count := reader.ReadU32()
		for range count {
			granted := &dto.CashShopGrantedItem{}
			granted.Deserialize(reader)
			p.Granted = append(p.Granted, granted)
		}
	}
}

func (p *CashShopResult) deserializeItems(reader *stream.StreamReader, count int) []*dto.CashShopItem {
	items := make([]*dto.CashShopItem, 0, count)
	for range count {
		item := &dto.CashShopItem{}
		item.Deserialize(reader)
		items = append(items, item)
	}
	return items
}
