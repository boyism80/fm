package dto

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type MiscItem struct {
	ItemId     uint32
	Count      uint16
	Expiration time.Time
	OwnerName  string
	Flags      uint16
}

func (i *MiscItem) GetCount() uint16 {
	return i.Count
}

func (i *MiscItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(constant.ITEM_TYPE_ETC))
	writer.WriteU32(i.ItemId)

	writer.WriteBoolean(false)

	writer.WriteDateTime(i.Expiration)
	writer.WriteU16(i.Count)
	writer.WriteStr16(i.OwnerName)
	writer.WriteU16(i.Flags)
}
