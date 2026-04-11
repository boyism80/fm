package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MoveSummon struct {
	CharacterID uint32
	OID         uint32
	StartPoint  types.Vector2[int16]
	Fragments   []dto.MoveFragment
}

func (p *MoveSummon) Opcode() uint16 {
	return 0x7D
}

func (p *MoveSummon) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.OID)
	w.Write16(p.StartPoint.X)
	w.Write16(p.StartPoint.Y)
	w.WriteU8(uint8(len(p.Fragments)))
	for _, m := range p.Fragments {
		m.Serialize(w)
	}
	return nil
}

func (p *MoveSummon) Deserialize(*stream.StreamReader) {}
