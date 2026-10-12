package request

import "github.com/boyism80/fm/stream"

type FamilyRemoveJunior struct {
	JuniorID uint32
}

func (*FamilyRemoveJunior) Opcode() byte { return 0x7d }

func (p *FamilyRemoveJunior) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.JuniorID)
	return nil
}

func (p *FamilyRemoveJunior) Deserialize(r *stream.StreamReader) {
	p.JuniorID = r.ReadU32()
}
