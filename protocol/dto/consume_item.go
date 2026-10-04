package dto

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
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

func (i *ConsumeItem) GetItemID() uint32 {
	return i.ItemId
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

func (i *ConsumeItem) Deserialize(reader *stream.StreamReader) {
	if reader.ReadBool() {
		reader.Skip(8)
	}
	i.Expiration = util.FromFileTime(reader.ReadU64())
	i.Count = reader.ReadU16()
	i.OwnerName = reader.ReadStr16()
	i.Flags = reader.ReadU16()
	if constant.ConsumeNeedsAmmoInventoryID(i.ItemId) {
		reader.Skip(8)
	}
}
