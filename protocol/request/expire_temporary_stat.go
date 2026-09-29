package request

import "github.com/boyism80/fm/stream"

type ExpireTemporaryStat struct{}

func (*ExpireTemporaryStat) Opcode() byte {
	return 0x52
}

func (p *ExpireTemporaryStat) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *ExpireTemporaryStat) Deserialize(reader *stream.StreamReader) {}
