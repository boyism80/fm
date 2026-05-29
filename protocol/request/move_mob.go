package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MoveMob struct {
	OID         uint32
	MovementId  uint16
	ActiveSkill bool
	Unknown2    bool
	Action      int
	CenterSplit int8
	SkillId     uint8
	SkillLevel  uint8
	Unknown     uint8
	SkillDelay  uint8
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
	m.ActiveSkill = flag&0xF != 0
	m.Unknown2 = flag&0xF0 != 0
	centerSplit := reader.Read8()
	m.CenterSplit = centerSplit
	if centerSplit < 0 {
		m.Action = -1
	} else {
		m.Action = int(centerSplit >> 1)
	}
	m.SkillId = reader.ReadU8()
	m.SkillLevel = reader.ReadU8()
	m.Unknown = reader.ReadU8()
	m.SkillDelay = reader.ReadU8()
	reader.Skip(9)
	m.Movements = dto.ReadMovements(reader)
	reader.Skip(9)
}
