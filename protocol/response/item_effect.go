package response

import "github.com/boyism80/fm/stream"

type ItemEffect struct {
	CharacterID uint32
	ItemID      uint32
}

func (*ItemEffect) Opcode() uint16 { return 0x8B }
func (p *ItemEffect) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.ItemID)
	return nil
}
func (p *ItemEffect) Deserialize(r *stream.StreamReader) {
	p.CharacterID = r.ReadU32()
	p.ItemID = r.ReadU32()
}
