package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

type NpcControl struct {
	NPC     *entity.Npc
	MiniMap bool
}

func (p *NpcControl) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteU32(p.NPC.OID)
	writer.WriteU32(p.NPC.Spec.ID)
	writer.Write16(p.NPC.Spec.Position.X)
	writer.Write16(p.NPC.Spec.CollisionY)
	writer.WriteU8(uint8(p.NPC.Spec.FacingDirection))
	writer.Write16(p.NPC.Spec.Foothold)
	writer.Write16(p.NPC.Spec.RenderX0)
	writer.Write16(p.NPC.Spec.RenderX1)
	writer.WriteBoolean(p.MiniMap)
	return nil
}

func (s *NpcControl) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *NpcControl) Opcode() uint16 {
	return 0xBD // NpcControl opcode
}
