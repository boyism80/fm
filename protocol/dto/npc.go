package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type Npc struct {
	OID      uint32
	NpcId    uint32
	Position types.Vector2[int16]
	Foothold int16
	Rx0      int16
	Rx1      int16
	Cy       int16
}

func (n *Npc) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(n.NpcId)
	writer.Write16(n.Position.X)
	writer.Write16(n.Position.Y)
	writer.WriteBoolean(false)
	writer.WriteU16(uint16(n.Foothold))
	writer.Write16(n.Rx0)
	writer.Write16(n.Rx1)
	writer.WriteBoolean(false)
	writer.Write16(n.Cy)
	return nil
}
