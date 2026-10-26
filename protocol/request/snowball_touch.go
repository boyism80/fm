package request

import "github.com/boyism80/fm/stream"

// SnowballTouch is sent when the character runs into a rolling snowball.
type SnowballTouch struct{}

func (*SnowballTouch) Opcode() byte {
	return 0xAB
}

func (p *SnowballTouch) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *SnowballTouch) Deserialize(reader *stream.StreamReader) {}
