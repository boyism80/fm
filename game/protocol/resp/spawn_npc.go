package resp

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/game/entity"
)

type SpawnNpc struct {
	NPC     *entity.Npc
	Visible bool
}

func (p *SpawnNpc) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.NPC.OID)
	writer.WriteU32(p.NPC.Wz.ID)
	writer.Write16(p.NPC.Wz.Position.X)
	writer.Write16(p.NPC.Wz.CollisionY)
	writer.WriteU8(uint8(p.NPC.Wz.FacingDirection))
	writer.Write16(p.NPC.Wz.Foothold)
	writer.Write16(p.NPC.Wz.RenderX0)
	writer.Write16(p.NPC.Wz.RenderX1)
	writer.WriteBoolean(p.Visible)
	return nil
}

func (s *SpawnNpc) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *SpawnNpc) Opcode() uint16 {
	return 0xBB // SpawnNpc opcode
}
