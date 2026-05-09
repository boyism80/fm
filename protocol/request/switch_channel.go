package request

import (
	"github.com/boyism80/fm/stream"
)

type SwitchChannel struct {
	Channel uint8
}

func (*SwitchChannel) Opcode() byte {
	return 0x16
}

func (a *SwitchChannel) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *SwitchChannel) Deserialize(reader *stream.StreamReader) {
	a.Channel = reader.ReadU8()
}
