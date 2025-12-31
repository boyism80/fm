package dto

import (
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// GeneralItem represents general item DTO
type GeneralItem struct {
	ItemId      uint32
	UniqueId    int64
	Count       uint16
	Expiration  time.Time
	OwnerName   string
	Flags       uint16
}

func (i *GeneralItem) GetCount() uint16 {
	return i.Count
}

func (i *GeneralItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(constant.ITEM_TYPE_ETC))
	writer.WriteU32(i.ItemId)

	hasUID := i.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(i.UniqueId)
	}

	writer.WriteDateTime(i.Expiration)
	writer.WriteU16(i.Count)
	writer.WriteStr16(i.OwnerName)
	writer.WriteU16(i.Flags)
}

