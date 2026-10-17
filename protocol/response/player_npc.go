package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type PlayerNpc struct {
	NpcID uint32
	Look  *dto.Character
}

func (p *PlayerNpc) Opcode() uint16 {
	return 0x40
}

func (p *PlayerNpc) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteU32(p.NpcID)
	writer.WriteStr16(p.Look.Name)
	p.Look.SerializeLook(writer)
	return nil
}

func (p *PlayerNpc) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.NpcID = reader.ReadU32()
	p.Look = &dto.Character{Name: reader.ReadStr16()}
	p.Look.DeserializeLook(reader)
}
