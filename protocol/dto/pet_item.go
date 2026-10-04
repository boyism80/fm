package dto

import (
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type PetItem struct {
	ItemId         uint32
	UniqueId       *uint64
	Expiration     time.Time
	PetName        string
	PetLevel       uint8
	PetCloseness   uint16
	PetFullness    uint8
	PetSpeed       uint16
	PetFlags       uint16
	PetExpiration  time.Time
	PetSecondsLeft uint32
}

func (i *PetItem) GetCount() uint16 {
	return 1
}

func (i *PetItem) GetItemID() uint32 {
	return i.ItemId
}

func (i *PetItem) Serialize(writer *stream.StreamWriter, opt ItemSerializeOption) {
	if opt.SlotMode != SlotEncodeOmit {
		writer.WriteU8(uint8(opt.Slot))
	}
	writer.WriteU8(3)
	writer.WriteU32(i.ItemId)

	hasUID := i.UniqueId != nil
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.WriteU64(*i.UniqueId)
	}

	writer.WriteDateTime(i.Expiration)
	writer.WriteStaticStr(i.PetName, 13)
	writer.WriteU8(i.PetLevel)
	writer.WriteU16(i.PetCloseness)
	writer.WriteU8(i.PetFullness)
	writer.WriteDateTime(i.PetExpiration)
	writer.WriteU16(i.PetSpeed)
	writer.WriteU16(i.PetFlags)
	if i.ItemId == 5000054 && i.PetSecondsLeft > 0 {
		writer.WriteU32(i.PetSecondsLeft)
	} else {
		writer.WriteU32(0)
	}
}

func (i *PetItem) Deserialize(reader *stream.StreamReader) {
	if reader.ReadBool() {
		uid := reader.ReadU64()
		i.UniqueId = &uid
	}
	i.Expiration = util.FromFileTime(reader.ReadU64())
	i.PetName = reader.ReadStaticStr(13)
	i.PetLevel = reader.ReadU8()
	i.PetCloseness = reader.ReadU16()
	i.PetFullness = reader.ReadU8()
	i.PetExpiration = util.FromFileTime(reader.ReadU64())
	i.PetSpeed = reader.ReadU16()
	i.PetFlags = reader.ReadU16()
	i.PetSecondsLeft = reader.ReadU32()
}
