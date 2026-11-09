package response

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type UpdateSelfRidding struct {
	BuffID  int32
	MountID int32
	Buffs   []dto.BuffEntry
}

func (p *UpdateSelfRidding) Opcode() uint16 { return 0x15 }

func (p *UpdateSelfRidding) Serialize(writer *stream.StreamWriter) error {
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
	// Every stat of the ride, the ride itself included, carries the mount and the skill where a buff carries
	// its source and duration.
	for i := range buffs {
		value := buffs[i].Value
		if buffs[i].Buff == constant.BuffFlagMonsterRiding {
			value = 0
		}
		writer.WriteU16(uint16(value))
		writer.Write32(p.MountID)
		writer.Write32(p.BuffID)
	}
	writer.WriteU16(0)
	writer.WriteU16(0)
	writer.Write32(0)
	writer.WriteU8(0)
	return nil
}

func (p *UpdateSelfRidding) Deserialize(_ *stream.StreamReader) {}
