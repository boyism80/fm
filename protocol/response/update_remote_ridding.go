package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type UpdateRemoteRidding struct {
	CharacterID int32
	MountID     int32
	Buffs       []dto.BuffEntry
}

func (p *UpdateRemoteRidding) Opcode() uint16 { return 0x90 }

func (p *UpdateRemoteRidding) Serialize(writer *stream.StreamWriter) error {
	writer.Write32(p.CharacterID)

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

	writer.WriteU16(0)
	writer.Write32(p.MountID)
	writer.Write32(int32(constant.SkillMonsterRider))
	writer.Write32(0)
	writer.Write32(0)
	writer.WriteU8(0)
	return nil
}

func (p *UpdateRemoteRidding) Deserialize(_ *stream.StreamReader) {}
