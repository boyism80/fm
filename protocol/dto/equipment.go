package dto

import (
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type Equipment struct {
	ItemId        uint32
	UniqueId      *uint64
	Expiration    time.Time
	EnhanceChance uint8
	EnhanceCount  uint8
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

func (e *Equipment) GetCount() uint16 {
	return 1
}

func (e *Equipment) Serialize(writer *stream.StreamWriter, opt ItemSerializeOption) {
	slot := opt.Slot
	if slot < 0 {
		slot = -slot
		if slot > 100 && slot < 1000 {
			slot -= 100
		}
	}
	switch opt.SlotMode {
	case SlotEncodeActual:
		if slot != 0 && !opt.Trade {
			writer.Write16(slot)
			break
		}
		writer.WriteU8(uint8(slot))
	case SlotEncodeZero:
		writer.WriteU8(uint8(slot))
	}

	writer.WriteU8(1)
	writer.WriteU32(e.ItemId)

	hasUID := e.UniqueId != nil
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.WriteU64(*e.UniqueId)
	}

	writer.WriteDateTime(e.Expiration)
	writer.WriteU8(e.EnhanceChance)
	writer.WriteU8(e.EnhanceCount)
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
	if e.UniqueId == nil {
		writer.Write64(-1)
	}
	writer.WriteDateTime(util.TimeZero)
	writer.Write32(-1)
}
