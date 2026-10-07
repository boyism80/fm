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
	Item           *dto.CashShopItem
	Slot           int16
	TakenOut       dto.Item
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
		writer.WriteU16(0)
	case constant.CashShopResultWishlist, constant.CashShopResultWishlistUpdated:
		for i := range 10 {
			if i < len(p.Wishlist) {
				writer.WriteU32(p.Wishlist[i])
			} else {
				writer.WriteU32(0)
			}
		}
	case constant.CashShopResultBuyFailed, constant.CashShopResultTakeOutFailed, constant.CashShopResultPutInFailed:
		writer.WriteU8(uint8(p.Failure))
	case constant.CashShopResultBought, constant.CashShopResultPutIn:
		p.Item.Serialize(writer)
	case constant.CashShopResultTakenOut:
		writer.WriteU16(uint16(p.Slot))
		p.TakenOut.Serialize(writer, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
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
	case constant.CashShopResultWishlist, constant.CashShopResultWishlistUpdated:
		p.Wishlist = make([]uint32, 10)
		for i := range p.Wishlist {
			p.Wishlist[i] = reader.ReadU32()
		}
	case constant.CashShopResultBuyFailed, constant.CashShopResultTakeOutFailed, constant.CashShopResultPutInFailed:
		p.Failure = constant.CashShopFailure(reader.ReadU8())
	case constant.CashShopResultBought, constant.CashShopResultPutIn:
		p.Item = &dto.CashShopItem{}
		p.Item.Deserialize(reader)
	case constant.CashShopResultTakenOut:
		p.Slot = int16(reader.ReadU16())
		p.TakenOut = dto.NewItemFromStream(reader)
	}
}
