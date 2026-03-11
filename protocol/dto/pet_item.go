package dto

import (
	"time"

	"github.com/boyism80/fm/stream"
)

// PetItem represents pet item DTO
type PetItem struct {
	ItemId         uint32
	UniqueId       int64
	Expiration     time.Time // ItemCore.Expiration
	PetName        string
	PetLevel       uint8
	PetCloseness   uint16
	PetFullness    uint8
	PetSpeed       uint16
	PetFlags       uint16
	PetExpiration  time.Time // Pet.Expiration
	PetSecondsLeft uint32
}

func (i *PetItem) GetCount() uint16 {
	return 1
}

func (i *PetItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(3) // ITEM_TYPE_PET
	writer.WriteU32(i.ItemId)

	hasUID := i.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(i.UniqueId)
	}

	writer.WriteDateTime(i.Expiration) // First expiration (ItemCore.Expiration)
	writer.WriteStaticStr(i.PetName, 13)
	writer.WriteU8(i.PetLevel)
	writer.WriteU16(i.PetCloseness)
	writer.WriteU8(i.PetFullness)
	writer.WriteDateTime(i.PetExpiration) // Pet expiration (second one, Pet.Expiration)
	writer.WriteU16(i.PetSpeed)
	writer.WriteU16(i.PetFlags)
	if i.ItemId == 5000054 && i.PetSecondsLeft > 0 {
		writer.WriteU32(i.PetSecondsLeft)
	} else {
		writer.WriteU32(0)
	}
}
