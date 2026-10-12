package request

import "github.com/boyism80/fm/stream"

type FamilyInvite struct {
	Name string
}

func (*FamilyInvite) Opcode() byte { return 0x7c }

func (p *FamilyInvite) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyInvite) Deserialize(r *stream.StreamReader) {
	p.Name = r.ReadStr16()
}
