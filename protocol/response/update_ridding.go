package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type UpdateRidding struct {
	BuffID  int32
	MountID int32
	Buffs   []dto.BuffEntry
}

func (p *UpdateRidding) Opcode() uint16 { return 0x15 }

func (p *UpdateRidding) Serialize(writer *stream.StreamWriter) error {
	buffs := p.Buffs
	if buffs == nil {
		buffs = []dto.BuffEntry{}
	}
	SortBuffEntries(buffs)
	flags := make([]constant.BuffFlag, 0, len(buffs))
	for i := range buffs {
		flags = append(flags, buffs[i].Buff)
	}
	WriteBuffs(writer, flags)

	ridingLevel := max(int32(1), p.MountID-1902000+1)

	writer.WriteU16(uint16(ridingLevel))
	writer.Write32(p.MountID)
	writer.Write32(p.BuffID)
	writer.WriteU16(0)
	writer.WriteU16(0)
	writer.Write32(0)
	writer.WriteU8(0)
	return nil
}

func (p *UpdateRidding) Deserialize(_ *stream.StreamReader) {}
