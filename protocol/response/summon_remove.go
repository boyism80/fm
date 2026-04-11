package response

import "github.com/boyism80/fm/stream"

type RemoveSummon struct {
	OwnerID  uint32
	OID      uint32
	Animated bool
}

func (p *RemoveSummon) Opcode() uint16 {
	return 0x7C
}

func (p *RemoveSummon) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.OwnerID)
	w.WriteU32(p.OID)
	if p.Animated {
		w.WriteU8(4)
	} else {
		w.WriteU8(1)
	}
	return nil
}

func (p *RemoveSummon) Deserialize(*stream.StreamReader) {}
