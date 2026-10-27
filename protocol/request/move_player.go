package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MovePlayer struct {
	FieldKey  uint8
	Fragments []dto.MoveFragment
}

func (*MovePlayer) Opcode() byte { return 0x18 }

func (m *MovePlayer) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(m.FieldKey)
	writer.Write(make([]byte, 4))
	writer.WriteU8(uint8(len(m.Fragments)))
	for _, frag := range m.Fragments {
		frag.Serialize(writer)
	}
	return nil
}

func (m *MovePlayer) Deserialize(reader *stream.StreamReader) {
	m.FieldKey = reader.ReadU8()
	reader.Read(4)
	m.Fragments = dto.ReadMovements(reader)
}
