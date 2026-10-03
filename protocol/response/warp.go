package response

import (
	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type Warp struct {
	Character *dto.Character
	Channel   uint32
}

func (p *Warp) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Channel)
	writer.WriteU16(2)
	writer.WriteU16(0)
	writer.WriteU32(p.Character.Map)
	writer.WriteU8(p.Character.SpawnPoint)
	writer.WriteU16(p.Character.Hp)
	writer.WriteDateTime(clock.Now())
	return nil
}

func (a *Warp) Deserialize(reader *stream.StreamReader) {
	a.Channel = reader.ReadU32()
	reader.ReadU16()
	reader.ReadU16()
	if a.Character == nil {
		a.Character = &dto.Character{}
	}
	a.Character.Map = reader.ReadU32()
	a.Character.SpawnPoint = reader.ReadU8()
	a.Character.Hp = reader.ReadU16()
	reader.Skip(8)
}

func (w *Warp) Opcode() uint16 {
	return 0x55
}
