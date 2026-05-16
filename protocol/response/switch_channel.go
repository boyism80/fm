package response

import (
	"github.com/boyism80/fm/stream"
)

type SwitchChannel struct {
	IP   string
	Port uint16
}

func (p *SwitchChannel) Opcode() uint16 {
	return 0x08
}

func (p *SwitchChannel) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteIPAddress(p.IP)
	writer.WriteU16(p.Port)
	writer.WriteU8(0)
	return nil
}

func (p *SwitchChannel) Deserialize(reader *stream.StreamReader) {
}
