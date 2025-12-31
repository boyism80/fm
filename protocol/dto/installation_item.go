package dto

import (
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// InstallationItem represents installation item DTO
type InstallationItem struct {
	ItemId      uint32
	UniqueId    int64
	Expiration  time.Time
	OwnerName   string
	Flags       uint16
}

func (i *InstallationItem) GetCount() uint16 {
	return 1
}

func (i *InstallationItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(constant.ITEM_TYPE_ETC))
	writer.WriteU32(i.ItemId)

	hasUID := i.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(i.UniqueId)
	}

	writer.WriteDateTime(i.Expiration)
	writer.WriteU16(1) // Installation always has count 1
	writer.WriteStr16(i.OwnerName)
	writer.WriteU16(i.Flags)
}

