package request

import "github.com/boyism80/fm/stream"

type FamilyPedigree struct {
	Name string
}

func (*FamilyPedigree) Opcode() byte { return 0x7a }

func (p *FamilyPedigree) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyPedigree) Deserialize(r *stream.StreamReader) {
	p.Name = r.ReadStr16()
}
