package resp

import "github.com/boyism80/fm/common/stream"

type ControlMoveMob struct {
	OID          uint32
	MoveId       uint16
	EnabledSkill bool
	MP           uint16
	SkillId      uint8
	SkillLevel   uint8
}

func (p *ControlMoveMob) Opcode() uint16 {
	return 0xAD
}

func (p *ControlMoveMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.WriteU16(p.MoveId)
	writer.WriteBoolean(p.EnabledSkill)
	writer.WriteU16(p.MP)
	writer.WriteU8(p.SkillId)
	writer.WriteU8(p.SkillLevel)
	return nil
}

func (p *ControlMoveMob) Deserialize(reader *stream.StreamReader) error {
	return nil
}
