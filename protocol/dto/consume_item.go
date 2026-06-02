package dto

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type ConsumeItem struct {
	ItemId     uint32
	Count      uint16
	Expiration time.Time
	OwnerName  string
	Flags      uint16
}

func (i *ConsumeItem) GetCount() uint16 {
	return i.Count
}

func (i *ConsumeItem) Serialize(writer *stream.StreamWriter, opt ItemSerializeOption) {
	if opt.SlotMode != SlotEncodeOmit {
		writer.WriteU8(uint8(opt.Slot))
	}
	writer.WriteU8(uint8(constant.ItemTypeETC))
	writer.WriteU32(i.ItemId)

	writer.WriteBoolean(false)

	writer.WriteDateTime(i.Expiration)
	writer.WriteU16(i.Count)
	writer.WriteStr16(i.OwnerName)
	writer.WriteU16(i.Flags)

	switch constant.GetConsumeType(i.ItemId) {
	case constant.ConsumeTypeShuriken, constant.ConsumeTypeBullet:
		writer.WriteU64(54399043)
	}
}
