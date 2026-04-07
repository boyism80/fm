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
	if err := w.WriteU32(p.OwnerID); err != nil {
		return err
	}
	if err := w.WriteU32(p.OID); err != nil {
		return err
	}
	if p.Animated {
		if err := w.WriteU8(4); err != nil {
			return err
		}
	} else {
		if err := w.WriteU8(1); err != nil {
			return err
		}
	}
	return nil
}

func (p *RemoveSummon) Deserialize(*stream.StreamReader) {}
