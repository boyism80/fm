package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MovePet struct {
	CharacterID uint32
	StartPoint  types.Vector2[int16]
	Fragments   []dto.MoveFragment
}

func (p *MovePet) Opcode() uint16 {
	return 0x76
}

func (p *MovePet) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.Write16(p.StartPoint.X)
	w.Write16(p.StartPoint.Y)
	w.WriteU8(uint8(len(p.Fragments)))
	for _, m := range p.Fragments {
		m.Serialize(w)
	}
	return nil
}

func (p *MovePet) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.StartPoint = types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
	p.Fragments = dto.ReadMovements(reader)
}
