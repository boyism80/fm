package request

import "github.com/boyism80/fm/stream"

type NpcAction struct {
	Bytes []byte
}

func (p *NpcAction) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *NpcAction) Deserialize(reader *stream.StreamReader) error {
	p.Bytes, _ = reader.Read(reader.Remaining())
	return nil
}
