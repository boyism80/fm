package req

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/protocol"
)

type MoveMob struct {
	OID          uint32
	MovementId   uint16
	EnabledSkill bool
	Unknown2     bool
	CenterSplit  int8
	Skill1       uint8
	Skill2       uint8
	Skill3       uint8
	Skill4       uint8
	Movements    []protocol.MoveFragment
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

	m.EnabledSkill = flag&0xF != 0
	m.Unknown2 = flag&0xF0 != 0
	m.CenterSplit, err = reader.Read8()
	if err != nil {
		return err
	}

	action := m.CenterSplit
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
	if action < 0 {
		action = -1
	} else {
		action = action >> 1
	}

	m.Movements, err = protocol.ReadMovements(reader)
	if err != nil {
		return err
	}

	err = reader.Skip(9)
	if err != nil {
		return err
	}

	return nil
}
