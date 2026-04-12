package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type SpawnNpc struct {
	NPC     *dto.Npc
	Visible bool
}

func (p *SpawnNpc) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.NPC.OID)
	writer.WriteU32(p.NPC.NpcId)
	writer.Write16(p.NPC.Position.X)
	writer.Write16(p.NPC.Cy)
	writer.WriteU8(0)
	writer.Write16(p.NPC.Foothold)
	writer.Write16(p.NPC.Rx0)
	writer.Write16(p.NPC.Rx1)
	writer.WriteBoolean(p.Visible)
	return nil
}

func (s *SpawnNpc) Deserialize(reader *stream.StreamReader) {
}

func (p *SpawnNpc) Opcode() uint16 {
	return 0xBB
}
