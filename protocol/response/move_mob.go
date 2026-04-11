package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MoveMob struct {
	IsAggroed   bool
	CenterSplit int8
	Skill1      uint8
	Skill2      uint8
	Skill3      uint8
	Skill4      uint8
	OID         uint32
	StartPoint  types.Vector2[int16]
	Movements   []dto.MoveFragment
}

func (p *MoveMob) Opcode() uint16 {
	return 0xAC
}

func (p *MoveMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.WriteBoolean(p.IsAggroed)
	writer.Write8(p.CenterSplit)
	writer.WriteU8(p.Skill1)
	writer.WriteU8(p.Skill2)
	writer.WriteU8(p.Skill3)
	writer.WriteU8(p.Skill4)
	writer.Write16(p.StartPoint.X)
	writer.Write16(p.StartPoint.Y)
	writer.WriteU8(uint8(len(p.Movements)))
	for _, move := range p.Movements {
		move.Serialize(writer)
	}
	return nil
}

func (p *MoveMob) Deserialize(reader *stream.StreamReader) {
}
