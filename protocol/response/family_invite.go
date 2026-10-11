package response

import "github.com/boyism80/fm/stream"

type FamilyInvite struct {
	InviterID uint32
	Name      string
}

func (p *FamilyInvite) Opcode() uint16 {
	return 0x4C
}

func (p *FamilyInvite) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.InviterID)
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyInvite) Deserialize(reader *stream.StreamReader) {
	p.InviterID = reader.ReadU32()
	p.Name = reader.ReadStr16()
}
