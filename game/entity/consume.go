package entity

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/data"
)

type Consume struct {
	*BaseItem
	Count     uint16
	OwnerName string
	Flags     uint16
}

func (item *Consume) GetInventoryType() InventoryType {
	return InventoryTypeConsume
}

func (item *Consume) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(ItemTypeEtc))
	writer.WriteU32(item.Template.GetID())

	hasUID := item.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(item.UniqueId)
	}

	writer.WriteDateTime(item.Expiration)
	template, ok := item.Template.(*data.ConsumeTemplate)
	if !ok {
		return // TODO: return error
	}

	writer.WriteU16(item.Count)
	writer.WriteStr16(item.OwnerName)
	writer.WriteU16(item.Flags)

	isThrowingStart := template.Id/10000 == 207
	isBullet := template.Id/10000 == 233
	isWhat := template.Id/10000 == 287
	inventoryId := uint64(54399043) // ??
	if isThrowingStart || isBullet || isWhat {
		writer.WriteU64(inventoryId)
	}
}
