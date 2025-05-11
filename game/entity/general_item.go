package entity

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/data"
)

type GeneralItem struct {
	*BaseItem
	Count     uint16
	OwnerName string
	Flags     uint16
}

func (item *GeneralItem) Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade bool, slot int16, itemType ItemType) {
	if zeroPosition {
		if !leaveOut {
			writer.WriteU8(0)
		}
	} else {
		if slot <= -1 {
			slot *= -1
			if slot > 100 && slot < 1000 {
				slot -= 100
			}
		}
		if !trade && itemType == ItemTypeEquipment { // equipment
			writer.WriteU16(uint16(slot))
		} else {
			writer.WriteU8(uint8(slot))
		}
	}

	writer.WriteU8(uint8(itemType))
	writer.WriteU32(item.Template.GetID())

	hasUID := item.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(item.UniqueId)
	}

	writer.WriteDateTime(item.Expiration)
	template, ok := item.Template.(*data.GeneralItemTemplate)
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
