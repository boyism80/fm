package response

import "github.com/boyism80/fm/stream"

type FamilyJoined struct {
	Name string
}

func (p *FamilyJoined) Opcode() uint16 {
	return 0x4E
}

func (p *FamilyJoined) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyJoined) Deserialize(reader *stream.StreamReader) {
	p.Name = reader.ReadStr16()
}
