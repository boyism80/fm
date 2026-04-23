package request

import "github.com/boyism80/fm/stream"

type NpcAction struct {
	Bytes []byte
}

func (*NpcAction) Opcode() byte { return 0x9E }

func (p *NpcAction) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *NpcAction) Deserialize(reader *stream.StreamReader) {
	p.Bytes = reader.Read(reader.Remaining())
}
