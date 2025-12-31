package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type SpawnMob struct {
	Mob       *dto.Mob
	SpawnType constant.MobSpawnType
	Link      uint32
}

func (p *SpawnMob) Opcode() uint16 {
	return 0xA9
}

func (p *SpawnMob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Mob.OID)
	writer.WriteU8(1)
	writer.WriteU32(p.Mob.MobId)
	p.Mob.Serialize(writer)
	writer.Write16(p.Mob.Position.X)
	writer.Write16(p.Mob.Position.Y)
	writer.WriteU8(p.Mob.Stance)
	writer.WriteU16(0)
	writer.Write16(p.Mob.Foothold)
	writer.Write8(int8(p.SpawnType))

	if p.SpawnType == -3 || p.SpawnType >= 0 {
		writer.WriteU32(p.Link)
	}
	writer.Write8(-1)
	writer.WriteU32(0)
	return nil
}

func (p *SpawnMob) Deserialize(reader *stream.StreamReader) error {
	return nil
}
