package request

import "github.com/boyism80/fm/stream"

// PulleyPull is sent when the character pulls the pulley that lowers the field's healer.
type PulleyPull struct{}

func (*PulleyPull) Opcode() byte {
	return 0xAE
}

func (p *PulleyPull) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *PulleyPull) Deserialize(reader *stream.StreamReader) {}
