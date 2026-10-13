package response

import "github.com/boyism80/fm/stream"

type CancelChargeSkillEffect struct {
	CharacterID uint32
	SkillID     uint32
}

func (*CancelChargeSkillEffect) Opcode() uint16 { return 0x88 }
func (p *CancelChargeSkillEffect) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.SkillID)
	return nil
}
func (p *CancelChargeSkillEffect) Deserialize(r *stream.StreamReader) {
	p.CharacterID = r.ReadU32()
	p.SkillID = r.ReadU32()
}
