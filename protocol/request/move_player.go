package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MovePlayer struct {
	Fragments []dto.MoveFragment
}

func (*MovePlayer) Opcode() byte { return 0x18 }

func (m *MovePlayer) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MovePlayer) Deserialize(reader *stream.StreamReader) {
	reader.Read(5)
	m.Fragments = dto.ReadMovements(reader)
}
