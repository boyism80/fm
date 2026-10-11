package response

import (
	"github.com/boyism80/fm/stream"
)

type ApplyMobBuff struct {
	OID        uint32
	Status     int32
	X          int16
	SkillID    uint16
	SkillLevel uint16
	BuffTime   int16
	Delay      int16
	StatusSize byte
}

func (p *ApplyMobBuff) Opcode() uint16 {
	return 0xAF
}

func (p *ApplyMobBuff) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(p.Status)
	writer.Write16(p.X)
	writer.WriteU16(p.SkillID)
	writer.WriteU16(p.SkillLevel)
	writer.Write16(p.BuffTime)
	writer.Write16(p.Delay)
	writer.WriteU8(p.StatusSize)
	return nil
}

func (p *ApplyMobBuff) Deserialize(reader *stream.StreamReader) {
	p.OID = reader.ReadU32()
	p.Status = reader.Read32()
	p.X = reader.Read16()
	p.SkillID = reader.ReadU16()
	p.SkillLevel = reader.ReadU16()
	p.BuffTime = reader.Read16()
	p.Delay = reader.Read16()
	p.StatusSize = reader.ReadU8()
}
