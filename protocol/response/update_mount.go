package response

import (
	"github.com/boyism80/fm/stream"
)

type UpdateMount struct {
	CharacterID uint32
	Level       uint32
	Exp         uint32
	Fatigue     uint32
	LevelUp     bool
}

func (p *UpdateMount) Opcode() uint16 {
	return 0x24
}

func (p *UpdateMount) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU32(p.Level)
	writer.WriteU32(p.Exp)
	writer.WriteU32(p.Fatigue)
	writer.WriteBoolean(p.LevelUp)
	return nil
}

func (p *UpdateMount) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.Level = reader.ReadU32()
	p.Exp = reader.ReadU32()
	p.Fatigue = reader.ReadU32()
	p.LevelUp = reader.ReadBool()
}
