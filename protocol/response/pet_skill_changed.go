package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type PetSkillChanged struct {
	SN    uint64
	Add   bool
	Skill constant.PetSkill
}

func (p *PetSkillChanged) Opcode() uint16 {
	return 0x9E
}

func (p *PetSkillChanged) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU64(p.SN)
	writer.WriteBoolean(p.Add)
	writer.WriteU16(uint16(p.Skill))
	return nil
}

func (p *PetSkillChanged) Deserialize(reader *stream.StreamReader) {
	p.SN = reader.ReadU64()
	p.Add = reader.ReadBool()
	p.Skill = constant.PetSkill(reader.ReadU16())
}
