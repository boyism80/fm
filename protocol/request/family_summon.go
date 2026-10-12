package request

import "github.com/boyism80/fm/stream"

type FamilySummon struct {
	Name     string
	Accepted bool
}

func (*FamilySummon) Opcode() byte { return 0x82 }

func (p *FamilySummon) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	w.WriteBoolean(p.Accepted)
	return nil
}

func (p *FamilySummon) Deserialize(r *stream.StreamReader) {
	p.Name = r.ReadStr16()
	p.Accepted = r.ReadBool()
}
