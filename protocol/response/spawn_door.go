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
		w.WriteU8(0)
	} else {
		w.WriteU8(1)
	}
	w.WriteU32(p.OwnerID)
	w.Write16(p.Position.X)
	w.Write16(p.Position.Y)
	return nil
}

func (p *SpawnDoor) Deserialize(*stream.StreamReader) {}
