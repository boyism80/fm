package response

import "github.com/boyism80/fm/stream"

type RemoveMist struct {
	OID      uint32
	Eruption bool
}

func (p *RemoveMist) Opcode() uint16 {
	return 0xCC
}

func (p *RemoveMist) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.OID)
	eruption := uint8(0)
	if p.Eruption {
		eruption = 1
	}
	w.WriteU8(eruption)
	return nil
}

func (p *RemoveMist) Deserialize(*stream.StreamReader) {}
