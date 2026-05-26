package response

import (
	"github.com/boyism80/fm/stream"
)

type ApplyMobBuffReflect struct {
	OID         uint32
	StatusMask  int32
	Entries     []ApplyMobBuffReflectEntry
	Reflections []int32
	Delay       int16
	StatusSize  byte
}

type ApplyMobBuffReflectEntry struct {
	X          int16
	SkillID    uint16
	SkillLevel uint16
	BuffTime   int16
}

func (p *ApplyMobBuffReflect) Opcode() uint16 {
	return 0xAF
}

func (p *ApplyMobBuffReflect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(p.StatusMask)
	for _, e := range p.Entries {
		writer.Write16(e.X)
		writer.WriteU16(e.SkillID)
		writer.WriteU16(e.SkillLevel)
		writer.Write16(e.BuffTime)
	}
	for _, r := range p.Reflections {
		writer.Write32(r)
	}
	writer.Write16(p.Delay)
	writer.WriteU8(p.StatusSize)
	return nil
}

func (p *ApplyMobBuffReflect) Deserialize(reader *stream.StreamReader) {
}
