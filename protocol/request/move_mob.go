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

func (*MoveMob) Opcode() byte { return 0x95 }

func (m *MoveMob) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MoveMob) Deserialize(reader *stream.StreamReader) {
	m.OID = reader.ReadU32()

	m.MovementId = reader.ReadU16()

	flag := reader.ReadU8()

	m.IsAggroed = flag&0xF != 0
	m.Unknown2 = flag&0xF0 != 0
	m.CenterSplit = reader.Read8()
	m.Skill1 = reader.ReadU8()

	m.Skill2 = reader.ReadU8()
	m.Skill3 = reader.ReadU8()

	m.Skill4 = reader.ReadU8()

	reader.Skip(9)

	m.Movements = dto.ReadMovements(reader)

	reader.Skip(9)

}
