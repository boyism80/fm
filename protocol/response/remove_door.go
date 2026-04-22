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
	w.WriteBoolean(!p.Animated)
	w.WriteU32(p.OwnerID)
	return nil
}

func (p *RemoveDoor) Deserialize(*stream.StreamReader) {}
