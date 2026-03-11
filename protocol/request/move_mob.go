package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MoveMob struct {
	OID         uint32
	MovementId  uint16
	IsAggroed   bool
	Unknown2    bool
	CenterSplit int8
	Skill1      uint8
	Skill2      uint8
	Skill3      uint8
	Skill4      uint8
	Movements   []dto.MoveFragment
}

func (m *MoveMob) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MoveMob) Deserialize(reader *stream.StreamReader) error {
	var err error
	m.OID, err = reader.ReadU32()
	if err != nil {
		return err
	}

	m.MovementId, err = reader.ReadU16()
	if err != nil {
		return err
	}

	flag, err := reader.ReadU8()
	if err != nil {
		return err
	}

	m.IsAggroed = flag&0xF != 0
	m.Unknown2 = flag&0xF0 != 0
	m.CenterSplit, err = reader.Read8()
	if err != nil {
		return err
	}

	centerSplit := m.CenterSplit
	m.Skill1, err = reader.ReadU8()
	if err != nil {
		return err
	}

	m.Skill2, err = reader.ReadU8()
	if err != nil {
		return err
	}
	m.Skill3, err = reader.ReadU8()
	if err != nil {
		return err
	}

	m.Skill4, err = reader.ReadU8()
	if err != nil {
		return err
	}

	reader.Skip(9)
	if centerSplit < 0 {
		centerSplit = -1
	} else {
		centerSplit = centerSplit >> 1
	}

	m.Movements, err = dto.ReadMovements(reader)
	if err != nil {
		return err
	}

	err = reader.Skip(9)
	if err != nil {
		return err
	}

	return nil
}
