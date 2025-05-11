package entity

import (
	"github.com/boyism80/fm/common/stream"
)

type GeneralItem struct {
	*BaseItem
	Count     uint16
	OwnerName string
	Flags     uint16
}

func (item *GeneralItem) GetInventoryType() InventoryType {
	return InventoryTypeEtc
}

func (item *GeneralItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
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
