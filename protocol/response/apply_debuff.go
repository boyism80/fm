package response

import (
	"github.com/boyism80/fm/stream"
)

type ApplyDebuff struct {
	OID        uint32
	Status     int32
	X          int16
	SkillID    uint32
	BuffTime   int16
	Delay      int16
	StatusSize byte
}

func (p *ApplyDebuff) Opcode() uint16 {
	return 0xAF
}

func (p *ApplyDebuff) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(p.Status)
	writer.Write16(p.X)
	if p.SkillID > 0 {
		writer.WriteU32(p.SkillID)
	}
	writer.Write16(p.BuffTime)
	writer.Write16(p.Delay)
	writer.WriteU8(p.StatusSize)
	return nil
}

func (p *ApplyDebuff) Deserialize(reader *stream.StreamReader) error {
	return nil
}
