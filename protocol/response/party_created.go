package response

import "github.com/boyism80/fm/stream"

type PartyCreated struct {
	PartyID uint32
}

func (p *PartyCreated) Opcode() uint16 {
	return 0x2D
}

func (p *PartyCreated) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(8)
	w.WriteU32(p.PartyID)
	w.WriteU32(999999999)
	w.WriteU32(999999999)
	w.Write64(0)
	w.WriteU8(0)
	w.WriteU8(1)
	return nil
}

func (p *PartyCreated) Deserialize(*stream.StreamReader) {}
