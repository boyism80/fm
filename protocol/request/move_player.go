package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MovePlayer struct {
	FieldKey  uint8
	Position  types.Vector2[int16]
	Fragments []dto.MoveFragment
	// One arrow key nibble per 30 ms tick of the path: left 8, right 4, up 2, down 1.
	KeyStates []uint8
	Bounds    types.Rect[int16]
}

func (*MovePlayer) Opcode() byte { return 0x18 }

func (m *MovePlayer) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(m.FieldKey)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.WriteU8(uint8(len(m.Fragments)))
	for _, frag := range m.Fragments {
		frag.Serialize(writer)
	}
	writer.WriteU8(uint8(len(m.KeyStates)))
	for i := 0; i < len(m.KeyStates); i += 2 {
		keys := m.KeyStates[i] & 0xF
		if i+1 < len(m.KeyStates) {
			keys |= m.KeyStates[i+1] << 4
		}
		writer.WriteU8(keys)
	}
	writer.Write16(m.Bounds.Left)
	writer.Write16(m.Bounds.Top)
	writer.Write16(m.Bounds.Right)
	writer.Write16(m.Bounds.Bottom)
	return nil
}

func (m *MovePlayer) Deserialize(reader *stream.StreamReader) {
	m.FieldKey = reader.ReadU8()
	m.Position.X = reader.Read16()
	m.Position.Y = reader.Read16()
	m.Fragments = dto.ReadMovements(reader)
}
