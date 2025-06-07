package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

type ControlMob struct {
	Mob   *entity.Mob
	Aggro bool
}

func (p *ControlMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xAB)
	if p.Aggro {
		writer.WriteU8(2)
	} else {
		writer.WriteU8(1)
	}
	writer.WriteU32(p.Mob.OID)
	writer.WriteU8(1)
	writer.WriteU32(p.Mob.Spec.ID)
	p.Mob.Serialize(writer)
	writer.Write16(p.Mob.Position.X)
	writer.Write16(p.Mob.Position.Y)
	writer.WriteU8(p.Mob.Stance)
	writer.WriteU16(0)
	writer.Write16(p.Mob.Foothold)
	writer.Write8(-1)
	writer.Write8(-1)
	writer.WriteU32(0)
	return nil
}

func (m *ControlMob) Deserialize(reader *stream.StreamReader) error {
	return nil
}
