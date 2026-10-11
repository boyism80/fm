package request

import "github.com/boyism80/fm/stream"

type SkillMacroEntry struct {
	Name   string
	Shout  uint8
	Skills [3]uint32
}

type SkillMacro struct {
	Count   uint8
	Entries []SkillMacroEntry
}

func (*SkillMacro) Opcode() byte { return 0x5D }

func (p *SkillMacro) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(len(p.Entries)))
	for _, entry := range p.Entries {
		w.WriteStr16(entry.Name)
		w.WriteU8(entry.Shout)
		for _, skillID := range entry.Skills {
			w.WriteU32(skillID)
		}
	}
	return nil
}

func (p *SkillMacro) Deserialize(r *stream.StreamReader) {
	p.Count = r.ReadU8()
	if p.Count > 5 {
		return
	}
	p.Entries = make([]SkillMacroEntry, p.Count)
	for i := range p.Entries {
		p.Entries[i].Name = r.ReadStr16()
		p.Entries[i].Shout = r.ReadU8()
		for j := range p.Entries[i].Skills {
			p.Entries[i].Skills[j] = r.ReadU32()
		}
	}
}
