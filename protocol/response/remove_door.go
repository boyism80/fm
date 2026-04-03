package response

import "github.com/boyism80/fm/stream"

type RemoveDoor struct {
	OwnerID  uint32
	Animated bool
}

func (p *RemoveDoor) Opcode() uint16 {
	return 0xCE
}

func (p *RemoveDoor) Serialize(w *stream.StreamWriter) error {
	if p.Animated {
		if err := w.WriteU8(0); err != nil {
			return err
		}
	} else {
		if err := w.WriteU8(1); err != nil {
			return err
		}
	}
	return w.WriteU32(p.OwnerID)
}

func (p *RemoveDoor) Deserialize(*stream.StreamReader) error { return nil }
