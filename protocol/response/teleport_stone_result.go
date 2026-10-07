package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type TeleportStoneResult struct {
	Result pconst.TeleportStoneResult
	VIP    bool
	Maps   []uint32
}

func (p *TeleportStoneResult) Opcode() uint16 {
	return 0x1E
}

func (p *TeleportStoneResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	writer.WriteBoolean(p.VIP)
	if p.Result != pconst.TeleportStoneResultList {
		return nil
	}
	for _, mapID := range p.Maps {
		writer.WriteU32(mapID)
	}
	return nil
}

func (p *TeleportStoneResult) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.TeleportStoneResult(reader.ReadU8())
	p.VIP = reader.ReadBool()
	if p.Result != pconst.TeleportStoneResultList {
		return
	}
	count := constant.TeleportStoneCount
	if p.VIP {
		count = constant.VipTeleportStoneCount
	}
	p.Maps = make([]uint32, count)
	for i := range p.Maps {
		p.Maps[i] = reader.ReadU32()
	}
}
