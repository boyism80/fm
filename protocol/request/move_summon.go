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

func (m *MoveSummon) Deserialize(reader *stream.StreamReader) error {
	var err error
	m.OID, err = reader.ReadU32()
	if err != nil {
		return err
	}
	x, err := reader.Read16()
	if err != nil {
		return err
	}
	y, err := reader.Read16()
	if err != nil {
		return err
	}
	m.Position = types.Vector2[int16]{X: x, Y: y}
	m.Fragments, err = dto.ReadMovements(reader)
	if err != nil {
		return err
	}

	return nil
}
