package entity

import (
	"github.com/boyism80/fm/common/stream"
)

type CashItem struct {
	*baseItem
	Count     uint16
	OwnerName string
	Flags     uint16
}

func (item *CashItem) GetInventoryType() InventoryType {
	return InventoryTypeCash
}

func (item *CashItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(ItemTypeEtc))
	writer.WriteU32(item.Template.GetID())

	hasUID := item.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(item.UniqueId)
	}

	writer.WriteDateTime(item.Expiration)
	writer.WriteU16(item.Count)
	writer.WriteStr16(item.OwnerName)
	writer.WriteU16(item.Flags)
}
