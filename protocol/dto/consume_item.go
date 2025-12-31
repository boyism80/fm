package dto

import (
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// ConsumeItem represents consume item DTO
type ConsumeItem struct {
	ItemId          uint32
	UniqueId        int64
	Count           uint16
	Expiration      time.Time
	OwnerName       string
	Flags           uint16
	IsThrowingStart bool
	IsBullet        bool
	IsWhat          bool
}

func (i *ConsumeItem) GetCount() uint16 {
	return i.Count
}

func (i *ConsumeItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
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

	// Special case: inventoryId for certain item types
	inventoryId := uint64(54399043)
	if i.IsThrowingStart || i.IsBullet || i.IsWhat {
		writer.WriteU64(inventoryId)
	}
}
