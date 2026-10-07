package request

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type TeleportStoneList struct {
	Action pconst.TeleportStoneAction
	VIP    bool
	MapID  uint32
}

func (*TeleportStoneList) Opcode() byte { return 0x55 }

func (p *TeleportStoneList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Action))
	writer.WriteBoolean(p.VIP)
	if p.Action == pconst.TeleportStoneActionRemove {
		writer.WriteU32(p.MapID)
	}
	return nil
}

func (p *TeleportStoneList) Deserialize(reader *stream.StreamReader) {
	p.Action = pconst.TeleportStoneAction(reader.ReadU8())
	p.VIP = reader.ReadBool()
	if p.Action == pconst.TeleportStoneActionRemove {
		p.MapID = reader.ReadU32()
	}
}
