package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MovePet struct {
	Position  types.Vector2[int16]
	Fragments []dto.MoveFragment
}

func (*MovePet) Opcode() byte { return 0x84 }

func (m *MovePet) Serialize(writer *stream.StreamWriter) error {
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.WriteU8(uint8(len(m.Fragments)))
	for _, frag := range m.Fragments {
		frag.Serialize(writer)
	}
	return nil
}

func (m *MovePet) Deserialize(reader *stream.StreamReader) {
	m.Position = types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
	m.Fragments = dto.ReadMovements(reader)
}
