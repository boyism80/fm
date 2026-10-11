package request

import "github.com/boyism80/fm/stream"

type SkillEffect struct {
	SkillID uint32
	Level   uint8
	Flags   uint8
	Speed   uint8
}

func (*SkillEffect) Opcode() byte { return 0x4C }
func (p *SkillEffect) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.SkillID)
	w.WriteU8(p.Level)
	w.WriteU8(p.Flags)
	w.WriteU8(p.Speed)
	return nil
}
func (p *SkillEffect) Deserialize(r *stream.StreamReader) {
	p.SkillID = r.ReadU32()
	p.Level = r.ReadU8()
	p.Flags = r.ReadU8()
	p.Speed = r.ReadU8()
}
