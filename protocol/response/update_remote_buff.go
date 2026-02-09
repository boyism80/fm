package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

// UpdateRemoteBuff notifies other clients on the map that a character gained buff(s).
// "Remote" = another character (observed by others). Broadcast to map for display.
type UpdateRemoteBuff struct {
	CharacterID int32
	Buffs       []dto.BuffEntry
}

func (p *UpdateRemoteBuff) Opcode() uint16 { return 0x5E }

func (p *UpdateRemoteBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []dto.BuffEntry{}
	}
	buffs := make([]dto.BuffEntry, len(p.Buffs))
	copy(buffs, p.Buffs)
	sortBuffs(buffs)

	writer.Write32(p.CharacterID)
	if err := writeBuffMask(writer, buffs); err != nil {
		return err
	}
	for _, e := range buffs {
		writer.WriteU16(uint16(e.Value))
	}
	writer.WriteU16(0)
	writer.WriteU16(0)
	return nil
}

func (p *UpdateRemoteBuff) Deserialize(_ *stream.StreamReader) error { return nil }
