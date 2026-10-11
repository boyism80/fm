package response

import "github.com/boyism80/fm/stream"

type FamilyInviteResult struct {
	Accepted bool
	Name     string
}

func (p *FamilyInviteResult) Opcode() uint16 {
	return 0x4D
}

func (p *FamilyInviteResult) Serialize(w *stream.StreamWriter) error {
	w.WriteBoolean(p.Accepted)
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyInviteResult) Deserialize(reader *stream.StreamReader) {
	p.Accepted = reader.ReadBool()
	p.Name = reader.ReadStr16()
}
