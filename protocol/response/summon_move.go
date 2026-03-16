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
	if err := w.WriteU32(p.CharacterID); err != nil {
		return err
	}
	if err := w.WriteU32(p.OID); err != nil {
		return err
	}
	if err := w.Write16(p.StartPoint.X); err != nil {
		return err
	}
	if err := w.Write16(p.StartPoint.Y); err != nil {
		return err
	}
	if err := w.WriteU8(uint8(len(p.Fragments))); err != nil {
		return err
	}
	for _, m := range p.Fragments {
		m.Serialize(w)
	}
	return nil
}

func (p *MoveSummon) Deserialize(*stream.StreamReader) error { return nil }
