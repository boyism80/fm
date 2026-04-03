package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type SpawnDoor struct {
	OwnerID  uint32
	Position types.Vector2[int16]
	Animated bool
}

func (p *SpawnDoor) Opcode() uint16 {
	return 0xCD
}

func (p *SpawnDoor) Serialize(w *stream.StreamWriter) error {
	if p.Animated {
		if err := w.WriteU8(0); err != nil {
			return err
		}
	} else {
		if err := w.WriteU8(1); err != nil {
			return err
		}
	}
	if err := w.WriteU32(p.OwnerID); err != nil {
		return err
	}
	if err := w.Write16(p.Position.X); err != nil {
		return err
	}
	return w.Write16(p.Position.Y)
}

func (p *SpawnDoor) Deserialize(*stream.StreamReader) error { return nil }
