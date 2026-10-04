package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
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
	writer.Write8(int8(p.Mob.CarnivalTeam))
	writer.WriteU32(0)
	return nil
}

func (p *SpawnMob) Deserialize(reader *stream.StreamReader) {
	p.Mob = &dto.Mob{OID: reader.ReadU32()}
	reader.Skip(1)
	p.Mob.MobId = reader.ReadU32()
	p.Mob.StatusMask = reader.Read32()
	if p.Mob.StatusMask != 0 {
		return
	}
	p.Mob.Position = types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
	p.Mob.Stance = reader.ReadU8()
	reader.Skip(2)
	p.Mob.Foothold = reader.Read16()
	p.SpawnType = constant.MobSpawnType(reader.Read8())
	if p.SpawnType == -3 || p.SpawnType >= 0 {
		p.Link = reader.ReadU32()
	}
	p.Mob.CarnivalTeam = constant.CarnivalTeam(reader.Read8())
	reader.Skip(4)
}
