package response

import "github.com/boyism80/fm/stream"

type RemoveMessageBox struct {
	OID      uint32
	Animated bool
}

func (*RemoveMessageBox) Opcode() uint16 { return 0xCA }

// The client plays the removal animation when the flag byte is 0.
func (p *RemoveMessageBox) Serialize(w *stream.StreamWriter) error {
	w.WriteBoolean(!p.Animated)
	w.WriteU32(p.OID)
	return nil
}

func (p *RemoveMessageBox) Deserialize(r *stream.StreamReader) {
	p.Animated = r.ReadU8() == 0
	p.OID = r.ReadU32()
}
