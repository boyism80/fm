package response

import "github.com/boyism80/fm/stream"

type Chalkboard struct {
	CharacterID uint32
	Text        string
}

func (*Chalkboard) Opcode() uint16 { return 0x71 }
func (p *Chalkboard) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteBoolean(p.Text != "")
	if p.Text != "" {
		w.WriteStr16(p.Text)
	}
	return nil
}
func (p *Chalkboard) Deserialize(r *stream.StreamReader) {
	p.CharacterID = r.ReadU32()
	if r.ReadU8() != 0 {
		p.Text = r.ReadStr16()
	}
}
