package request

import "github.com/boyism80/fm/stream"

type FamilyPrecept struct {
	Precept string
}

func (*FamilyPrecept) Opcode() byte { return 0x81 }

func (p *FamilyPrecept) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Precept)
	return nil
}

func (p *FamilyPrecept) Deserialize(r *stream.StreamReader) {
	p.Precept = r.ReadStr16()
}
