package resp

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/game/entity"
)

type StartControlMob struct {
	Mob   *entity.Mob
	Aggro bool
}

func (p *StartControlMob) Opcode() uint16 {
	return 0xAB
}

type StopControlMob struct {
	OID uint32
}

func (p *StopControlMob) Opcode() uint16 {
	return 0xAB
}

func (p *StartControlMob) Serialize(writer *stream.StreamWriter) error {
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

func (p *StartControlMob) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *StopControlMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.Opcode())
	writer.WriteU8(0)
	writer.WriteU32(p.OID)
	return nil
}

func (p *StopControlMob) Deserialize(reader *stream.StreamReader) error {
	return nil
}
