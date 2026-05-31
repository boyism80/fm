package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildMessage struct {
	Code constant.GuildResponseCode
}

func (p *GuildMessage) Opcode() uint16 {
	return 0x30
}

func (p *GuildMessage) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Code))
	return nil
}

func (p *GuildMessage) Deserialize(*stream.StreamReader) {}
