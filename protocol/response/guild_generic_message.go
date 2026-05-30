package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildGenericMessage struct {
	Code constant.GuildResponseCode
}

func (p *GuildGenericMessage) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildGenericMessage) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Code))
	return nil
}

func (p *GuildGenericMessage) Deserialize(*stream.StreamReader) {}
