package response

import "github.com/boyism80/fm/stream"

type FamilySummon struct {
	Name    string
	MapName string
}

func (p *FamilySummon) Opcode() uint16 {
	return 0x53
}

func (p *FamilySummon) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	w.WriteStr16(p.MapName)
	return nil
}

func (p *FamilySummon) Deserialize(reader *stream.StreamReader) {
}
