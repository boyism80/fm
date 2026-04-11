package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MoveSummon struct {
	OID       uint32
	Position  types.Vector2[int16]
	Fragments []dto.MoveFragment
}

func (m *MoveSummon) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MoveSummon) Deserialize(reader *stream.StreamReader) {
	m.OID = reader.ReadU32()
	m.Position = types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
	m.Fragments = dto.ReadMovements(reader)

}
