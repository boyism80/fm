package response

import (
	"time"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

// UpdateBuff notifies the client to apply buff(s) to self (buff on).
type UpdateBuff struct {
	BuffID   int32
	Duration time.Duration
	Buffs    []dto.BuffEntry
}

func (p *UpdateBuff) Opcode() uint16 { return 0x15 }

func (p *UpdateBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []dto.BuffEntry{}
	}
	buffs := make([]dto.BuffEntry, len(p.Buffs))
	copy(buffs, p.Buffs)
	sortBuffs(buffs)

	if err := writeBuffMask(writer, buffs); err != nil {
		return err
	}
	for _, buff := range buffs {
		writer.WriteU16(uint16(buff.Value))
		writer.Write32(p.BuffID)
		writer.Write32(int32(p.Duration.Milliseconds()))
	}
	writer.WriteU16(0)
	writer.WriteU16(0)
	writer.Write32(0)
	writer.WriteU8(0)
	return nil
}

func (p *UpdateBuff) Deserialize(_ *stream.StreamReader) error { return nil }
