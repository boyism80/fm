package request

import (
	"github.com/boyism80/fm/stream"
)

type PetReviveInquiry struct {
	Tick uint32
	SN   uint64
}

func (*PetReviveInquiry) Opcode() byte { return 0x3F }

func (p *PetReviveInquiry) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.WriteU64(p.SN)
	return nil
}

func (p *PetReviveInquiry) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.SN = reader.ReadU64()
}
