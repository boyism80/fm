package dto

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type InstallationItem struct {
	ItemId     uint32
	Expiration time.Time
	OwnerName  string
	Flags      uint16
}

func (i *InstallationItem) GetCount() uint16 {
	return 1
}

func (i *InstallationItem) GetItemID() uint32 {
	return i.ItemId
}

func (i *InstallationItem) Serialize(writer *stream.StreamWriter, opt ItemSerializeOption) {
	if opt.SlotMode != SlotEncodeOmit {
		writer.WriteU8(uint8(opt.Slot))
	}
	writer.WriteU8(uint8(constant.ItemTypeETC))
	writer.WriteU32(i.ItemId)

	writer.WriteBoolean(false)

	writer.WriteDateTime(i.Expiration)
	writer.WriteU16(1)
	writer.WriteStr16(i.OwnerName)
	writer.WriteU16(i.Flags)
}

func (i *InstallationItem) Deserialize(reader *stream.StreamReader) {
	if reader.ReadBool() {
		reader.Skip(8)
	}
	i.Expiration = util.FromFileTime(reader.ReadU64())
	reader.ReadU16()
	i.OwnerName = reader.ReadStr16()
	i.Flags = reader.ReadU16()
}
