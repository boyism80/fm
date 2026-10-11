package response

import "github.com/boyism80/fm/stream"

type SkillMacroEntry struct {
	Name   string
	Shout  bool
	Skills [3]uint32
}
type SkillMacro struct{ Entries []SkillMacroEntry }

func (*SkillMacro) Opcode() uint16 { return 0x54 }
func (p *SkillMacro) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(len(p.Entries)))
	for _, entry := range p.Entries {
		w.WriteStr16(entry.Name)
		w.WriteBoolean(entry.Shout)
		for _, skillID := range entry.Skills {
			w.WriteU32(skillID)
		}
	}
	return nil
}
func (p *SkillMacro) Deserialize(r *stream.StreamReader) {
	count := r.ReadU8()
	p.Entries = make([]SkillMacroEntry, count)
	for i := range p.Entries {
		p.Entries[i].Name = r.ReadStr16()
		p.Entries[i].Shout = r.ReadU8() != 0
		for j := range p.Entries[i].Skills {
			p.Entries[i].Skills[j] = r.ReadU32()
		}
	}
}
