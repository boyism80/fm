package dto

import (
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

// Equipment represents equipment data for protocol
type Equipment struct {
	ItemId        uint32
	UniqueId      int64
	Expiration    time.Time // Time.Time for proper serialization
	EnchantChance uint8
	Level         uint8
	Str           uint16
	Dex           uint16
	Int           uint16
	Luk           uint16
	MaxHP         uint16
	MaxMP         uint16
	PAD           uint16
	MAD           uint16
	PDD           uint16
	MDD           uint16
	ACC           uint16
	Avoid         uint16
	Hands         uint16
	Speed         uint16
	Jump          uint16
	OwnerName     string
	Flag          uint16
	SkillBonus    uint8
}

// GetCount returns the equipment count (always 1 for equipment)
func (e *Equipment) GetCount() uint16 {
	return 1
}

// Serialize serializes equipment data
func (e *Equipment) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	if slot <= -1 {
		slot *= -1
		if slot > 100 && slot < 1000 {
			slot -= 100
		}
	}
	if slot != 0 && !trade {
		writer.WriteU16(uint16(slot))
	} else {
		writer.WriteU8(uint8(slot))
	}

	writer.WriteU8(1) // ITEM_TYPE_EQUIPMENT
	writer.WriteU32(e.ItemId)

	hasUID := e.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(e.UniqueId)
	}

	writer.WriteDateTime(e.Expiration)
	writer.WriteU8(e.EnchantChance)
	writer.WriteU8(e.Level)
	writer.WriteU16(e.Str)
	writer.WriteU16(e.Dex)
	writer.WriteU16(e.Int)
	writer.WriteU16(e.Luk)
	writer.WriteU16(e.MaxHP)
	writer.WriteU16(e.MaxMP)
	writer.WriteU16(e.PAD)
	writer.WriteU16(e.MAD)
	writer.WriteU16(e.PDD)
	writer.WriteU16(e.MDD)
	writer.WriteU16(e.ACC)
	writer.WriteU16(e.Avoid)
	writer.WriteU16(e.Hands)
	writer.WriteU16(e.Speed)
	writer.WriteU16(e.Jump)
	writer.WriteStr16(e.OwnerName)
	writer.WriteU16(e.Flag)
	writer.WriteBoolean(e.SkillBonus > 0)
	writer.WriteU8(1)
	writer.WriteU32(0)
	if e.UniqueId <= 0 {
		writer.Write64(-1)
	}
	writer.WriteDateTime(util.TimeZero)
	writer.Write32(-1)
}

