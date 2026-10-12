package request

import "github.com/boyism80/fm/stream"

type FamilyAccept struct {
	SeniorID uint32
	Name     string
	Accepted bool
}

func (*FamilyAccept) Opcode() byte { return 0x7f }

func (p *FamilyAccept) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.SeniorID)
	w.WriteStr16(p.Name)
	w.WriteBoolean(p.Accepted)
	return nil
}

func (p *FamilyAccept) Deserialize(r *stream.StreamReader) {
	p.SeniorID = r.ReadU32()
	p.Name = r.ReadStr16()
	p.Accepted = r.ReadBool()
}
