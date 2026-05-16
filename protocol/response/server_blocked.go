package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type ServerBlocked struct {
	Reason constant.ServerBlockedReason
}

func (p *ServerBlocked) Opcode() uint16 {
	return 0x5A
}

func (p *ServerBlocked) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Reason))
	return nil
}

func (p *ServerBlocked) Deserialize(reader *stream.StreamReader) {
}
