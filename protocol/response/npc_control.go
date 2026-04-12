package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type NpcControl struct {
	NPC     *dto.Npc
	MiniMap bool
}

func (p *NpcControl) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteU32(p.NPC.OID)
	writer.WriteU32(p.NPC.NpcId)
	writer.Write16(p.NPC.Position.X)
	writer.Write16(p.NPC.Cy)
	writer.WriteU8(0)
	writer.Write16(p.NPC.Foothold)
	writer.Write16(p.NPC.Rx0)
	writer.Write16(p.NPC.Rx1)
	writer.WriteBoolean(p.MiniMap)
	return nil
}

func (s *NpcControl) Deserialize(reader *stream.StreamReader) {
}

func (p *NpcControl) Opcode() uint16 {
	return 0xBD
}
