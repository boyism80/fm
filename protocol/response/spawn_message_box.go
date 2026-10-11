package response

import "github.com/boyism80/fm/stream"

type SpawnMessageBox struct {
	OID    uint32
	ItemID uint32
	Text   string
	Owner  string
	X      int16
	Y      int16
}

func (*SpawnMessageBox) Opcode() uint16 { return 0xC9 }

func (p *SpawnMessageBox) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.OID)
	w.WriteU32(p.ItemID)
	w.WriteStr16(p.Text)
	w.WriteStr16(p.Owner)
	w.Write16(p.X)
	w.Write16(p.Y)
	return nil
}

func (p *SpawnMessageBox) Deserialize(r *stream.StreamReader) {
	p.OID = r.ReadU32()
	p.ItemID = r.ReadU32()
	p.Text = r.ReadStr16()
	p.Owner = r.ReadStr16()
	p.X = r.Read16()
	p.Y = r.Read16()
}
