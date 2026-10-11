package response

import "github.com/boyism80/fm/stream"

type UseSkillBook struct {
	CharacterID uint32
	SkillID     uint32
	MasterLevel uint32
	CanUse      bool
	Success     bool
}

func (*UseSkillBook) Opcode() uint16 { return 0x27 }

func (p *UseSkillBook) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(1)
	writer.WriteU32(p.SkillID)
	writer.WriteU32(p.MasterLevel)
	writer.WriteBoolean(p.CanUse)
	writer.WriteBoolean(p.Success)
	return nil
}

func (p *UseSkillBook) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	reader.ReadU8()
	p.SkillID = reader.ReadU32()
	p.MasterLevel = reader.ReadU32()
	p.CanUse = reader.ReadBool()
	p.Success = reader.ReadBool()
}
