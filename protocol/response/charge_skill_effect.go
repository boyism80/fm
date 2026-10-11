package response

import "github.com/boyism80/fm/stream"

type ChargeSkillEffect struct {
	CharacterID uint32
	SkillID     uint32
	Level       uint8
	Flags       uint8
	Speed       uint8
}

func (*ChargeSkillEffect) Opcode() uint16 { return 0x87 }
func (p *ChargeSkillEffect) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.SkillID)
	w.WriteU8(p.Level)
	w.WriteU8(p.Flags)
	w.WriteU8(p.Speed)
	return nil
}
func (p *ChargeSkillEffect) Deserialize(r *stream.StreamReader) {
	p.CharacterID = r.ReadU32()
	p.SkillID = r.ReadU32()
	p.Level = r.ReadU8()
	p.Flags = r.ReadU8()
	p.Speed = r.ReadU8()
}
