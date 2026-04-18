package dto

import (
	"time"

	"github.com/boyism80/fm/stream"
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

func (i *PetItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
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
