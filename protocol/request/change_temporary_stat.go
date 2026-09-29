package request

import "github.com/boyism80/fm/stream"

type ChangeTemporaryStat struct{}

func (*ChangeTemporaryStat) Opcode() byte {
	return 0x5B
}

func (p *ChangeTemporaryStat) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *ChangeTemporaryStat) Deserialize(reader *stream.StreamReader) {}
