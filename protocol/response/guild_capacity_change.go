package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildCapacityChange struct {
	GuildID  uint32
	Capacity uint8
}

func (p *GuildCapacityChange) Opcode() uint16 {
	return 0x30
}

func (p *GuildCapacityChange) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CCapacityChange))
	w.WriteU32(p.GuildID)
	w.WriteU8(p.Capacity)
	return nil
}

func (p *GuildCapacityChange) Deserialize(*stream.StreamReader) {}
